package logs

import (
	"os"
	"os/signal"
	"sync"
)

// Std 全局函数
var Std = struct {
	LogOut  [2]os.File
	LogMap  sync.Map
	LogCode sync.Map
	Exit    chan struct{}
	Pool    Pool
}{Exit: make(chan struct{}), LogOut: [2]os.File{*os.Stdout, *os.Stderr}, Pool: NewPool(1024)}

// Exit 退出程序
func Exit() {
	close(Std.Exit)
	os.Exit(0)
}

// IsExit 程序是否结束
func IsExit() <-chan struct{} {
	return Std.Exit
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go func() {
		<-c
		close(Std.Exit)
		os.Exit(0)
	}()
	// 更改默认输出
	os.Stdout = InfoLevel.Info().LogFile
	os.Stderr = ErrorLevel.Info().LogFile
}

// Close 用于处理代码中存在多次关联的 Close 过程
func Close(close *func()) (func(), *bool, func(fn any)) {
	var i int
	var ok bool
	add := make([]any, 0, 10)
	(*close) = func() {
		close = nil
		for ; i < len(add); i++ {
			switch fn := add[i].(type) {
			case func():
				fn()
			case func() error:
				fn()
			}
		}
	}
	return func() {
			if !ok {
				if close != nil {
					(*close)()
				}
			}
		}, &ok, func(fn any) {
			add = append(add, fn)
		}
}
