// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stencil

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
	store    ReaderSaver
}

func New(name string, rs ReaderSaver) *Stencil {
	s := new(Stencil)
	s.name = name
	s.store = rs
	s.isStale = true
	return s
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

	b, err := s.store.Read(r)
	if err != nil {
		return s, err
	}

	t := template.New(s.name)
	// Load base template
	if s.base != nil {
		if t, err = t.Parse(s.base.Root.String()); err != nil {
			return s, err
		}
	}
	// Load this template
	if t, err = t.Parse(string(b)); err != nil {
		return s, err
	}
	s.Template = t
	s.isStale = false

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
