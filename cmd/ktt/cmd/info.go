package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"transmutate.io/pkg/ktt/ktt"
)

var InfoCmd = &cobra.Command{
	Use:     "info <project-path>",
	Aliases: []string{"i"},
	Short:   "show a templating project info",
	Args:    cobra.ExactArgs(1),
	RunE:    showProjectInfo,
}

func showProjectInfo(cmd *cobra.Command, args []string) error {
	tpl, err := ktt.LoadTemplate(args[0])
	if err != nil {
		return err
	}
	s := []string{
		"",
		"Info:",
		indentLines(tpl.Info.String()),
		"Template files:",
		indentLines(tpl.TemplateFiles.String()),
		"Files:",
		indentLines(tpl.Files.String()),
		"",
		"Values:",
		indentLines(ktt.ValuesYAML(tpl.Values)),
		"",
	}
	fmt.Fprint(os.Stdout, strings.Join(s, "\n"))
	return nil
}

func indentLines(s string) string {
	l := strings.Split(s, "\n")
	r := make([]string, 0, len(l))
	for _, i := range l {
		r = append(r, "  "+i)
	}
	return strings.Join(r, "\n")
}
