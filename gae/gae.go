// +build appengine

// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gae

import (
	"errors"
	"fmt"

	"appengine"
	"appengine/datastore"
	"bitbucket.org/juztin/stencil"
)

type gae struct {
	collection, name string
}

func New(collection string) stencil.StencilFn {
	return func(name string) *stencil.Stencil {
		return stencil.New(name, &gae{collection, name})
	}
}

func (g *gae) Read(r stencil.Requestor) ([]byte, error) {
	return ioutil.ReadFile(g.name)
}

func (g *gae) Save(r stencil.Requestor, data []byte) error {
	c := appengine.NewContext(req)
	k := datastore.NewKey(c, g.collection, g.name, 0, nil)

	if _, err := datastore.Put(c, k, b); err != nil {
		return errors.New(fmt.Sprintf("Failed to save %s, %s", g.name, err))
	}
}
