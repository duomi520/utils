package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

//IMux 路由
type IMux interface {
	Vars(*http.Request) map[string]string
	WarpHandleFunc(string, string, func(http.ResponseWriter, *http.Request))
	PathPrefix(string)
	Handler() http.Handler
}

//HTTPContext 上下文
type HTTPContext struct {
	index   *int
	chain   []func(HTTPContext)
	do      func(HTTPContext)
	Vars    map[string]string
	Writer  http.ResponseWriter
	Request *http.Request
	//validator
	Verify func(any, string) error
	Struct func(any) error
}

//Params 请求参数
func (c HTTPContext) Params(s string) string {
	if v, ok := c.Vars[s]; ok {
		return v
	}
	return c.Request.FormValue(s)
}

//BindJSON 绑定JSON数据
func (c HTTPContext) BindJSON(v any) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("readAll faile: %w", err)
	}
	err = json.Unmarshal(buf, v)
	if err != nil {
		return fmt.Errorf("unmarshal %v faile: %w", v, err)
	}
	err = c.Struct(v)
	if err != nil {
		return fmt.Errorf("validator %v faile: %w", v, err)
	}
	return nil
}

//String 带有状态码的纯文本响应
func (c HTTPContext) String(status int, msg string) {
	c.Writer.WriteHeader(status)
	io.WriteString(c.Writer, msg)
}

//JSON 带有状态码的JSON 数据
func (c HTTPContext) JSON(status int, v any) {
	d, err := json.Marshal(v)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	c.Writer.Write(d)
}

//Next 下一个
func (c HTTPContext) Next() {
	*c.index++
	if *c.index >= len(c.chain) {
		c.do(c)
	} else {
		c.chain[*c.index](c)
	}
}

//WRoute w
type WRoute struct {
	muxVars           func(*http.Request) map[string]string
	muxWarpHandleFunc func(string, string, func(http.ResponseWriter, *http.Request))
	muxPathPrefix     func(string)
	muxHandler        func() http.Handler
	validatorVar      func(any, string) error
	validatorStruct   func(any) error
	DebugMode         bool
}

//NewRoute 新建
func NewRoute(ops *Options) *WRoute {
	r := &WRoute{}
	if ops.Mux == nil {
		panic("Mux is nil")
	}
	r.muxVars = ops.Mux.Vars
	r.muxWarpHandleFunc = ops.Mux.WarpHandleFunc
	r.muxPathPrefix = ops.Mux.PathPrefix
	r.muxHandler = ops.Mux.Handler
	if ops.Validator == nil {
		panic("Validator is nil")
	}
	r.validatorVar = ops.Validator.Var
	r.validatorStruct = ops.Validator.Struct
	if ops.Logger == nil {
		panic("Logger is nil")
	}
	defaultHTTPLogger = ops.Logger
	return r
}

//PathPrefix 前缀
func (r *WRoute) PathPrefix(tpl string) {
	r.muxPathPrefix(tpl)
}

//HandleFunc 处理
func (r *WRoute) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.muxWarpHandleFunc("GET", pattern, handler)
}

//API a
func (r *WRoute) API(g []func(HTTPContext), url string, fn func(any) (any, error)) {
	// TODO
}

//GET g
func (r *WRoute) GET(g []func(HTTPContext), url string, fn func(HTTPContext)) {
	r.muxWarpHandleFunc("GET", url, r.Warp(g, fn))
}

//POST p
func (r *WRoute) POST(g []func(HTTPContext), url string, fn func(HTTPContext)) {
	r.muxWarpHandleFunc("POST", url, r.Warp(g, fn))
}

//DELETE d
func (r *WRoute) DELETE(g []func(HTTPContext), url string, fn func(HTTPContext)) {
	r.muxWarpHandleFunc("DELETE", url, r.Warp(g, fn))
}

//Warp 封装
func (r *WRoute) Warp(g []func(HTTPContext), fn func(HTTPContext)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				buf := make([]byte, 4096)
				lenght := runtime.Stack(buf, false)
				defaultHTTPLogger.Error(fmt.Sprintf("WRoute.warp %v \n%s", v, buf[:lenght]))
				if r.DebugMode {
					rw.Write([]byte("\n"))
					rw.Write(buf[:lenght])
				}
			}
		}()
		var index int
		c := HTTPContext{index: &index, chain: g, do: fn, Writer: rw, Request: req, Verify: r.validatorVar, Struct: r.validatorStruct}
		c.Vars = r.muxVars(req)
		if len(g) > 0 {
			c.chain[0](c)
		} else {
			fn(c)
		}
	}
}

//Middleware 使用中间件
func Middleware(m ...func(HTTPContext)) []func(HTTPContext) {
	return m
}

//AppendValidator 追加
func AppendValidator(base []func(HTTPContext), a ...string) []func(HTTPContext) {
	g := append(base, ValidatorMiddleware(a...))
	return g
}

//ValidatorMiddleware 输入参数验证
func ValidatorMiddleware(a ...string) func(HTTPContext) {
	var s0, s1 []string
	for _, v := range a {
		s := strings.Split(v, ":")
		if len(s) != 2 {
			panic(fmt.Sprintf("validate:bad describe %s", v))
		}
		s0 = append(s0, s[0])
		s1 = append(s1, s[1])
	}
	return func(c HTTPContext) {
		for i := range s0 {
			if err := c.Verify(c.Params(s0[i]), s1[i]); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("validate %s %s faile: %s", s0[i], s1[i], err.Error()))
				return
			}
		}
		c.Next()
	}
}

var defaultHTTPLogger ILogger

//LoggerMiddleware 日志
func LoggerMiddleware() func(HTTPContext) {
	var startTime time.Time
	var rw HTTPLoggerResponseWriter
	return func(c HTTPContext) {
		startTime = time.Now()
		rw.w = c.Writer
		c.Writer = &rw
		c.Next()
		if rw.status > 299 {
			if rw.err != nil {
				defaultHTTPLogger.Error(fmt.Sprintf("| %13v | %15s | %5d | %7s | %s | %s ", time.Since(startTime), c.Request.RemoteAddr, rw.status, c.Request.Method, c.Request.URL, rw.err.Error()))
			} else {
				defaultHTTPLogger.Warn(fmt.Sprintf("| %13v | %15s | %5d | %7s | %s | %s ", time.Since(startTime), c.Request.RemoteAddr, rw.status, c.Request.Method, c.Request.URL, rw.result))
			}
		} else {
			defaultHTTPLogger.Debug(fmt.Sprintf("| %13v | %15s | %5d | %7s | %s | %s ", time.Since(startTime), c.Request.RemoteAddr, rw.status, c.Request.Method, c.Request.URL, rw.result))
		}
	}
}

//HTTPLoggerResponseWriter d
type HTTPLoggerResponseWriter struct {
	status int
	result []byte
	err    error
	w      http.ResponseWriter
}

//Header 返回一个Header类型值
func (h *HTTPLoggerResponseWriter) Header() http.Header {
	return h.w.Header()
}

//WriteHeader 该方法发送HTTP回复的头域和状态码
func (h *HTTPLoggerResponseWriter) WriteHeader(s int) {
	h.status = s
	h.w.WriteHeader(s)
}

//Write 向连接中写入作为HTTP的一部分回复的数据
func (h *HTTPLoggerResponseWriter) Write(d []byte) (int, error) {
	n, err := h.w.Write(d)
	h.result = d
	h.err = err
	return n, err
}

// https://github.com/julienschmidt/httprouter
// https://mp.weixin.qq.com/s/9P1AV6d_Cc4pH9DNJeEHHg
