package ngspice

import (
	"unsafe"
)
import "C"

/* Dvec flags. */
const (
	Vf_real      int16 = 1 << iota /* The data is real. */
	Vf_complex                     /* The data is complex. */
	Vf_accum                       /* writedata should save this vector. */
	Vf_plot                        /* writedata should incrementally plot it. */
	Vf_print                       /* writedata should print this vector. */
	Vf_mingiven                    /* The v_minsignal value is valid. */
	Vf_maxgiven                    /* The v_maxsignal value is valid. */
	Vf_permanent                   /* Don't garbage collect this vector. */
)

/* Plot types. */
const (
	Plot_lin   = 1
	Plot_comb  = 2
	Plot_point = 3
)

type VectorInfo struct {
	VName     string
	VType     int
	VFlags    int16
	VRealData []float64
	VCompData []complex128
	VLength   int
	_         [4]byte
}

func (vecInfo *VectorInfo) storeVectorInfo(pvectori pvector_info) {
	if pvectori != nil {
		vecInfo.VName = C.GoString(pvectori.v_name)
		vecInfo.VType = int(pvectori.v_type)
		vecInfo.VFlags = int16(pvectori.v_flags)
		vecInfo.VLength = int(pvectori.v_length)
		if vecInfo.VFlags&Vf_real != 0 {
			vecInfo.VRealData = unsafe.Slice((*float64)(pvectori.v_realdata), vecInfo.VLength)
		}
		if vecInfo.VFlags&Vf_complex != 0 {
			vecInfo.VCompData = make([]complex128, vecInfo.VLength)
			for i, v := range unsafe.Slice(pvectori.v_compdata, vecInfo.VLength) {
				vecInfo.VCompData[i] = complex(v.cx_real, v.cx_imag)
			}
		}
	}
}

type VecInfo struct {
	Number     int
	VecName    string
	IsReal     bool
	PDVec      unsafe.Pointer
	PDVecScale unsafe.Pointer
}

type VecInfoAll struct {
	Name     string
	Title    string
	Date     string
	Type     string
	VecCount int
	Vecs     []*VecInfo
}

type VecValues struct {
	Name      string
	CReal     float64
	CImag     float64
	IsScale   bool
	IsComplex bool
	_         [6]byte
}

type VecValuesAll struct {
	VecCount int
	VecIndex int
	VecsA    []*VecValues
}
