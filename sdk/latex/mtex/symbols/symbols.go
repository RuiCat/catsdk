// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package symbols contains logic about TeX symbols.
package symbols // import "sdk/latex/mtex/symbols"

var (
	SpacedSymbols = UnionOf(BinaryOperators, RelationSymbols, ArrowSymbols)
)

func IsSpaced(s string) bool {
	return SpacedSymbols.Has(s)
}
