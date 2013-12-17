// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencils

import (
	"bytes"
	"io"
	"net/http"
)

type Stencils struct {
	col    map[string]*Stencil
	create StencilFn
	FourOh *Stencil
	FiveOh *Stencil
}

func (e Error) Error() string {
	return e.error.Error()
}

type StencilFn func(name string) *Stencil

func New(fn StencilFn) *Stencils {
	c := new(Stencils)
	c.col = make(map[string]*Stencil)
	c.create = fn
	return c
}

func etchErr(status int, wr io.Writer, r *http.Request, err error, s *Stencil) *Error {
	Log(err)

	if r, ok := wr.(http.ResponseWriter); ok {
		r.WriteHeader(status)
	}

	if s != nil {
		// Todo, concat errors
		s.Etch(wr, r, data)
	}
	return NewError(status, err)
}

func (c *Stencils) New(name string) *Stencil {
	s := c.create(name)
	c.Add(s)
	return s
}

func (c *Stencils) Add(s *Stencil) *Stencils {
	c.col[s.name] = s
	return c
}

func (c *Stencils) Etch(n string, wr io.Writer, r *http.Request, data interface{}) *Error {
	if s, ok := c.Name(n); ok {
		var buf bytes.Buffer
		if err := s.Etch(&buf, r, data); err != nil {
			return etchErr(500, wr, r, err, c.FiveOh)
		} else if _, err = buf.WriteTo(wr); err != nil {
			return etchErr(500, wr, r, err, c.FiveOh)
		}
		return nil
	}
	return etchErr(404, wr, r, fmtErr("Stencil does not exist: %s", n), c.FourOh)
}

func (c *Stencils) Name(n string) (*Stencil, bool) {
	s, ok := c.col[n]
	return s, ok
}

func (c *Stencils) Remove(n string) (ok bool) {
	if _, ok = c.col[n]; ok {
		delete(c.col, n)
	}
	return ok
}
