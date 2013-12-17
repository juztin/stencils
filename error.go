// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencils

type Error struct {
	error
	Status int
}

func NewError(status int, err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{err, status}
}
