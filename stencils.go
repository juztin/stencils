// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencil

import (
	"bytes"
	"io"
)

type Stencils struct {
	col    map[string]*Stencil
	create StencilFn
	FourOh *Stencil
	FiveOh *Stencil
}

type StencilFn func(name string) *Stencil

func NewStencils(fn StencilFn) *Stencils {
	c := new(Stencils)
	c.col = make(map[string]*Stencil)
	c.create = fn
	return c
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

func (c *Stencils) Etch(n string, r Requestor, wr io.Writer, data interface{}) error {
	if s, ok := c.Name(n); ok {
		var buf bytes.Buffer
		if err := s.Etch(r, &buf, data); err != nil {
			return writeErr(c.FiveOh, r, wr, data, err)
		}
		_, err := buf.WriteTo(wr)
		return err
	}
	return writeErr(c.FourOh, r, wr, data, fmtErr("Stencils does not exist: %s", n))
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

func writeErr(s *Stencil, r Requestor, wr io.Writer, data interface{}, err error) error {
	if s != nil {
		if e := s.Etch(r, wr, data); e != nil {
			err = fmtErr("Stencil: %s : %s", err, e)
		}
	}
	return err
}
