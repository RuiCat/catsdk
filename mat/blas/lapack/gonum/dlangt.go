// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"mat/blas"
	"math"
)

// Dlangt returns the value of the given norm of an n×n tridiagonal matrix
// represented by the three diagonals.
//
// d must have length at least n and dl and du must have length at least n-1.
func (impl Implementation) Dlangt(norm blas.MatrixNorm, n int, dl, d, du []float64) float64 {
	switch {
	case norm != blas.MaxAbs && norm != blas.MaxRowSum && norm != blas.MaxColumnSum && norm != blas.Frobenius:
		panic(badNorm)
	case n < 0:
		panic(nLT0)
	}

	if n == 0 {
		return 0
	}

	switch {
	case len(dl) < n-1:
		panic(shortDL)
	case len(d) < n:
		panic(shortD)
	case len(du) < n-1:
		panic(shortDU)
	}

	dl = dl[:n-1]
	d = d[:n]
	du = du[:n-1]

	var anorm float64
	switch norm {
	case blas.MaxAbs:
		for _, diag := range [][]float64{dl, d, du} {
			for _, di := range diag {
				if math.IsNaN(di) {
					return di
				}
				di = math.Abs(di)
				if di > anorm {
					anorm = di
				}
			}
		}
	case blas.MaxColumnSum:
		if n == 1 {
			return math.Abs(d[0])
		}
		anorm = math.Abs(d[0]) + math.Abs(dl[0])
		if math.IsNaN(anorm) {
			return anorm
		}
		tmp := math.Abs(du[n-2]) + math.Abs(d[n-1])
		if math.IsNaN(tmp) {
			return tmp
		}
		if tmp > anorm {
			anorm = tmp
		}
		for i := 1; i < n-1; i++ {
			tmp = math.Abs(du[i-1]) + math.Abs(d[i]) + math.Abs(dl[i])
			if math.IsNaN(tmp) {
				return tmp
			}
			if tmp > anorm {
				anorm = tmp
			}
		}
	case blas.MaxRowSum:
		if n == 1 {
			return math.Abs(d[0])
		}
		anorm = math.Abs(d[0]) + math.Abs(du[0])
		if math.IsNaN(anorm) {
			return anorm
		}
		tmp := math.Abs(dl[n-2]) + math.Abs(d[n-1])
		if math.IsNaN(tmp) {
			return tmp
		}
		if tmp > anorm {
			anorm = tmp
		}
		for i := 1; i < n-1; i++ {
			tmp = math.Abs(dl[i-1]) + math.Abs(d[i]) + math.Abs(du[i])
			if math.IsNaN(tmp) {
				return tmp
			}
			if tmp > anorm {
				anorm = tmp
			}
		}
	case blas.Frobenius:
		scale := 0.0
		ssq := 1.0
		scale, ssq = impl.Dlassq(n, d, 1, scale, ssq)
		if n > 1 {
			scale, ssq = impl.Dlassq(n-1, dl, 1, scale, ssq)
			scale, ssq = impl.Dlassq(n-1, du, 1, scale, ssq)
		}
		anorm = scale * math.Sqrt(ssq)
	}
	return anorm
}
