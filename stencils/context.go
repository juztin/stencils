// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencils

import (
	"net/http"

	"code.minty.io/stencils"
	"code.minty.io/stencils/file"
)

var s *stencils.Stencils = stencils.New(file.New("./templates"))

func Reload(n string) func() {
	return func() {
		if o, ok := s.Name(n); ok {
			o.Reload(nil)
		}
	}
}

func New(name string) *stencils.Stencil {
	return s.New(name)
}

func Name(name string) (*stencils.Stencil, bool) {
	return s.Name(name)
}

func Add(st *stencils.Stencil) *stencils.Stencils {
	return s.Add(st)
}

func Remove(name string) (ok bool) {
	return s.Remove(name)
}

func Etch(name string, data interface{}, w http.ResponseWriter, r *http.Request) {
	s.Etch(name, w, r, data)
}

func ServerError(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	s.FiveOh.Etch(w, r, nil)
}

func NotFound(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.WriteHeader(http.StatusNotFound)
	s.FiveOh.Etch(w, r, nil)
}
