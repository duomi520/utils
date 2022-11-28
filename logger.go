package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"
)

//ILogger 日志接口
type ILogger interface {
	SetLevel(int)
	Debug(a ...string)
	Info(a ...string)
	Warn(a ...string)
	Error(a ...string)
	Fatal(a ...string)
	Close()
}

// levels
const (
	UnknownLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	CrashLevel
)

//ErrWLoggerCrash 日志崩溃
var ErrWLoggerCrash = errors.New("Logger.Write：crash")

//ErrWLoggerUndefinedLevel 未定义等级
var ErrWLoggerUndefinedLevel = errors.New("Logger.NewWLogger：undefined Logger level")

//LevelName 等级名称
var LevelName = []string{"[?????]", "[Debug]", "[Info ]", "[Warn ]", "[Error]", "[Fatal]", "[Crash]"}

type fileWrite struct {
	fd              *os.File
	path            string
	LogFileSize     int //单个日志文件大小
	MaxLogFileCount int //最大日志文件数
	count           int //计数
}

//Write 写
func (f *fileWrite) Write(b []byte) (n int, err error) {
	size, err := f.fd.Write(b)
	if err == nil {
		f.count = f.count + size
		if f.count > f.LogFileSize {
			f.count = 0
			//关闭现有文件
			if err := f.fd.Close(); err != nil {
				log.Printf("fileWrite.Write：close faile, %s\n", err.Error())
				return size, ErrWLoggerCrash
			}
			//新建文件
			if err := f.createFile(); err != nil {
				log.Printf("fileWrite.Write：create file faile, %s\n", err.Error())
				return size, ErrWLoggerCrash
			}
			//遍历文件数
			lm := make([]int, 0, 2*f.MaxLogFileCount)
			//遍历目录，读出日志文件名
			filepath.Walk(f.path, func(path string, fi os.FileInfo, err error) error {
				if fi == nil {
					return err
				}
				if fi.IsDir() {
					return nil
				}
				name := fi.Name()
				c, err := strconv.Atoi(name[:len(name)-4])
				if err != nil {
					log.Printf("fileWrite.Write：name format error, %s, %s\n", name, err.Error())
				}
				lm = append(lm, c)
				return nil
			})
			//删除超出的日志文件
			if len(lm) > f.MaxLogFileCount {
				sort.Ints(lm)
				for _, v := range lm[:len(lm)-f.MaxLogFileCount] {
					name := fmt.Sprintf("./%s/%d.log", f.path, v)
					err := os.Remove(name)
					if err != nil {
						log.Printf("fileWrite.Write：remove failed, %s, %s\n", name, err.Error())
					}
				}
			}
		}
	}
	return size, err
}

//Close 关闭
func (f *fileWrite) Close() error {
	if f.fd != nil {
		return f.fd.Close()
	}
	return nil
}
func (f *fileWrite) createFile() error {
	now := time.Now()
	filename := fmt.Sprintf("%d%02d%02d%02d%02d%02d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
	var err error
	f.fd, err = os.Create(path.Join(f.path, filename))
	return err
}

//WLogger 日志
type WLogger struct {
	level      int //输出等级
	bufferPool sync.Pool
	mu         sync.Mutex
	out        io.WriteCloser
}

//NewWLogger 新日志
func NewWLogger(level int, logPath string) (*WLogger, error) {
	if level > FatalLevel || level <= UnknownLevel {
		return nil, ErrWLoggerUndefinedLevel
	}
	wlogger := &WLogger{
		level: level,
	}
	wlogger.bufferPool.New = func() any {
		return &bytes.Buffer{}
	}
	if len(logPath) != 0 {
		fw := &fileWrite{count: 0, path: logPath, LogFileSize: 10 * 1024 * 1024,
			MaxLogFileCount: 10}
		if err := fw.createFile(); err != nil {
			return nil, err
		}
		wlogger.out = fw
	} else {
		wlogger.out = os.Stdout
	}

	return wlogger, nil
}

//SetLogFileSize 设置log文件尺寸
func (logger *WLogger) SetLogFileSize(LogFileSize, MaxLogFileCount int) {
	_, ok := logger.out.(*fileWrite)
	if ok {
		logger.out.(*fileWrite).LogFileSize = LogFileSize
		logger.out.(*fileWrite).MaxLogFileCount = MaxLogFileCount
	}
}

//SetLevel 设置等级
func (logger *WLogger) SetLevel(l int) {
	logger.level = l
}

//Close 关闭
func (logger *WLogger) Close() {
	if err := logger.out.Close(); err != nil {
		log.Printf("WLogger.Close: %s\n", err.Error())
	}
}

//doPrint 输出
func (logger *WLogger) doPrint(level int, a ...string) {
	if level < logger.level {
		return
	}
	b := logger.bufferPool.Get().(*bytes.Buffer)
	b.WriteString(LevelName[level])
	b.WriteString(time.Now().Format(" 2006-01-02 15:04:05 "))
	for _, v := range a {
		b.WriteString(v)
	}
	b.WriteString("\n")
	logger.mu.Lock()
	if _, err := b.WriteTo(logger.out); err != nil {
		//日志崩溃,禁止写入
		logger.level = CrashLevel
		log.Printf("WLogger.doPrint：WriteTo faile, %s\n", err.Error())
	}
	logger.mu.Unlock()
	b.Reset()
	logger.bufferPool.Put(b)
	if level == FatalLevel {
		logger.Close()
		os.Exit(1)
	}
}

//Debug Debug
func (logger *WLogger) Debug(a ...string) {
	logger.doPrint(DebugLevel, a...)
}

//Info Info
func (logger *WLogger) Info(a ...string) {
	logger.doPrint(InfoLevel, a...)
}

//Warn Warn
func (logger *WLogger) Warn(a ...string) {
	logger.doPrint(WarnLevel, a...)
}

//Error Error
func (logger *WLogger) Error(a ...string) {
	logger.doPrint(ErrorLevel, a...)
}

//Fatal Fatal
func (logger *WLogger) Fatal(a ...string) {
	logger.doPrint(FatalLevel, a...)
}

// https://github.com/siddontang/go-log/blob/master/log/log.go
