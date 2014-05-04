// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencils

import (
	"fmt"
	"html/template"
	"io"

	"net/http"
)

// Reader reads stencil data.
type Reader interface {
	Read(*http.Request) ([]byte, error)
}

// Saver saves stencil data.
type Saver interface {
	Save(*http.Request, []byte) error
}

// ReaderSaver reads and saves stencil data.
type ReaderSaver interface {
	Reader
	Saver
}

type store struct {
	Reader
	Saver
}

// Stencil is a wrapper around html/templatea.
type Stencil struct {
	*template.Template
	base     *Stencil
	children []*Stencil
	isStale  bool
	name     string
	text     string
	store    ReaderSaver
	fm       template.FuncMap
}

// NewStencil returns a new stencil for the given name, using the given reader.
func NewStencil(name string, rs ReaderSaver) *Stencil {
	return &Stencil{
		name:    name,
		store:   rs,
		isStale: true,
		fm:      make(map[string]interface{}),
	}
}


// Etch writes the stencil to the given writer.
func (s *Stencil) Etch(w io.Writer, r *http.Request, data interface{}) error {
	if s.isStale {
		if _, err := s.Load(r); err != nil {
			return err
		}
	}
	return s.Template.Execute(w, data)
}

// Extend extends the given parent stencil.
func (s *Stencil) Extend(parent *Stencil) (*Stencil, error) {
	for _, c := range parent.children {
		if c.name == s.name {
			return s, fmt.Errorf("%s is already exteded by %s", parent.name, s.name)
		}
	}
	s.base = parent
	s.isStale = true
	parent.children = append(parent.children, s)
	return s, nil
}

// Load loads the stencil from it's reader.
func (s *Stencil) Load(r *http.Request) (*Stencil, error) {
	// If there is a base, reload it (will reload this as part of it's children)
	if s.base != nil && s.base.isStale {
		return s.base.Load(r)
	}

	// Get template text
	b, err := s.store.Read(r)
	if err != nil {
		return s, err
	}

	// Store text, create template
	s.text = string(b)
	t := template.New(s.name)

	// Add template funcs
	if len(s.fm) > 0 {
		t.Funcs(s.fm)
	}

	// Load full templates (traverse parents)
	if t, err = t.Parse(data(s)); err != nil {
		return s, err
	}
	s.Template = t
	s.isStale = false

	// Load children
	return s, Reload(r, s.children...)
}

// Reload reloads this stencil, and notifies all base templates.
func (s *Stencil) Reload(r *http.Request) (*Stencil, error) {
	if s.base != nil {
		return s.base.Reload(r)
	}
	return s.Load(r)
}

// String returns the templates data.
func (s *Stencil) String() string {
	return s.text
}

// Funcs adds a function to the internal html.Template funcs.
func (s *Stencil) Funcs(fm template.FuncMap) *Stencil {
	for name, fn := range fm {
		s.fm[name] = fn
	}
	return s
}

// Reload reloads all of the given child stencils.
func Reload(r *http.Request, children ...*Stencil) (err error) {
	for _, s := range children {
		if _, e := s.Load(r); e != nil {
			err = fmt.Errorf("%s; [%s] - %s", err, s.name, e)
		}
	}
	return
}

func data(s *Stencil) string {
	if s.base != nil {
		return data(s.base) + s.text
	}
	return s.text
}
