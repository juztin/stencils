// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package file

import (
	"io/ioutil"
	"path/filepath"

	"bitbucket.org/juztin/stencils"
)

type file struct {
	path string
	name string
}

func New(path string) stencils.StencilFn {
	return func(name string) *stencils.Stencil {
		p := filepath.Join(path, name)
		return stencils.NewStencil(name, &file{p, name})
	}
}

func (f *file) Read(r stencils.Requestor) ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *file) Save(r stencils.Requestor, data []byte) error {
	return ioutil.WriteFile(f.path, data, 0600)
}
