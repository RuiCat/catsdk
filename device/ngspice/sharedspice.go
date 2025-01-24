package ngspice

/*
#include <stdbool.h>
#include <stdlib.h>
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
	"unsafe"
)

// 结构体
type ngcomplex_t = C.ngcomplex_t
type pvector_info = *C.vector_info
type pvecvalues = *C.vecvalues
type pvecvaluesall = *C.vecvaluesall
type pvecinfo = *C.vecinfo
type pvecinfoall = *C.vecinfoall

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
	value := (*Ngspice)(userdata).NgspiceValue
	if value.SendChar != nil {
		return C.int(value.SendChar(C.GoString(outputreturn), int(ident), value))
	}
	return -1
}

//export eSendStat
func eSendStat(outputreturn *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.SendStat != nil {
		return C.int(value.SendStat(C.GoString(outputreturn), int(ident), value))
	}
	return -1
}

//export eControlledExit
func eControlledExit(exitstatus C.int, immediate C._Bool, quitexit C._Bool, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.ControlledExit != nil {
		return C.int(value.ControlledExit(int(exitstatus), bool(immediate), bool(quitexit), int(ident), value))
	}
	return -1
}

//export eSendData
func eSendData(vdata pvecvaluesall, numvecs C.int, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	value.ValueMutex.Lock()
	defer value.ValueMutex.Unlock()
	val := value.VecValuesAll
	val.VecCount = int(vdata.veccount)
	val.VecIndex = int(vdata.vecindex)
	vecsa := unsafe.Slice(vdata.vecsa, val.VecCount)
	for k, v := range value.ValueMap {
		vecs := val.VecsA[v]
		vecs.Name = k
		vecs.CReal = float64(vecsa[v].creal)
		vecs.CImag = float64(vecsa[v].cimag)
		vecs.IsScale = bool(vecsa[v].is_scale)
		vecs.IsComplex = bool(vecsa[v].is_complex)

	}
	for i := range len(value.ValueMap) {
		vecs := val.VecsA[i]
		value.Value = append(value.Value, complex(vecs.CReal, vecs.CImag))
	}
	if value.SendData != nil {
		return C.int(value.SendData(val, int(numvecs), int(ident), value))
	}
	return -1
}

//export eSendInitData
func eSendInitData(intdata pvecinfoall, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.SendInitData != nil {
		value.ValueMutex.Lock()
		defer value.ValueMutex.Unlock()
		val := &VecInfoAll{}
		// 初始化参数
		val.Name = C.GoString(intdata.name)
		val.Title = C.GoString(intdata.title)
		val.Date = C.GoString(intdata.date)
		val.Type = C.GoString(intdata._type)
		val.VecCount = int(intdata.veccount)
		val.Vecs = make([]*VecInfo, val.VecCount)
		// 获取全局连接信息
		for i, v := range unsafe.Slice(intdata.vecs, val.VecCount) {
			val.Vecs[i] = &VecInfo{
				Number:     int(v.number),
				VecName:    C.GoString(v.vecname),
				IsReal:     bool(v.is_real),
				PDVec:      v.pdvec,
				PDVecScale: v.pdvecscale,
			}
			// 记录值索引
			value.ValueMap[val.Vecs[i].VecName] = i
		}
		value.VecInfoAll = val
		value.VecValuesAll = &VecValuesAll{
			VecsA: make([]*VecValues, len(value.ValueMap)),
		}
		for i := range value.VecValuesAll.VecsA {
			value.VecValuesAll.VecsA[i] = &VecValues{}
		}
		return C.int(value.SendInitData(val, int(ident), value))
	}
	return -1
}

//export eBGThreadRunning
func eBGThreadRunning(noruns C._Bool, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.BGThreadRunning != nil {
		return C.int(value.BGThreadRunning(bool(noruns), int(ident), value))
	}
	return -1
}

//export eGetVSRCData
func eGetVSRCData(retvoltval *C.double, acttime C.double, nodename *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.GetVSRCData != nil {
		return C.int(value.GetVSRCData((*float64)(retvoltval), float64(acttime), C.GoString(nodename), int(ident), value))
	}
	return -1
}

//export eGetISRCData
func eGetISRCData(retcurrval *C.double, acttime C.double, nodename *C.char, ident C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.GetISRCData != nil {
		return C.int(value.GetISRCData((*float64)(retcurrval), float64(acttime), C.GoString(nodename), int(ident), value))
	}
	return -1
}

//export eGetSyncData
func eGetSyncData(acttime C.double, deltatime *C.double, olddeltatime C.double, redostep C.int, ident C.int, location C.int, userdata unsafe.Pointer) C.int {
	value := (*Ngspice)(userdata).NgspiceValue
	if value.GetSyncData != nil {
		return C.int(value.GetSyncData(float64(acttime), (*float64)(deltatime), float64(olddeltatime), int(redostep), int(ident), int(location), value))
	}
	return 0
}
