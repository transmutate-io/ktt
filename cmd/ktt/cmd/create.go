package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/transmutateio/ktt/ktt"
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
	files := map[string][]byte{
		"info.yaml":   []byte(tplInfo.String()),
		"values.yaml": []byte(valuesFileComment),
		filepath.Join("templates", "helpers.tpl"): []byte(templateHelpers),
	}
	for i, b := range files {
		err := ioutil.WriteFile(filepath.Join(args[0], i), b, 0644)
		if err != nil {
			return err
		}
	}
	return nil
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
	templateHelpers = `# helpers
{{- define "util.name" -}}
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
