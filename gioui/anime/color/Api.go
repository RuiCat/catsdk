package color

// BlendType 混合类型
type BlendType uint

// ColorBlend 像素混合接口
type ColorBlend func(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64)

// 混合类型常量
const (
	BlendTypeAdd BlendType = iota
	BlendTypeClear
	BlendTypeDarken
	BlendTypeDst
	BlendTypeDstAtop
	BlendTypeDstIn
	BlendTypeDstOut
	BlendTypeDstOver
	BlendTypeLighten
	BlendTypeMultiply
	BlendTypeOverlay
	BlendTypeScreen
	BlendTypeSrc
	BlendTypeSrcAtop
	BlendTypeSrcIn
	BlendTypeSrcOut
	BlendTypeSrcOver
	BlendTypeXor
	BlendTypeMix
)

// Blend 绘制
var Blend = map[BlendType]ColorBlend{
	BlendTypeAdd:      Add,
	BlendTypeClear:    Clear,
	BlendTypeDarken:   Darken,
	BlendTypeDst:      Dst,
	BlendTypeDstAtop:  DstAtop,
	BlendTypeDstIn:    DstIn,
	BlendTypeDstOut:   DstOut,
	BlendTypeDstOver:  DstOver,
	BlendTypeLighten:  Lighten,
	BlendTypeMultiply: Multiply,
	BlendTypeOverlay:  Overlay,
	BlendTypeScreen:   Screen,
	BlendTypeSrc:      Src,
	BlendTypeSrcAtop:  SrcAtop,
	BlendTypeSrcIn:    SrcIn,
	BlendTypeSrcOut:   SrcOut,
	BlendTypeSrcOver:  SrcOver,
	BlendTypeXor:      Xor,
	BlendTypeMix:      Mix,
}

// BlendStr 名称
var BlendStr = map[BlendType]string{
	BlendTypeAdd:      "Add",
	BlendTypeClear:    "Clear",
	BlendTypeDarken:   "Darken",
	BlendTypeDst:      "Dst",
	BlendTypeDstAtop:  "DstAtop",
	BlendTypeDstIn:    "DstIn",
	BlendTypeDstOut:   "DstOut",
	BlendTypeDstOver:  "DstOver",
	BlendTypeLighten:  "Lighten",
	BlendTypeMultiply: "Multiply",
	BlendTypeOverlay:  "Overlay",
	BlendTypeScreen:   "Screen",
	BlendTypeSrc:      "Src",
	BlendTypeSrcAtop:  "SrcAtop",
	BlendTypeSrcIn:    "SrcIn",
	BlendTypeSrcOut:   "SrcOut",
	BlendTypeSrcOver:  "SrcOver",
	BlendTypeXor:      "Xor",
	BlendTypeMix:      "Mix",
}
