package itabs

import (
	"reflect"
	"runtime"
	"unsafe"
)

var FtabMap = map[string]*runtime.Func{}
var RtypeMap = map[any]reflect.Type{}

// -ldflags="-s"
func init() {
	pc := make([]uintptr, 255)
	frames := runtime.CallersFrames([]uintptr{pc[runtime.Callers(0, pc)-1]})
	frame, ok := frames.Next()
	if !ok {
		funcInfo := reflect.ValueOf(frame).FieldByName("funcInfo")
		datap := funcInfo.FieldByName("datap").Elem()
		ftab := datap.FieldByName("ftab")
		pclntable := datap.FieldByName("pclntable").Pointer()
		for i := range ftab.Len() {
			funcPtr := (*runtime.Func)(unsafe.Pointer(pclntable + uintptr(ftab.Index(i).FieldByName("funcoff").Uint())))
			FtabMap[funcPtr.Name()] = funcPtr
		}
		// 数据底层类型
		rtype := reflect.TypeOf(reflect.TypeOf(0)).Elem()
		// 获取类型
		types := uintptr(datap.FieldByName("types").Uint())
		typelinks := datap.FieldByName("typelinks")
		for i := range typelinks.Len() {
			typeAt := reflect.NewAt(rtype, unsafe.Pointer(types+uintptr(typelinks.Index(i).Int())))
			RtypeMap[typeAt.Interface()] = typeAt.Type()
		}
	}
}

func GetFunc[T any](name string) *T {
	ftab, ok := FtabMap[name]
	if ok {
		entry := ftab.Entry()
		ptr := &entry
		return (*T)(unsafe.Pointer(&ptr))
	}
	return nil
}
