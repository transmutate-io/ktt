package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version, Commit string

var VersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "show version",
	Args:    cobra.MaximumNArgs(0),
	RunE:    showVersion,
}

func showVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("Version: %s\nCommit: %s\n", Version, Commit)
	return nil
}
