package ngspice

import "unsafe"
import "C"

/* Dvec flags. */
const (
	vf_real      = (1 << 0) /* The data is real. */
	vf_complex   = (1 << 1) /* The data is complex. */
	vf_accum     = (1 << 2) /* writedata should save this vector. */
	vf_plot      = (1 << 3) /* writedata should incrementally plot it. */
	vf_print     = (1 << 4) /* writedata should print this vector. */
	vf_mingiven  = (1 << 5) /* The v_minsignal value is valid. */
	vf_maxgiven  = (1 << 6) /* The v_maxsignal value is valid. */
	vf_permanent = (1 << 7) /* Don't garbage collect this vector. */
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
		if vecInfo.VFlags&vf_real != 0 {
			vecInfo.VRealData = unsafe.Slice((*float64)(pvectori.v_realdata), vecInfo.VLength)
		}
		if vecInfo.VFlags&vf_complex != 0 {
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
