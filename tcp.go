package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

//TCPBufferSize 缓存大小
const TCPBufferSize = 4096

//TCPServer TCP服务
type TCPServer struct {
	Ctx         context.Context
	tcpAddress  *net.TCPAddr
	tcpListener *net.TCPListener
	tcpPort     string
	Handler     func(func([]byte) error, []byte) error
	Logger      ILogger
}

//NewTCPServer 新建
func NewTCPServer(ctx context.Context, port string, h func(func([]byte) error, []byte) error, logger ILogger) *TCPServer {
	if h == nil {
		logger.Fatal("TCP Handler不为nil")
		return nil
	}
	tcpAddress, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		logger.Fatal("ResolveTCPAddr失败:", err.Error())
		return nil
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		logger.Fatal("TCP监听端口失败:", err.Error())
		return nil
	}
	s := &TCPServer{
		Ctx:         ctx,
		tcpAddress:  tcpAddress,
		tcpListener: listener,
		tcpPort:     port,
		Handler:     h,
		Logger:      logger,
	}
	return s
}

//Run 运行
func (s *TCPServer) Run() {
	go func() {
		var closeOnce sync.Once
		<-s.Ctx.Done()
		closeOnce.Do(func() {
			if err := s.tcpListener.Close(); err != nil {
				s.Logger.Error("TCP监听端口关闭失败:", err.Error())
			} else {
				s.Logger.Debug("TCP监听端口关闭。")
			}
		})
	}()
	s.Logger.Debug("TCP监听端口", s.tcpPort)
	s.Logger.Debug("TCP已初始化连接，等待客户端连接……")
	tcpGroup := sync.WaitGroup{}
	//设置重试延迟时间
	var tempDelay time.Duration
	for {
		conn, err := s.tcpListener.AcceptTCP()
		if err != nil {
			if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
				s.Logger.Warn("Temporary error when accepting new connections:", err.Error())
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					//设置延迟时间，开始递增，避免无意义的重试
					tempDelay *= 2
				}
				//限制了最大重试时间在1s
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				s.Logger.Warn("Permanent error when accepting new connections:", err.Error())
			}
			break
		}
		//如果没有网络错误，下次网络错误的重试时间重新开始
		tempDelay = 0
		if err = conn.SetNoDelay(false); err != nil {
			s.Logger.Error("TCP设定操作系统是否应该延迟数据包传递失败:" + err.Error())
		}
		//设置TCP保持长连接
		conn.SetKeepAlive(true)
		//设置TCP探测时间间隔时间为3分钟，如果客户端3分钟没有和服务端通信，则开始探测
		conn.SetKeepAlivePeriod(3 * time.Minute)
		conn.SetLinger(10)
		session := &TCPSession{
			Ctx:        s.Ctx,
			conn:       conn,
			Handler:    s.Handler,
			serverNode: true,
			Logger:     s.Logger,
		}
		tcpGroup.Add(1)
		go func() {
			session.tcpReceive()
			tcpGroup.Done()
		}()
	}
	s.Logger.Debug("TCP等待子协程关闭……")
	tcpGroup.Wait()
	s.Logger.Debug("TCPServer关闭。")
}

//TCPSession 会话
type TCPSession struct {
	Ctx        context.Context
	conn       *net.TCPConn
	Handler    func(func([]byte) error, []byte) error
	serverNode bool
	stopChan   chan struct{}
	Logger     ILogger
}

func (ts *TCPSession) tcpReceive() {
	defer func() {
		if ts.serverNode {
			if err := ts.conn.Close(); err != nil {
				ts.Logger.Error(err.Error())
			}
		} else {
			close(ts.stopChan)
		}
		if r := recover(); r != nil {
			ts.Logger.Error(fmt.Sprintf("tcpReceive：%v \n%s", r, string(debug.Stack())))
		}
		ts.Logger.Debug(ts.conn.RemoteAddr().String(), " tcpReceive stop")
	}()
	//IO读缓存
	buf := make([]byte, TCPBufferSize)
	//buf 写、读位置序号
	var iw, ir int
	var err error
	for {
		//TODO 优化
		select {
		case <-ts.Ctx.Done():
			if ts.serverNode {
				goto end
			}
		default:
		}
		if err = ts.conn.SetReadDeadline(time.Now().Add(DefaultDeadlineDuration)); err != nil {
			goto end
		}
		n, errX := ts.conn.Read(buf[iw:])
		iw += n
		if errX != nil {
			err = errX
			goto end
		}
		//处理数据
		for iw >= (ir + 6) {
			length := int(BytesToUint32(buf[ir : ir+4]))
			tail := ir + length
			if tail > iw {
				break
			}
			flag := BytesToUint16(buf[ir+4 : ir+6])
			if flag == StatusGoaway16 {
				ir += 6
				goto end
			}
			if flag == StatusPing16 {
				ir += 6
				err = ts.Send(framePong)
				if err != nil {
					goto end
				}
				continue
			}
			if flag == StatusPong16 {
				ir += 6
				continue
			}
			if err := ts.Handler(ts.Send, buf[ir:tail]); err != nil {
				ts.Logger.Error("tcpReceive：" + err.Error())
			}
			ir = tail
		}
		if ir < iw {
			copy(buf[0:], buf[ir:iw])
			iw -= ir
		} else {
			iw = 0
		}
		ir = 0
	}
end:
	if err != nil && err != io.EOF {
		if !strings.Contains(err.Error(), "wsarecv: An existing connection was forcibly closed by the remote host.") && !strings.Contains(err.Error(), "use of closed network connection") {
			ts.Logger.Error(ts.conn.RemoteAddr().String(), " tcpReceive：", err.Error())
		} else {
			ts.Logger.Debug(ts.conn.RemoteAddr().String(), " tcpReceive：", err.Error())
		}
	}
}

//TCPDial 连接
func TCPDial(ctx context.Context, url string, h func(func([]byte) error, []byte) error, logger ILogger) (*TCPSession, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", url)
	if err != nil {
		return nil, errors.New("TCPDial tcpAddr fail:" + err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, errors.New("TCPDial 连接服务端失败:" + err.Error())
	}
	if conn.SetNoDelay(false) != nil {
		return nil, errors.New("TCPDial 设定操作系统是否应该延迟数据包传递失败:" + err.Error())
	}
	c := &TCPSession{
		Ctx:        ctx,
		conn:       conn,
		Handler:    h,
		serverNode: false,
		stopChan:   make(chan struct{}),
		Logger:     logger,
	}
	go c.tcpReceive()
	go c.tcpSend()
	return c, nil
}

func (ts *TCPSession) tcpSend() {
	ticker := time.NewTicker(DefaultHeartbeatPeriodDuration)
	defer func() {
		ticker.Stop()
		if r := recover(); r != nil {
			ts.Logger.Error(fmt.Sprintf("tcpSendC：%v %s", r, string(debug.Stack())))
		}
		ts.Logger.Debug(ts.conn.RemoteAddr().String(), " tcpSend stop")
	}()
	for {
		select {
		case <-ticker.C:
			if err := ts.Send(framePing); err != nil {
				ts.Logger.Error(err.Error())
				return
			}
		case <-ts.Ctx.Done():
			if err := ts.Send(frameGoaway); err != nil {
				ts.Logger.Error(err.Error())
			}
			return
		case <-ts.stopChan:
			return
		}
	}
}

//Send 发送
func (ts *TCPSession) Send(message []byte) error {
	if err := ts.conn.SetWriteDeadline(time.Now().Add(DefaultDeadlineDuration)); err != nil {
		return err
	}
	if _, err := ts.conn.Write(message); err != nil {
		return err
	}
	return nil
}

// http://xiaorui.cc/archives/6402
