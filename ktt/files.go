package ktt

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

type Files map[string][]byte

func (f Files) Merge(files Files) {
	for k, v := range files {
		f[k] = v
	}
}

func LoadFilesDir(dir string) (Files, error) {
	r := make(Files, 64)
	files, err := listFiles(dir)
	if err != nil {
		return nil, err
	}
	for _, i := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, i))
		if err != nil {
			return nil, err
		}
		r[i] = b
	}
	return r, nil
}

func (f Files) String() string {
	r := make([]string, 0, len(f))
	for k := range f {
		r = append(r, k)
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}
