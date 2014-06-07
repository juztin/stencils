// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package inline

import (
	"io/ioutil"
	"net/http"

	"code.minty.io/stencils"
)

type inline struct {
	name string
}

func New(path string) stencils.StencilFn {
	return func(name string) *stencils.Stencil {
		return stencils.NewStencil(name, &inline{name})
	}
}

func (f *inline) Read(r *http.Request) ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *inline) Save(r *http.Request, data []byte) error {
	return nil
}
