// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencils

import (
	"errors"
	"fmt"
	"io"
	"text/template"

	"net/http"
)

type Requestor interface {
	Request() *http.Request
}

type Reader interface {
	Read(Requestor) ([]byte, error)
}

type Saver interface {
	Save(Requestor, []byte) error
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
}

func NewStencil(name string, rs ReaderSaver) *Stencil {
	return &Stencil{
		name:    name,
		store:   rs,
		isStale: true,
	}
}

func fmtErr(format string, fields ...interface{}) error {
	return errors.New(fmt.Sprintf(format, fields...))
}

func (s *Stencil) Etch(r Requestor, wr io.Writer, data interface{}) error {
	if s.isStale {
		if _, err := s.Load(r); err != nil {
			return err
		}
	}
	return s.Template.Execute(wr, data)
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

func (s *Stencil) Load(r Requestor) (*Stencil, error) {
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

	// Load full templates (traverse parents)
	if t, err = t.Parse(data(s)); err != nil {
		return s, err
	}
	s.Template = t
	s.isStale = false

	// Load children
	return s, Reload(r, s.children...)
}

func Reload(r Requestor, children ...*Stencil) (err error) {
	for _, s := range children {
		if _, e := s.Load(r); e != nil {
			err = fmtErr("%s; [%s] - %s", err, s.name, e)
		}
	}
	return
}

func data(s *Stencil) string {
	if s.base != nil {
		return s.text + data(s.base)
	}
	return s.text
}
