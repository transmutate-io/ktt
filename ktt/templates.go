package ktt

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"text/template"
)

type Templates map[string]*template.Template

func loadTemplatesDir(dir string) (Templates, error) {
	files, err := listFiles(dir)
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Op == "lstat" && e.Err == syscall.ENOENT {
			return nil, ErrNoTemplatesDir
		}
		return nil, err
	}
	tplFiles := make([]string, 0, len(files))
	for i := 0; i < len(files); i++ {
		ff := strings.ToLower(files[i])
		if strings.HasSuffix(ff, ".yaml") || strings.HasSuffix(ff, ".yml") {
			tplFiles = append(tplFiles, files[i])
			files = append(files[:i], files[i+1:]...)
			i--
		}
	}
	for i := range files {
		files[i] = filepath.Join(dir, files[0])
	}
	r := make(Templates, len(tplFiles))
	for _, i := range tplFiles {
		tpl := template.New(i)
		tpl = tpl.Funcs(template.FuncMap{
			"include": func(name string, data interface{}) string {
				out := bytes.NewBuffer(make([]byte, 0, 1024))
				if err := tpl.ExecuteTemplate(out, name, data); err != nil {
					return err.Error()
				}
				return out.String()
			},
		}).Funcs(FuncMap())
		if len(files) > 0 {
			tpl, err = tpl.ParseFiles(files...)
		}
		if err != nil {
			return nil, err
		}
		if tpl, err = tpl.ParseFiles(filepath.Join(dir, i)); err != nil {
			return nil, err
		}
		r[i] = tpl
	}
	return r, nil
}

func (t Templates) String() string {
	r := make([]string, 0, len(t))
	for k := range t {
		r = append(r, k)
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}
