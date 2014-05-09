// +build appengine

// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gae

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.minty.io/stencils"

	"appengine"
	"appengine/datastore"
)

type gae struct {
	collection, name string
}

func New(collection string) stencils.StencilFn {
	return func(name string) *stencils.Stencil {
		return stencils.NewStencil(name, &gae{collection, name})
	}
}

func (g *gae) Read(r *http.Request) ([]byte, error) {
	return ioutil.ReadFile(g.name)
}

func (g *gae) Save(r *http.Request, data []byte) error {
	c := appengine.NewContext(req)
	k := datastore.NewKey(c, g.collection, g.name, 0, nil)

	if _, err := datastore.Put(c, k, b); err != nil {
		return errors.New(fmt.Sprintf("Failed to save %s, %s", g.name, err))
	}
}
