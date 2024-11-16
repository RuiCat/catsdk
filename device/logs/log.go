package logs

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

// LogLevel 日志等级
type LogLevel uint8

// 默认日志等级
const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	UnknownLevel
)

var (
	CodePrint = InfoLevel.NewCode(0)
	CodeError = ErrorLevel.NewCode(1)
)

// ConfingLog 日志配置
type ConfingLog struct {
	Name        string   // 定义名称
	Level       LogLevel // 日志等级
	IsConfingDB bool     // 是否储存配置到数据库
	IsNet       bool     // 启用远程日志
	IsFile      bool     // 启用本地文件
	NetAddr     net.Addr // 远程日志
	FileName    string   // 文件路径
}

// New 初始化
func (confing ConfingLog) New() {
	if _, ok := Std.LogMap.Load(confing.Level); !ok {
		info := &LogInfo{ConfingLog: confing, LogWriter: &Std.LogOut[0]}
		lw := []LogWriter{}
		// 将日志写入文件
		if info.IsFile {
			f, err := os.OpenFile(info.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			Panicf("无法创建日志文件: %s", err != nil, err)
			lw = append(lw, f)
		}
		// 将日志写入远程
		if info.IsNet {
			Panic("目前功能没有实现")
		}
		// 处理日志转发
		if len(lw) != 0 {
			r, w, err := os.Pipe()
			Panicf("日志 Pipe 创建失败: %s", err != nil, err)
			info.LogWriter = w
			info.close = func() {
				w.Close()
			}
			go func() {
				scanner := bufio.NewScanner(r)
				for scanner.Scan() {
					b := scanner.Bytes()
					if b[len(b)-1] != '\n' {
						b = append(b, '\n')
					}
					// 转发日志信息
					for _, w := range lw {
						w.Write(b)
					}
				}
				for _, w := range lw {
					w.Close()
				}
			}()
		}
		Std.LogMap.Store(confing.Level, info)
	}
}

// LogWriter 日志写入接口
type LogWriter interface {
	io.Writer
	io.StringWriter
	io.Closer
	Name() string
}

// LogInfo 日志信息
type LogInfo struct {
	ConfingLog        // 日志配置
	LogWriter         // 日志写入接口
	close      func() // 用于关闭
}

// Close 关闭连接
func (Info *LogInfo) Close() error {
	if Info.close != nil {
		Info.close()
	}
	return nil
}

func init() {
	// 初始化配置
	confing := []ConfingLog{
		{Name: "DebugLevel", IsConfingDB: true, Level: DebugLevel, FileName: "Debug.log"},
		{Name: "InfoLevel", IsConfingDB: true, Level: InfoLevel, FileName: "Info.log"},
		{Name: "WarnLevel", IsConfingDB: true, Level: WarnLevel, FileName: "Warn.log"},
		{Name: "ErrorLevel", IsConfingDB: true, Level: ErrorLevel, FileName: "Error.log"},
		{Name: "FatalLevel", IsConfingDB: true, Level: FatalLevel, FileName: "Fatal.log"},
		{Name: "UnknownLevel", IsConfingDB: true, Level: UnknownLevel, FileName: "Unknown.log"},
	}
	err := GetConfing("LogConfing", &confing)
	Panicf("默认日志初始化失败:%s ", err != nil, err)
	// 创建日志
	for _, cfg := range confing {
		cfg.New()
	}
}

// NewCode 创建新日志代码
func (Level LogLevel) NewCode(Code LogCode) LogCode {
	if _, ok := Std.LogCode.Load(Code); !ok {
		Std.LogCode.Store(Code, Level)
		return Code
	}
	panic(fmt.Sprintf("函数 NewCode(Code:%d,Level:%d) 重复创建", Code, Level))
}

// Info 日志信息
func (level LogLevel) Info() *LogInfo {
	if info, ok := Std.LogMap.Load(level); ok {
		return info.(*LogInfo)
	}
	return nil
}

// String 等级标志
func (level LogLevel) String() string {
	if info, ok := Std.LogMap.Load(level); ok {
		return (info.(*LogInfo)).ConfingLog.Name
	}
	return fmt.Sprintf("LogLevel(%d?)", level)
}

// LogCode 日志代码
type LogCode uint64

// Level 日志等级
func (code LogCode) Level() LogLevel {
	if level, ok := Std.LogCode.Load(code); ok {
		return level.(LogLevel)
	}
	return LogLevel(UnknownLevel)
}

// Info 日志信息
func (code LogCode) Info() *LogInfo {
	return code.Level().Info()
}

// String 代码标志
func (code LogCode) String() string {
	return fmt.Sprintf("Code:%d Level:%s", code, code.Level().String())
}

// Print 写入信息
func (code LogCode) Print(a ...any) {
	code.Info().LogWriter.WriteString(fmt.Sprintf("%s,Details:%v\n", code.String(), a))
}

// Error 输出错误
func (code LogCode) Error(err error) {
	code.Print(err.Error())
}

// IsError 检测错误
func (code LogCode) IsError(err error, info ...string) {
	if err != nil {
		code.Print(fmt.Sprintf("%s,info:%v\n", err.Error(), info))
	}
}

// Defer 拦截错误
func (code LogCode) Defer() {
	if err := recover(); err != nil {
		code.Print(err)
	}
}

// Details 构建错误
func (code LogCode) Details(a any, Info ...any) error {
	return &LogDetails{Code: code, Details: a, Info: Info}
}

// LogDetails 日志记录
type LogDetails struct {
	Code    LogCode // 日志代码
	Details any     // 日志内容
	Info    []any   // 附加信息
}

func (log LogDetails) Error() string {
	buffer := Std.Pool.Get()
	defer buffer.Free()
	buffer.WriteString("Code:[")
	buffer.WriteTypeString(log.Code)
run:
	if len(log.Info) > 0 {
		buffer.WriteString(",Info:[")
		for _, v := range log.Info {
			buffer.WriteTypeString(v)
			buffer.WriteString(",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("]")
	}
	switch v := log.Details.(type) {
	case *LogDetails:
		log = *v
	case LogDetails:
		log = v
	default:
		goto exit
	}
	buffer.WriteString("->")
	buffer.WriteTypeString(log.Code)
	goto run
exit:
	buffer.WriteString("],Details:")
	buffer.WriteTypeString(log.Details)
	return buffer.String()
}

// IfPrint 输出错误/日志信息
func IfPrint(is bool, err error) {
	if is && err != nil {
		CodePrint.Print(err)
	}
}

// Print 输出错误/日志信息
func Print(err any) {
	if err != nil {
		CodePrint.Print(err)
	}
}

// Printf 输出错误/日志信息
func Printf(format string, is bool, v ...any) {
	if is {
		CodePrint.Print(fmt.Sprintf(format, v...))
	}
}

// Panic 输出错误/日志信息
// 遇到无法修复错误退出程序执行
func Panic(err any) {
	if err != nil {
		panic(err)
	}
}

// Panicf 输出错误/日志信息
// 遇到无法修复错误退出程序执行
func Panicf(format string, is bool, v ...any) {
	if is {
		panic(fmt.Errorf(format, v...))
	}
}

// Recover 拦截错误
func Recover() {
	if err := recover(); err != nil {
		CodeError.Print(err)
	}
}

// Recovere 拦截错误
func Recovere(err *error) {
	if e := recover(); e != nil {
		(*err) = fmt.Errorf("%s", e)
	}
}

// Recoverf 拦截错误
func Recoverf(format string) {
	if e := recover(); e != nil {
		CodeError.Print(fmt.Errorf(format, e))
	}
}

// Recoverfe 拦截错误
func Recoverfe(format string, err *error) {
	if e := recover(); e != nil {
		(*err) = fmt.Errorf(format, e)
	}
}
