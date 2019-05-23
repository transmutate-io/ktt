package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var LintCmd = &cobra.Command{
	Use:     "lint <project-path>",
	Aliases: []string{"l"},
	Short:   "lint a project",
	Args:    cobra.ExactArgs(1),
	RunE:    lintProject,
}

func init() {
	fs = LintCmd.PersistentFlags()
	addCommonFlags(fs)
}

func lintProject(cmd *cobra.Command, args []string) error {
	tpl, err := loadProject(cmd, args)
	if err != nil {
		return err
	}
	outputDir = false
	output = os.DevNull
	return render(tpl, tpl.Info.App.Name)
}
