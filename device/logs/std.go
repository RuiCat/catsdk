package logs

import (
	"fmt"
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
	// 退出程序
	select {
	case <-Std.Exit:
		return
	default:
		// 关闭通道
		close(Std.Exit)
	}
	if r := recover(); r != nil {
		ErrorLevel.Info().WriteString(fmt.Sprint(r))
	}
	// 储存并关闭日志
	confing := []ConfingLog{}
	Std.LogMap.Range(func(key, value any) bool {
		if info, ok := value.(*LogInfo); ok {
			if info.IsDB {
				confing = append(confing, info.ConfingLog)
			}
			info.Close()
		}
		return true
	})
	if len(confing) != 0 {
		Panicf("储存日志失败:%s ", SetConfing("LogConfing", &confing))
	}
	// 结束程序
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
		Exit()
	}()
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
