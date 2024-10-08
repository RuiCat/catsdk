// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//go:build !appengine
// +build !appengine

package appengine

import (
	"device/util/net/context"

	"device/util/appengine/internal"
)

// BackgroundContext returns a context not associated with a request.
// This should only be used when not servicing a request.
// This only works in App Engine "flexible environment".
func BackgroundContext() context.Context {
	return internal.BackgroundContext()
}
