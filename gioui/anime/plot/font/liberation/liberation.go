// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package liberation exports the Liberation fonts as a font.Collection.
package liberation // import "gioui/anime/plot/font/liberation"

import (
	"fmt"
	"sync"

	stdfnt "sdk/x/image/font"
	"sdk/x/image/font/opentype"

	"sdk/liberation/liberationmonobold"
	"sdk/liberation/liberationmonobolditalic"
	"sdk/liberation/liberationmonoitalic"
	"sdk/liberation/liberationmonoregular"
	"sdk/liberation/liberationsansbold"
	"sdk/liberation/liberationsansbolditalic"
	"sdk/liberation/liberationsansitalic"
	"sdk/liberation/liberationsansregular"
	"sdk/liberation/liberationserifbold"
	"sdk/liberation/liberationserifbolditalic"
	"sdk/liberation/liberationserifitalic"
	"sdk/liberation/liberationserifregular"

	"gioui/anime/plot/font"
)

var (
	once       sync.Once
	collection font.Collection
)

func Collection() font.Collection {
	once.Do(func() {
		addColl(font.Font{}, liberationserifregular.TTF)
		addColl(font.Font{Style: stdfnt.StyleItalic}, liberationserifitalic.TTF)
		addColl(font.Font{Weight: stdfnt.WeightBold}, liberationserifbold.TTF)
		addColl(font.Font{
			Style:  stdfnt.StyleItalic,
			Weight: stdfnt.WeightBold,
		}, liberationserifbolditalic.TTF)

		// mono variant
		addColl(font.Font{Variant: "Mono"}, liberationmonoregular.TTF)
		addColl(font.Font{
			Variant: "Mono",
			Style:   stdfnt.StyleItalic,
		}, liberationmonoitalic.TTF)
		addColl(font.Font{
			Variant: "Mono",
			Weight:  stdfnt.WeightBold,
		}, liberationmonobold.TTF)
		addColl(font.Font{
			Variant: "Mono",
			Style:   stdfnt.StyleItalic,
			Weight:  stdfnt.WeightBold,
		}, liberationmonobolditalic.TTF)

		// sans-serif variant
		addColl(font.Font{Variant: "Sans"}, liberationsansregular.TTF)
		addColl(font.Font{
			Variant: "Sans",
			Style:   stdfnt.StyleItalic,
		}, liberationsansitalic.TTF)
		addColl(font.Font{
			Variant: "Sans",
			Weight:  stdfnt.WeightBold,
		}, liberationsansbold.TTF)
		addColl(font.Font{
			Variant: "Sans",
			Style:   stdfnt.StyleItalic,
			Weight:  stdfnt.WeightBold,
		}, liberationsansbolditalic.TTF)

		n := len(collection)
		collection = collection[:n:n]
	})

	return collection
}

func addColl(fnt font.Font, ttf []byte) {
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Errorf("vg: could not parse font: %+v", err))
	}
	fnt.Typeface = "Liberation"
	if fnt.Variant == "" {
		fnt.Variant = "Serif"
	}
	collection = append(collection, font.Face{
		Font: fnt,
		Face: face,
	})
}
