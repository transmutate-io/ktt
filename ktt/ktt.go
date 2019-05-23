package ktt

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"
)

type Info struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type TemplateInfo struct {
	App      *Info `yaml:"app"`
	Template *Info `yaml:"template"`
}

type Template struct {
	Info          *TemplateInfo
	Values        Values
	TemplateFiles Templates
	Files         Files
}

const (
	infoFilename   = "info"
	valuesFilename = "values"
	templatesDir   = "templates"
	filesDir       = "files"
)

func LoadTemplate(dir string) (*Template, error) {
	var err error
	r := &Template{}
	if r.Info, err = LoadTemplateInfo(dir); err != nil {
		return nil, err
	}
	var fileOpen bool
	for _, i := range []string{"yaml", "yml"} {
		if r.Values, err = LoadValuesFile(dir, valuesFilename+"."+i); err != nil {
			if e, ok := err.(*os.PathError); ok && e.Op == "open" && e.Err == syscall.ENOENT {
				continue
			}
			return nil, err
		}
		fileOpen = true
		break
	}
	if !fileOpen {
		return nil, ErrNoValuesFile
	}
	if r.Files, err = LoadFilesDir(filepath.Join(dir, filesDir)); err != nil {
		exists := true
		if e, ok := err.(*os.PathError); ok && e.Op == "lstat" && e.Err == syscall.ENOENT {
			exists = false
		}
		if exists {
			return nil, err
		}
	}
	if r.TemplateFiles, err = loadTemplatesDir(filepath.Join(dir, templatesDir)); err != nil {
		return nil, err
	}
	if len(r.TemplateFiles) == 0 {
		return nil, ErrNoTemplates
	}
	return r, nil
}

func loadYAMLFile(fn string, r interface{}) error {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, r)
}

var (
	ErrNoInfoFile       = errors.New("info (info.yaml) file not found")
	ErrAppInfoNotSet    = errors.New("app information not set")
	ErrAppNameNotSet    = errors.New("app name not set")
	ErrAppVersionNotSet = errors.New("app version not set")
	ErrNoValuesFile     = errors.New("values (values.yaml) file not found")
	ErrNoTemplatesDir   = errors.New("templates directory not found")
	ErrNoTemplates      = errors.New("no templates found")
)

func LoadTemplateInfo(dir string) (*TemplateInfo, error) {
	r := &TemplateInfo{}
	var fileOpen bool
	for _, i := range []string{"yaml", "yml"} {
		if err := loadYAMLFile(filepath.Join(dir, infoFilename+"."+i), r); err != nil {
			if e, ok := err.(*os.PathError); ok && e.Op == "open" && e.Err == syscall.ENOENT {
				continue
			}
			return nil, err
		}
		fileOpen = true
		break
	}
	if !fileOpen {
		return nil, ErrNoInfoFile
	}
	if r.App == nil {
		return nil, ErrAppInfoNotSet
	}
	if r.App.Name == "" {
		return nil, ErrAppNameNotSet
	}
	if r.App.Version == "" {
		return nil, ErrAppVersionNotSet
	}
	if r.Template == nil {
		r.Template = r.App
	} else {
		if r.Template.Name == "" {
			r.Template.Name = r.App.Name
		}
		if r.Template.Version == "" {
			r.Template.Version = r.App.Version
		}
	}
	return r, nil
}

func listFiles(dir string) ([]string, error) {
	r := make([]string, 0, 64)
	err := filepath.Walk(filepath.Join(dir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		r = append(r, strings.TrimPrefix(strings.TrimPrefix(path, dir), string([]rune{filepath.Separator})))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (ti *TemplateInfo) String() string {
	b, err := yaml.Marshal(ti)
	if err != nil {
		panic(err)
	}
	return string(b)
}
