package color

func Add(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	saturate := func(a, b float64) float64 {
		if a+b > 1 {
			return 1
		}
		return a + b
	}
	return saturate(ScR, DcR),
		saturate(ScG, DcG),
		saturate(ScB, DcB),
		saturate(Sa, Da)
}

func Clear(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return 0, 0, 0, 0
}

func Darken(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	min := func(a, b float64) float64 {
		if a > b {
			return b
		}
		return a
	}
	da := 1 - Da
	sa := 1 - Sa
	return ScR*da + DcR*sa + min(ScR, DcR),
		ScG*da + DcG*sa + min(ScG, DcG),
		ScB*da + DcB*sa + min(ScB, DcB),
		Sa + Da - Sa*Da
}

func Dst(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return DcR, DcG, DcB, Da
}

func DstAtop(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	da := 1 - Da
	return Sa*DcR + ScR*da,
		Sa*DcG + ScG*da,
		Sa*DcB + ScB*da,
		Sa
}
func DstIn(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return Sa * DcR,
		Sa * DcG,
		Sa * DcB,
		Sa * Da
}

func DstOut(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	sa := 1 - Sa
	return DcR * sa,
		DcG * sa,
		DcB * sa,
		Da * sa
}

func DstOver(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	sa := 1 - Sa
	da := 1 - Da
	return DcR + da*ScR,
		DcG + da*ScG,
		DcB + da*ScB,
		Sa + sa*Da
}

func Lighten(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	max := func(a, b float64) float64 {
		if a < b {
			return b
		}
		return a
	}
	da := 1 - Da
	sa := 1 - Sa
	return ScR*da + DcR*sa + max(ScR, DcR),
		ScG*da + DcG*sa + max(ScG, DcG),
		ScB*da + DcB*sa + max(ScB, DcB),
		Sa + Da - Sa*Da
}

func Multiply(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return ScR * DcR,
		ScG * DcG,
		ScB * DcB,
		Sa * Da
}

func Overlay(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	Ra = Sa + Da - Sa*Da
	a := Sa * Da
	out := func(s, d float64) float64 {
		if 2*d < Da {
			return 2 * s * d
		}
		return a - 2*(Da-s)*(Sa-d)
	}
	Rr = out(ScR, DcR)
	Rg = out(ScG, DcG)
	Rb = out(ScB, DcB)
	return Rr, Rg, Rb, Ra
}

func Screen(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return ScR + DcR - ScR*DcR,
		ScG + DcG - ScG*DcG,
		ScB + DcB - ScB*DcB,
		Sa + Da - Sa*Da
}

func Src(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return ScR, ScG, ScB, Sa
}

func SrcAtop(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	sa := 1 - Sa
	return ScR*Da + sa*DcR,
		ScG*Da + sa*DcG,
		ScB*Da + sa*DcB,
		Da
}
func SrcIn(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	return ScR * Da,
		ScG * Da,
		ScB * Da,
		Sa * Da
}

func SrcOut(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	da := 1 - Da
	return ScR * da,
		ScG * da,
		ScB * da,
		Sa * da
}

func SrcOver(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	sa := 1 - Sa
	return ScR + sa*DcR,
		ScG + sa*DcG,
		ScB + sa*DcB,
		Sa + sa*Da
}

func Xor(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	sa := 1 - Sa
	da := 1 - Da
	return ScR*da + sa*DcR,
		ScG*da + sa*DcG,
		ScB*da + sa*DcB,
		Sa + Da - 2*Sa*Da
}

func Mix(ScR, ScG, ScB, Sa, DcR, DcG, DcB, Da float64) (Rr, Rg, Rb, Ra float64) {
	a := (1 - Sa)
	R := ScR*Sa + DcR*Da*a
	G := ScG*Sa + DcG*Da*a
	B := ScB*Sa + DcB*Da*a
	Alpha := 1 - a*(1-Da)
	R = R / Alpha
	G = G / Alpha
	B = B / Alpha
	return R, G, B, Alpha
}
