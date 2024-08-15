// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !minimal
// +build !minimal

package plot

import (
	_ "mat/plot/vg/vgeps"
	_ "mat/plot/vg/vgimg"
	_ "mat/plot/vg/vgpdf"
	_ "mat/plot/vg/vgsvg"
	_ "mat/plot/vg/vgtex"
)
