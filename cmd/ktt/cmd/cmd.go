package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/transmutateio/ktt/ktt"
	"gopkg.in/yaml.v2"
)

var (
	fs            *pflag.FlagSet
	vals          = []string{}
	valsYAML      = []string{}
	valsFromFiles = []string{}
	valsFiles     = []string{}
	onlyFiles     = []string{}
)

func addCommonFlags(fs *pflag.FlagSet) {
	fs.StringSliceVarP(&vals, "val", "v", nil, "set a value (NAME=value)")
	fs.StringSliceVarP(&valsYAML, "yaml", "y", nil, "set a value from a YAML file (NAME=/file/path.yaml)")
	fs.StringSliceVarP(&valsFromFiles, "file", "f", nil, "set a value from a file (NAME=/file/path)")
	fs.StringSliceVarP(&valsFiles, "merge", "m", nil, "merge a values file")
	fs.StringSliceVarP(&onlyFiles, "template", "t", nil, "parse only specific files")
}

var (
	errInvalidFormatValue = errors.New("missing name/value")
	errInvalidFormatFile  = errors.New("missing name/file")
)

func loadProject(cmd *cobra.Command, args []string) (*ktt.Template, error) {
	tpl, err := ktt.LoadTemplate(args[0])
	if err != nil {
		return nil, err
	}
	// merge provided values files
	for _, i := range valsFiles {
		b, err := ioutil.ReadFile(i)
		if err != nil {
			return nil, err
		}
		v := map[interface{}]interface{}{}
		if err := yaml.Unmarshal(b, v); err != nil {
			return nil, err
		}
		ktt.MergeValues(tpl.Values, v)
	}
	// merge fields with values from files
	for _, i := range valsFromFiles {
		n := strings.SplitN(i, "=", 2)
		if len(n) != 2 {
			return nil, errInvalidFormatFile
		}
		b, err := ioutil.ReadFile(n[1])
		if err != nil {
			return nil, err
		}
		ktt.SetValuesKey(tpl.Values, n[0], string(b))
	}
	// merge other fields from YAML files
	for _, i := range valsYAML {
		n := strings.SplitN(i, "=", 2)
		if len(n) != 2 {
			return nil, errInvalidFormatFile
		}
		b, err := ioutil.ReadFile(n[1])
		if err != nil {
			return nil, err
		}
		v := map[interface{}]interface{}{}
		if err = yaml.Unmarshal(b, v); err != nil {
			return nil, err
		}
		ktt.SetValuesKey(tpl.Values, n[0], v)
	}
	// merge other fields
	for _, i := range vals {
		n := strings.SplitN(i, "=", 2)
		if len(n) != 2 {
			return nil, errInvalidFormatValue
		}
		ktt.SetValuesKey(tpl.Values, n[0], n[1])
	}
	return tpl, nil
}

func render(tpl *ktt.Template, name string) error {
	var fout io.Writer
	if !outputDir {
		if output == "" || output == "-" {
			fout = os.Stdout
		} else {
			f, err := os.Create(output)
			if err != nil {
				return err
			}
			defer f.Close()
			fout = f
		}
	}
	if len(onlyFiles) > 0 {
		t := make(ktt.Templates, len(tpl.TemplateFiles))
		for _, filename := range onlyFiles {
			tt, ok := tpl.TemplateFiles[filename]
			if !ok {
				return fmt.Errorf("template not found: %s", filename)
			}
			t[filename] = tt
		}
		tpl.TemplateFiles = t
	}
	for filename, t := range tpl.TemplateFiles {
		var f *os.File
		if outputDir {
			fn := filepath.Join(output, filename)
			d, _ := filepath.Split(fn)
			if d != "" {
				if err := os.MkdirAll(d, os.ModeDir|0755); err != nil {
					return err
				}
			}
			f, err := os.Create(fn)
			if err != nil {
				return err
			}
			fout = f
		}
		fmt.Fprintf(fout, "---\n# source: %s\n\n", filename)
		err := t.Execute(fout, &templateData{
			Name:     name,
			App:      tpl.Info.App,
			Template: tpl.Info.Template,
			Values:   tpl.Values,
			Files:    tpl.Files,
		})
		if err != nil {
			return err
		}
		fmt.Fprint(fout, "\n\n")
		if outputDir {
			f.Close()
		}
	}
	return nil
}

type templateData struct {
	Name     string
	App      *ktt.Info
	Template *ktt.Info
	Values   ktt.Values
	Files    ktt.Files
}
