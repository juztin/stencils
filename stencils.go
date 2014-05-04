// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stencils is a simple wrapper around html/template adding template inheritance, and other helpers.
package stencils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Stencils is a collection of stencil objects.
type Stencils struct {
	col    map[string]*Stencil
	create StencilFn
	FourOh *Stencil
	FiveOh *Stencil
	Logger *log.Logger
}

// StencilFn is a func type that returns a new stencil.
type StencilFn func(name string) *Stencil

// New Returns a Stencils collection.
func New(fn StencilFn) *Stencils {
	c := new(Stencils)
	c.col = make(map[string]*Stencil)
	c.create = fn
	//c.LogTo(os.Stderr)
	c.Logger = log.New(os.Stderr, "[stencils] ", log.LstdFlags)
	return c
}

// New adds a new stencil, by name, to the collection.
func (c *Stencils) New(name string) *Stencil {
	s := c.create(name)
	c.Add(s)
	return s
}

// Add adds the given stencil to the collection.
func (c *Stencils) Add(s *Stencil) *Stencils {
	c.col[s.name] = s
	return c
}

// Etch writes the stencil to the given writer.
func (c *Stencils) Etch(n string, wr io.Writer, r *http.Request, data interface{}) *Error {
	if s, ok := c.Name(n); ok {
		var buf bytes.Buffer
		if err := s.Etch(&buf, r, data); err != nil {
			return c.etchErr(500, wr, r, err, c.FiveOh)
		} else if _, err = buf.WriteTo(wr); err != nil {
			return c.etchErr(500, wr, r, err, c.FiveOh)
		}
		return nil
	}
	return c.etchErr(404, wr, r, fmt.Errorf("Stencil does not exist: %s", n), c.FourOh)
}

// Name returns a stencil by name.
func (c *Stencils) Name(name string) (*Stencil, bool) {
	s, ok := c.col[name]
	return s, ok
}

// Remove removes a stencil from the collection by name.
func (c *Stencils) Remove(name string) (ok bool) {
	if _, ok = c.col[name]; ok {
		delete(c.col, name)
	}
	return ok
}

/*func (c *Stencils) LogTo(w io.Writer) {
	if w == nil {
		c.logger = nil
	} else {
		c.logger = log.New(w, "[stencils] ", log.LstdFlags)
	}
}*/

/*func (c *Stencils) log(err error) {
	if err != nil && c.Logger != nil {
		c.Logger.Println(err.Error())
	}
}*/

func (c *Stencils) etchErr(status int, wr io.Writer, r *http.Request, err error, s *Stencil) *Error {
	//c.log(err)

	/*if w, ok := wr.(http.ResponseWriter); ok {
		w.WriteHeader(status)
	}*/

	if s != nil {
		// Todo, concat errors
		s.Etch(wr, r, data)
	}
	return NewError(status, err)
}
