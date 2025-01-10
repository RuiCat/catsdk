package itabs

import (
	"reflect"
	"runtime"
	"unsafe"
)

var FtabMap = map[string]*runtime.Func{}

func init() {
	pc := make([]uintptr, 255)
	frames := runtime.CallersFrames([]uintptr{pc[runtime.Callers(0, pc)-1]})
	frame, ok := frames.Next()
	if !ok {
		funcInfo := reflect.ValueOf(frame).FieldByName("funcInfo")
		datap := funcInfo.FieldByName("datap")
		ftab := datap.Elem().FieldByName("ftab")
		pclntable := datap.Elem().FieldByName("pclntable").Pointer()
		for i := range ftab.Len() {
			funcPtr := (*runtime.Func)(unsafe.Pointer(pclntable + uintptr(ftab.Index(i).FieldByName("funcoff").Uint())))
			FtabMap[funcPtr.Name()] = funcPtr
		}
	}
}
