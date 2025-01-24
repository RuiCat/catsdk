package ngspice

import (
	"device/gousb/dlopen"
	"sync"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

var ngspicePath string

type symbol[T any] struct {
	name string
	ptr  unsafe.Pointer
	Func T
}

func (symbol *symbol[T]) GetPointer(handle dlopen.Handle) {
	var err error
	symbol.ptr, err = handle.GetSymbolPointer(symbol.name)
	if err != nil {
		panic(err)
	}
	symbol.Func = dlopen.GetFn[T](symbol.ptr, 0)

}

type NgspiceValue struct {
	SendChar        func(string, int, *NgspiceValue) int
	SendStat        func(string, int, *NgspiceValue) int
	ControlledExit  func(int, bool, bool, int, *NgspiceValue) int
	SendData        func(*VecValuesAll, int, int, *NgspiceValue) int
	SendInitData    func(*VecInfoAll, int, *NgspiceValue) int
	BGThreadRunning func(bool, int, *NgspiceValue) int
	GetVSRCData     func(*float64, float64, string, int, *NgspiceValue) int
	GetISRCData     func(*float64, float64, string, int, *NgspiceValue) int
	GetSyncData     func(float64, *float64, float64, int, int, int, *NgspiceValue) int
	Value           []complex128
	ValueMap        map[string]int
	ValueMutex      sync.Mutex
	VecInfoAll      *VecInfoAll
	VecValuesAll    *VecValuesAll
}

type Ngspice struct {
	*NgspiceValue
	id           int
	handle       dlopen.Handle
	ngGetVecInfo *symbol[ngGet_Vec_Info]
	ngInit       *symbol[ngSpice_Init]
	ngInitSync   *symbol[ngSpice_Init_Sync]
	ngCommand    *symbol[ngSpice_Command]
	ngCirc       *symbol[ngSpice_Circ]
	ngCurPlot    *symbol[ngSpice_CurPlot]
	ngAllPlots   *symbol[ngSpice_AllPlots]
	ngAllVecs    *symbol[ngSpice_AllVecs]
	ngrunning    *symbol[ngSpice_running]
	ngSetBkpt    *symbol[ngSpice_SetBkpt]
}

func NewNgspice(id int) *Ngspice {
	ng := &Ngspice{
		id:           id,
		ngGetVecInfo: &symbol[ngGet_Vec_Info]{name: "ngGet_Vec_Info"},
		ngInit:       &symbol[ngSpice_Init]{name: "ngSpice_Init"},
		ngInitSync:   &symbol[ngSpice_Init_Sync]{name: "ngSpice_Init_Sync"},
		ngCommand:    &symbol[ngSpice_Command]{name: "ngSpice_Command"},
		ngCirc:       &symbol[ngSpice_Circ]{name: "ngSpice_Circ"},
		ngCurPlot:    &symbol[ngSpice_CurPlot]{name: "ngSpice_CurPlot"},
		ngAllPlots:   &symbol[ngSpice_AllPlots]{name: "ngSpice_AllPlots"},
		ngAllVecs:    &symbol[ngSpice_AllVecs]{name: "ngSpice_AllVecs"},
		ngrunning:    &symbol[ngSpice_running]{name: "ngSpice_running"},
		ngSetBkpt:    &symbol[ngSpice_SetBkpt]{name: "ngSpice_SetBkpt"},
		NgspiceValue: &NgspiceValue{
			Value:    []complex128{},
			ValueMap: map[string]int{},
		},
	}
	var err error
	ng.handle, err = dlopen.GetHandle("/lib/libngspice.so")
	if err != nil {
		panic(err)
	}
	ng.ngGetVecInfo.GetPointer(ng.handle)
	ng.ngInit.GetPointer(ng.handle)
	ng.ngInitSync.GetPointer(ng.handle)
	ng.ngCommand.GetPointer(ng.handle)
	ng.ngrunning.GetPointer(ng.handle)
	ng.ngCirc.GetPointer(ng.handle)
	ng.ngCurPlot.GetPointer(ng.handle)
	ng.ngAllPlots.GetPointer(ng.handle)
	ng.ngAllVecs.GetPointer(ng.handle)
	ng.ngSetBkpt.GetPointer(ng.handle)
	ng.ngInit.Func(callSendChar, callSendStat, callControlledExit, callSendData, callSendInitData, callBGThreadRunning, unsafe.Pointer(ng))
	var ident C.int = C.int(ng.id)
	ng.ngInitSync.Func(callGetVSRCData, callGetISRCData, callGetSyncData, &ident, unsafe.Pointer(ng))
	return ng
}
func (ng *Ngspice) Command(command string) int {
	c := C.CString(command)
	defer C.free(unsafe.Pointer(c))
	return int(ng.ngCommand.Func(c))
}
func (ng *Ngspice) Circ(circa []string) int {
	n := len(circa)
	cStrings := make([]*C.char, n+1)
	for i := range circa {
		cStrings[i] = C.CString(circa[i])
		defer C.free(unsafe.Pointer(cStrings[i]))
	}
	cStrings[n] = nil
	return int(ng.ngCirc.Func((**C.char)(unsafe.Pointer(&cStrings[0]))))
}
func (ng *Ngspice) Quit() {
	ng.Command("quit")
	ng.handle.Close()
	(*ng) = Ngspice{}
}
func (ng *Ngspice) Clear() int {
	return int(ng.ngCirc.Func(nil))
}
func (ng *Ngspice) Running() bool {
	return bool(ng.ngrunning.Func())
}
func (ng *Ngspice) GetVecInfo(name string) (ret *VectorInfo) {
	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))
	ret = &VectorInfo{}
	if pvec := ng.ngGetVecInfo.Func(c); pvec != nil {
		ret.storeVectorInfo(pvec)
	}
	return ret
}
func (ng *Ngspice) CurPlot() string {
	var ret string
	if pchar := ng.ngCurPlot.Func(); pchar != nil {
		ret = C.GoString(pchar)
	}
	return ret
}
func (ng *Ngspice) AllPlots() []string {
	var ret []string
	if ppchar := ng.ngAllPlots.Func(); ppchar != nil {
		ret = ppcharToStringSlice(ppchar)
	}
	return ret
}
func (ng *Ngspice) AllVecs(s string) []string {
	var ret []string
	c := C.CString(s)
	defer C.free(unsafe.Pointer(c))
	if ppchar := ng.ngAllVecs.Func(c); ppchar != nil {
		ret = ppcharToStringSlice(ppchar)
	}
	return ret
}
func ppcharToStringSlice(ppchar **C.char) []string {
	var goStrings []string
	for i := 0; ; i++ {
		cString := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(ppchar)) + uintptr(i)*unsafe.Sizeof(ppchar)))
		if cString == nil {
			break
		}
		goStrings = append(goStrings, C.GoString(cString))
	}
	return goStrings
}
