package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func TestStack(t *testing.T) {
	err1 := errors.New("a error")
	err2 := WrapStack(1, err1)
	fmt.Println(err2)
	fmt.Println(ErrorWithStack(err2))
	//wrap
	err3 := fmt.Errorf("warp %w", err2)
	fmt.Println(err3)
	//unwarp
	err4 := errors.Unwrap(err3)
	fmt.Println(ErrorWithStack(err4))
	//nil
	fmt.Println(ErrorWithStack(nil))
	err5 := WrapStack(1, nil)
	fmt.Println(err5)
}

/*
a error
[errors_test.go:13] a error
warp a error
[errors_test.go:13] a error

<nil>
*/

func TestReplaceAttr(t *testing.T) {
	err1 := errors.New("something")
	err2 := WrapStack(1, err1)
	err3 := fmt.Errorf("warp1 %w", err2)
	err4 := WrapStack(1, err3)
	err5 := fmt.Errorf("warp2 %w", err4)
	err6 := fmt.Errorf("warp3 %w", err5)
	err7 := WrapStack(1, err6)
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: ReplaceAttr,
	})
	fmt.Println(ErrorWithFullStack(err7))
	logger := slog.New(h)
	logger.Error("message", err7)
}

/*
[errors_test.go:44] warp3 warp2 warp1 something
warp2 warp1 something
[errors_test.go:41] warp1 something
[errors_test.go:39] somethin
time=2024-04-30T00:00:12.360+08:00 level=ERROR msg=message trace="[errors_test.go:44] warp3 warp2 warp1 something\nwarp2 warp1 something\n[errors_test.go:41] warp1 something\n[errors_test.go:39] somethin"
*/
