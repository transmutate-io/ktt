package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"transmutate.io/pkg/ktt/ktt"
)

var (
	CreateCmd = &cobra.Command{
		Use:     "create <project-path> <app-name> <app-version>",
		Aliases: []string{"c"},
		Short:   "create a new templating project",
		Args:    cobra.ExactArgs(3),
		RunE:    createProject,
	}
	tplInfo = &ktt.TemplateInfo{
		App:      &ktt.Info{},
		Template: &ktt.Info{},
	}
)

func init() {
	fs := CreateCmd.PersistentFlags()
	fs.StringVarP(&tplInfo.Template.Name, "name", "n", "", "template name. defaults to app name if not set")
	fs.StringVarP(&tplInfo.Template.Version, "version", "v", "", "template version. defaults to app version if not set")
}

func createProject(cmd *cobra.Command, args []string) error {
	for _, i := range []string{"templates", "files"} {
		err := os.MkdirAll(filepath.Join(args[0], i), os.ModeDir|0755)
		if err != nil {
			return err
		}
	}
	// create info file
	tplInfo.App.Name = args[1]
	tplInfo.App.Version = args[2]
	if tplInfo.Template.Name == "" {
		tplInfo.Template.Name = tplInfo.App.Name
	}
	if tplInfo.Template.Version == "" {
		tplInfo.Template.Version = tplInfo.App.Version
	}
	files := map[string][]byte{"info.yaml": []byte(tplInfo.String())}
	var (
		hasCustomValues    bool
		hasCustomTemplates bool
	)
	cfgDir := configDirectory()
	templatesDir := filepath.Join(cfgDir, "templates")
	fi, err := os.Stat(templatesDir)
	if err == nil {
		if !fi.IsDir() {
			return errors.New("invalid custom user templates directory")
		}
		hasCustomTemplates = true
	} else {
		if e, ok := err.(*os.PathError); !ok || e.Err != syscall.ENOENT {
			return err
		}
		files[filepath.Join("templates", "helpers.tpl")] = []byte(templateHelpers)
	}
	valuesYamlFile := filepath.Join(cfgDir, "values.yaml")
	if fi, err = os.Stat(valuesYamlFile); err == nil {
		if fi.IsDir() {
			return errors.New("invalid custom user values.yaml file")
		}
		hasCustomValues = true
	} else {
		if e, ok := err.(*os.PathError); !ok || e.Err != syscall.ENOENT {
			return err
		}
		files["values.yaml"] = []byte(valuesFileComment)
	}
	for i, b := range files {
		err := ioutil.WriteFile(filepath.Join(args[0], i), b, 0644)
		if err != nil {
			return err
		}
	}
	if !hasCustomTemplates && !hasCustomValues {
		return nil
	}
	cpArgs := append(make([]string, 0, 5), "-r", "-f")
	if hasCustomTemplates {
		cpArgs = append(cpArgs, templatesDir)
	}
	if hasCustomValues {
		cpArgs = append(cpArgs, valuesYamlFile)
	}
	cpArgs = append(cpArgs, args[0])
	if err = exec.Command("cp", cpArgs...).Run(); err != nil {
		return err
	}
	return nil
}

func configDirectory() string {
	r := os.ExpandEnv("$XDG_CONFIG_HOME")
	if r == "" {
		r = filepath.Join(os.ExpandEnv("$HOME"), ".config")
	}
	return filepath.Join(r, "ktt")
}

const (
	valuesFileComment = `# configure metadata
meta:
  # override default name
  name: ""
  # override namespace
  namespace: ""
  # add extra common labels
  labels: {}
  # add extra common annotations
  annotations: {}
`
	templateHelpers = `{{- define "util.name" -}}
{{- default .Name .Values.meta.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "util.namespace" -}}
{{ if .Values.meta.namespace -}}namespace: {{ .Values.meta.namespace }}{{- end }}
{{- end -}}

{{- define "util.labels" -}}
app.kubernetes.io/instance: {{ include "util.name" . }}
app.kubernetes.io/name: {{ .App.Name }}
app.kubernetes.io/version: {{ .App.Version }}
template/name: {{ .Template.Name }}
template/version: {{ .Template.Version }}
{{ if .Values.meta.labels -}}{{- toYaml .Values.meta.labels -}}{{- end }}
{{- end -}}

{{- define "util.annotations" -}}
{{ if .Values.meta.annotations -}}{{- toYaml .Values.meta.annotations -}}{{- end }}
{{- end -}}
`
)
