package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

/*
Errors are just values not sentinel errors 防止包引入
handle not just check errors
Only handle errors once
要求传入的字符串首字母小写，结尾不带标点符号
将 errors 看成黑盒，判断它的行为，而不是类型
使用 Wrap Unwrap
最上层打印
*/

type wError struct {
	stack string
	err   error
}

func (m wError) Error() string {
	return m.err.Error()
}

func ErrorWithStack(err error) string {
	if err == nil {
		return ""
	}
	m, ok := err.(wError)
	if ok {
		return m.stack
	}
	return err.Error()

}

func WrapStack(skip int, err error) error {
	if err == nil {
		return nil
	}
	_, f, l, _ := runtime.Caller(skip)
	result := strings.Split(f, "/")
	m := wError{
		stack: fmt.Sprintf("[%s:%d] %s", result[len(result)-1], l, err),
		err:   err,
	}
	return m
}
func ErrorWithFullStack(w error) string {
	v, ok := w.(wError)
	if ok {
		var builder strings.Builder
		builder.WriteString(v.stack)
		builder.WriteString("\n")
		e := errors.Unwrap(v.err)
		for e != nil {
			w, ok := e.(wError)
			if ok {
				e = w.err
				builder.WriteString(w.stack)
				builder.WriteString("\n")
			} else {
				builder.WriteString(e.Error())
				builder.WriteString("\n")
			}
			e = errors.Unwrap(e)
		}
		return builder.String()[:builder.Len()-2]
	}
	return w.Error()
}
func ReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		switch v := a.Value.Any().(type) {
		case wError:
			a = slog.String("trace", ErrorWithFullStack(v))
		}
	}
	return a
}

func FormatRecover() ([]byte, any) {
	if r := recover(); r != nil {
		const size = 65536
		buf := make([]byte, size)
		end := min(runtime.Stack(buf, false), size)
		return buf[:end], r
	}
	return nil, nil
}

// https://zhuanlan.zhihu.com/p/82985617
