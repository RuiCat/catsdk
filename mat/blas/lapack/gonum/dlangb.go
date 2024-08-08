// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"mat/asm/f64"
	"mat/blas"
	"math"
)

// Dlangb returns the given norm of an m×n band matrix with kl sub-diagonals and
// ku super-diagonals.
func (impl Implementation) Dlangb(norm blas.MatrixNorm, m, n, kl, ku int, ab []float64, ldab int) float64 {
	ncol := kl + 1 + ku
	switch {
	case norm != blas.MaxAbs && norm != blas.MaxRowSum && norm != blas.MaxColumnSum && norm != blas.Frobenius:
		panic(badNorm)
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case kl < 0:
		panic(klLT0)
	case ku < 0:
		panic(kuLT0)
	case ldab < ncol:
		panic(badLdA)
	}

	// Quick return if possible.
	if m == 0 || n == 0 {
		return 0
	}

	switch {
	case len(ab) < min(m, n+kl)*ldab:
		panic(shortAB)
	}

	var value float64
	switch norm {
	case blas.MaxAbs:
		for i := 0; i < min(m, n+kl); i++ {
			l := max(0, kl-i)
			u := min(n+kl-i, ncol)
			for _, aij := range ab[i*ldab+l : i*ldab+u] {
				aij = math.Abs(aij)
				if aij > value || math.IsNaN(aij) {
					value = aij
				}
			}
		}
	case blas.MaxRowSum:
		for i := 0; i < min(m, n+kl); i++ {
			l := max(0, kl-i)
			u := min(n+kl-i, ncol)
			sum := f64.L1Norm(ab[i*ldab+l : i*ldab+u])
			if sum > value || math.IsNaN(sum) {
				value = sum
			}
		}
	case blas.MaxColumnSum:
		for j := 0; j < min(m+ku, n); j++ {
			jb := min(kl+j, ncol-1)
			ib := max(0, j-ku)
			jlen := min(j+kl, m-1) - ib + 1
			sum := f64.L1NormInc(ab[ib*ldab+jb:], jlen, max(1, ldab-1))
			if sum > value || math.IsNaN(sum) {
				value = sum
			}
		}
	case blas.Frobenius:
		scale := 0.0
		sum := 1.0
		for i := 0; i < min(m, n+kl); i++ {
			l := max(0, kl-i)
			u := min(n+kl-i, ncol)
			ilen := u - l
			scale, sum = impl.Dlassq(ilen, ab[i*ldab+l:], 1, scale, sum)
		}
		value = scale * math.Sqrt(sum)
	}
	return value
}
