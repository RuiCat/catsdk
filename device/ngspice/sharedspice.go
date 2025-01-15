package ngspice

/*
#include <stdbool.h>
#include "sharedspice.h"
extern int eSendChar(char*, int, void*);
extern int eSendStat(char*, int, void*);
extern int eControlledExit(int, bool, bool, int, void*);
extern int eSendData(pvecvaluesall, int, int, void*);
extern int eSendInitData(pvecinfoall, int, void*);
extern int eBGThreadRunning(bool, int, void*);
extern int eGetVSRCData(double*, double, char*, int, void*);
extern int eGetISRCData(double*, double, char*, int, void*);
extern int eGetSyncData(double, double*, double, int, int, int, void*);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// 结构体
type ngcomplex_t = C.ngcomplex_t
type pvector_info = C.pvector_info
type pvecvalues = C.pvecvalues
type pvecvaluesall = C.pvecvaluesall
type pvecinfo = C.pvecinfo
type pvecinfoall = C.pvecinfoall

// 回调函数类型
type sendChar = func(*C.char, C.int, unsafe.Pointer) C.int
type sendStat = func(*C.char, C.int, unsafe.Pointer) C.int
type controlledExit = func(C.int, C._Bool, C._Bool, C.int, unsafe.Pointer) C.int
type sendData = func(unsafe.Pointer, C.int, C.int, unsafe.Pointer) C.int
type sendInitData = func(unsafe.Pointer, C.int, unsafe.Pointer) C.int
type bGThreadRunning = func(C._Bool, C.int, unsafe.Pointer) C.int
type getVSRCData = func(*C.double, C.double, *C.char, C.int, unsafe.Pointer) C.int
type getISRCData = func(*C.double, C.double, *C.char, C.int, unsafe.Pointer) C.int
type getSyncData = func(C.double, *C.double, C.double, C.int, C.int, C.int, unsafe.Pointer) C.int

// 导出函数定义
type ngSpice_Init = func(printfcn *sendChar, statfcn *sendStat, ngexit *controlledExit, sdata *sendData, sinitdata *sendInitData, bgtrun *bGThreadRunning, userData unsafe.Pointer) C.int
type ngSpice_Init_Sync = func(vsrcdat *getVSRCData, isrcdat *getISRCData, syncdat *getSyncData, ident *C.int, userData unsafe.Pointer) C.int
type ngSpice_Command = func(command *C.char) C.int
type ngGet_Vec_Info = func(vecname *C.char) pvector_info
type ngSpice_Circ = func(circarray **C.char) C.int
type ngSpice_CurPlot = func() *C.char
type ngSpice_AllPlots = func() **C.char
type ngSpice_AllVecs = func(plotname *C.char) **C.char
type ngSpice_running = func() C._Bool
type ngSpice_SetBkpt = func(C.double) C._Bool

// 回调函数实现
var callSendChar = (*sendChar)(C.eSendChar)
var callSendStat = (*sendStat)(C.eSendStat)
var callControlledExit = (*controlledExit)(C.eControlledExit)
var callSendData = (*sendData)(C.eSendData)
var callSendInitData = (*sendInitData)(C.eSendInitData)
var callBGThreadRunning = (*bGThreadRunning)(C.eBGThreadRunning)
var callGetVSRCData = (*getVSRCData)(C.eGetVSRCData)
var callGetISRCData = (*getISRCData)(C.eGetISRCData)
var callGetSyncData = (*getSyncData)(C.eGetSyncData)

//export eSendChar
func eSendChar(outputreturn *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("SendChar", C.GoString(outputreturn), ident, userdata)
	return -1
}

//export eSendStat
func eSendStat(outputreturn *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("SendStat", C.GoString(outputreturn), ident, userdata)
	return -1
}

//export eControlledExit
func eControlledExit(exitstatus C.int, immediate C._Bool, quitexit C._Bool, ident C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("ControlledExit", exitstatus, immediate, quitexit, ident, userdata)
	return -1
}

//export eSendData
func eSendData(vdata pvecvaluesall, numvecs C.int, ident C.int, userdata unsafe.Pointer) C.int {
	a := &VecValuesAll{}
	a.VecCount = int(vdata.veccount)
	a.VecIndex = int(vdata.vecindex)
	a.VecsA = unsafe.Slice((**VecValues)(unsafe.Pointer(vdata.vecsa)), a.VecCount)
	fmt.Printf("SendData: %#v %d %d %d\n", *a, numvecs, ident, userdata)
	return -1
}

//export eSendInitData
func eSendInitData(intdata pvecinfoall, ident C.int, userdata unsafe.Pointer) C.int {
	a := &VecInfoAll{}
	a.Name = C.GoString(intdata.name)
	a.Title = C.GoString(intdata.title)
	a.Date = C.GoString(intdata.date)
	a.Type = C.GoString(intdata._type)
	a.VecCount = int(intdata.veccount)
	a.Vecs = unsafe.Slice((**VecInfo)(unsafe.Pointer(intdata.vecs)), a.VecCount)
	fmt.Printf("SendInitData: %#v %d %d\n", *a, ident, userdata)
	return -1
}

//export eBGThreadRunning
func eBGThreadRunning(noruns C._Bool, ident C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("BGThreadRunning", noruns, ident, userdata)
	return -1
}

//export eGetVSRCData
func eGetVSRCData(retvoltval *C.double, acttime C.double, nodename *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("GetVSRCData", retvoltval, acttime, nodename, ident, userdata)
	return -1
}

//export eGetISRCData
func eGetISRCData(retcurrval *C.double, acttime C.double, nodename *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("GetISRCData", retcurrval, acttime, nodename, ident, userdata)
	return -1
}

//export eGetSyncData
func eGetSyncData(acttime C.double, deltatime *C.double, olddeltatime C.double, redostep C.int, ident C.int, location C.int, userdata unsafe.Pointer) C.int {
	fmt.Println("GetSyncData", acttime, *deltatime, olddeltatime, redostep, ident, location, userdata)
	return 0
}
