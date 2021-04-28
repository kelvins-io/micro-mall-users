package goroutine

import (
	"context"
	"gitee.com/kelvins-io/kelvins"
	"runtime/debug"
)

func AttachPanicHandle(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				kelvins.ErrLogger.Errorf(context.Background(), "goroutine panic: %v, stacktrace:%v", err, string(debug.Stack()))
			}
		}()
		f()
	}
}
