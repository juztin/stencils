// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencils

import (
	"errors"
	"fmt"
	"html/template"
	"io"

	"net/http"
)

type Reader interface {
	Read(*http.Request) ([]byte, error)
}

type Saver interface {
	Save(*http.Request, []byte) error
}

type ReaderSaver interface {
	Reader
	Saver
}

type store struct {
	Reader
	Saver
}

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

func NewStencil(name string, rs ReaderSaver) *Stencil {
	return &Stencil{
		name:    name,
		store:   rs,
		isStale: true,
		fm:      make(map[string]interface{}),
	}
}

func fmtErr(format string, fields ...interface{}) error {
	return errors.New(fmt.Sprintf(format, fields...))
}

func (s *Stencil) Etch(w io.Writer, r *http.Request, data interface{}) error {
	if s.isStale {
		if _, err := s.Load(r); err != nil {
			return err
		}
	}
	return s.Template.Execute(w, data)
}

func (s *Stencil) Extend(t *Stencil) (*Stencil, error) {
	for _, c := range t.children {
		if c.name == s.name {
			return s, fmtErr("%s is already exteded by %s", t.name, s.name)
		}
	}
	s.base = t
	s.isStale = true
	t.children = append(t.children, s)
	return s, nil
}

func (s *Stencil) Load(r *http.Request) (*Stencil, error) {
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

func (s *Stencil) Reload(r *http.Request) (*Stencil, error) {
	if s.base != nil {
		return s.base.Reload(r)
	}
	return s.Load(r)
}

func (s *Stencil) String() string {
	return s.text
}

func (s *Stencil) Funcs(fm template.FuncMap) *Stencil {
	for name, fn := range fm {
		s.fm[name] = fn
	}
	return s
}

func Reload(r *http.Request, children ...*Stencil) (err error) {
	for _, s := range children {
		if _, e := s.Load(r); e != nil {
			err = fmtErr("%s; [%s] - %s", err, s.name, e)
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
