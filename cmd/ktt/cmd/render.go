package cmd

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var RenderCmd = &cobra.Command{
	Use:     "render <project-path> <deployment-name>",
	Aliases: []string{"r"},
	Short:   "render a project",
	Args:    cobra.ExactArgs(2),
	RunE:    renderProject,
}

func init() {
	fs = RenderCmd.PersistentFlags()
	addCommonFlags(fs)
	fs.BoolVarP(&outputDir, "dir", "d", false, "output as a directory structure")
	fs.StringVarP(&output, "out", "o", "", "output file/directory")
	fs.BoolVarP(&valuesOnly, "values", "s", false, "render values and show")
}

var (
	output     string
	outputDir  bool
	valuesOnly bool
)

func renderProject(cmd *cobra.Command, args []string) error {
	tpl, err := loadProject(cmd, args)
	if err != nil {
		return err
	}
	if valuesOnly {
		b, err := yaml.Marshal(tpl.Values)
		if err != nil {
			return err
		}
		var fout io.Writer
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
		_, err = fmt.Fprintln(fout, string(b))
		return err
	}
	if !cmd.PersistentFlags().Changed("dir") {
		if fi, err := os.Stat(output); err == nil && fi.IsDir() {
			outputDir = true
		}
	}
	if outputDir && output == "" {
		output = "."
	}
	return render(tpl, args[1])
}
