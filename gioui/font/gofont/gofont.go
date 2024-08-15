// SPDX-License-Identifier: Unlicense OR MIT

// Package gofont exports the Go fonts as a text.Collection.
//
// See https://blog.golang.org/go-fonts for a description of the
// fonts, and the golang.org/x/image/font/gofont packages for the
// font data.
package gofont

import (
	"fmt"
	"sync"

	"sdk/x/image/font/gofont/gobold"
	"sdk/x/image/font/gofont/gobolditalic"
	"sdk/x/image/font/gofont/goitalic"
	"sdk/x/image/font/gofont/gomedium"
	"sdk/x/image/font/gofont/gomediumitalic"
	"sdk/x/image/font/gofont/gomono"
	"sdk/x/image/font/gofont/gomonobold"
	"sdk/x/image/font/gofont/gomonobolditalic"
	"sdk/x/image/font/gofont/gomonoitalic"
	"sdk/x/image/font/gofont/goregular"
	"sdk/x/image/font/gofont/gosmallcaps"
	"sdk/x/image/font/gofont/gosmallcapsitalic"

	"gioui/font"
	"gioui/font/opentype"
)

var (
	regOnce    sync.Once
	reg        []font.FontFace
	once       sync.Once
	collection []font.FontFace
)

func loadRegular() {
	regOnce.Do(func() {
		faces, err := opentype.ParseCollection(goregular.TTF)
		if err != nil {
			panic(fmt.Errorf("failed to parse font: %v", err))
		}
		reg = faces
		collection = append(collection, reg[0])
	})
}

// Regular returns a collection of only the Go regular font face.
func Regular() []font.FontFace {
	loadRegular()
	return reg
}

// Regular returns a collection of all available Go font faces.
func Collection() []font.FontFace {
	loadRegular()
	once.Do(func() {
		register(goitalic.TTF)
		register(gobold.TTF)
		register(gobolditalic.TTF)
		register(gomedium.TTF)
		register(gomediumitalic.TTF)
		register(gomono.TTF)
		register(gomonobold.TTF)
		register(gomonobolditalic.TTF)
		register(gomonoitalic.TTF)
		register(gosmallcaps.TTF)
		register(gosmallcapsitalic.TTF)
		// Ensure that any outside appends will not reuse the backing store.
		n := len(collection)
		collection = collection[:n:n]
	})
	return collection
}

func register(ttf []byte) {
	faces, err := opentype.ParseCollection(ttf)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	collection = append(collection, faces[0])
}
