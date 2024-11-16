package logs

import (
	"io"
	"os"
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
	// 拦截错误
	if r := recover(); r != nil {
		CodeError.Print(r)
	}
	// 储存并关闭日志
	confing := []ConfingLog{}
	Std.LogMap.Range(func(key, value any) bool {
		if info, ok := value.(*LogInfo); ok {
			// 是否储存配置到数据库
			if info.IsConfingDB {
				confing = append(confing, info.ConfingLog)
			}
			info.Close()
		}
		return true
	})
	if len(confing) != 0 {
		err := SetConfing("LogConfing", &confing)
		Panicf("储存日志失败:%s ", err != nil, err)
	}
	// 结束程序
	os.Exit(0)
}

// IsExit 程序是否结束
func IsExit() <-chan struct{} {
	return Std.Exit
}

func init() {
	// os.Stdout 转发
	go func() {
		r, w, err := os.Pipe()
		Panicf("日志 Pipe 创建失败: %s", err != nil, err)
		os.Stdout = w
		for err != io.EOF {
			_, err = io.Copy(CodePrint.Info(), r)
		}
	}()
	// os.Stderr 转发
	go func() {
		r, w, err := os.Pipe()
		Panicf("日志 Pipe 创建失败: %s", err != nil, err)
		os.Stderr = w
		for err != io.EOF {
			_, err = io.Copy(CodeError.Info(), r)
		}
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
