// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"mat/blas"
	"mat/blas/blas64"
	"math"
)

// Dlanhs returns the value of the one norm, or the Frobenius norm, or the
// infinity norm, or the element of largest absolute value of a Hessenberg
// matrix A.
//
// If norm is blas.MaxColumnSum, work must have length at least n.
func (impl Implementation) Dlanhs(norm blas.MatrixNorm, n int, a []float64, lda int, work []float64) float64 {
	switch {
	case norm != blas.MaxRowSum && norm != blas.MaxAbs && norm != blas.MaxColumnSum && norm != blas.Frobenius:
		panic(badNorm)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	if n == 0 {
		return 0
	}

	switch {
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case norm == blas.MaxColumnSum && len(work) < n:
		panic(shortWork)
	}

	bi := blas64.Implementation()
	var value float64
	switch norm {
	case blas.MaxAbs:
		for i := 0; i < n; i++ {
			minj := max(0, i-1)
			for _, v := range a[i*lda+minj : i*lda+n] {
				value = math.Max(value, math.Abs(v))
			}
		}
	case blas.MaxColumnSum:
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		for i := 0; i < n; i++ {
			for j := max(0, i-1); j < n; j++ {
				work[j] += math.Abs(a[i*lda+j])
			}
		}
		for _, v := range work[:n] {
			value = math.Max(value, v)
		}
	case blas.MaxRowSum:
		for i := 0; i < n; i++ {
			minj := max(0, i-1)
			sum := bi.Dasum(n-minj, a[i*lda+minj:], 1)
			value = math.Max(value, sum)
		}
	case blas.Frobenius:
		scale := 0.0
		sum := 1.0
		for i := 0; i < n; i++ {
			minj := max(0, i-1)
			scale, sum = impl.Dlassq(n-minj, a[i*lda+minj:], 1, scale, sum)
		}
		value = scale * math.Sqrt(sum)
	}
	return value
}
