// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package file

import (
	"io/ioutil"
	"net/http"
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

/*func new(path string, name string) *stencils.Stencil {
	p := filepath.Join(path, name)
	return stencils.NewStencil(name, &file{p, name})
}*/

func (f *file) Read(r *http.Request) ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *file) Save(r *http.Request, data []byte) error {
	return ioutil.WriteFile(f.path, data, 0600)
}

/*func LoadAll(path string) *stencils.Stencils {
	col := stencils.New(New(path))
	loadAll(col, path, nil)
	return col
}

func loadAll(col *stencils.Stencils, path string, base *stencils.Stencil) {
	f, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, fi := range f {
		if fi.IsDir() {
			loadAll(col, fi.Name(), base)
		}
		t := new(path, fi.Name())
		col.Add(t)
		if base != nil {
			base.Extend(t)
		}
	}

}*/
