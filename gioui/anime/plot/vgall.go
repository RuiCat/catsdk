// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !minimal
// +build !minimal

package plot

import (
	_ "gioui/anime/plot/vg/vgeps"
	_ "gioui/anime/plot/vg/vgimg"
	_ "gioui/anime/plot/vg/vgpdf"
	_ "gioui/anime/plot/vg/vgsvg"
	_ "gioui/anime/plot/vg/vgtex"
)
