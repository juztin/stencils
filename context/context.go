// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"net/http"

	"bitbucket.org/juztin/stencils"
	"bitbucket.org/juztin/stencils/file"
)

var Stencils *stencils.Stencils

type Context struct {
	*http.Request
	Resp     http.ResponseWriter
	stencils *stencils.Stencils
}

func init() {
	Stencils = stencils.New(file.New("./stencils/"))
	Stencils.FourOh = Stencils.New("errors/404.html")
	Stencils.FiveOh = Stencils.New("errors/500.html")
	Stencils.New("base.html")
}

func New(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{r, w, Stencils}
}

func Reload(n string) func() {
	return func() {
		if o, ok := Stencils.Name(n); ok {
			o.Reload(nil)
		}
	}
}

func NotFound(ctx *Context) {
	ctx.NotFound(nil)
}

func (c *Context) Etch(name string, data interface{}) {
	// log
	c.stencils.Etch(name, c.Resp, c.Request, data)
}

func (c *Context) ServerError(data interface{}) {
	c.Resp.WriteHeader(http.StatusInternalServerError)
	c.stencils.FiveOh.Etch(c.Resp, c.Request, nil)
}

func (c *Context) NotFound(data interface{}) {
	c.Resp.WriteHeader(http.StatusNotFound)
	c.stencils.FourOh.Etch(c.Resp, c.Request, nil)
}
