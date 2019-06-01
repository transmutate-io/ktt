package main

import (
	"os"

	"github.com/spf13/cobra"
	"transmutate.io/pkg/ktt/cmd/ktt/cmd"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "ktt is a Kubernetes Templating Tool",
}

func init() {
	rootCmd.AddCommand(
		cmd.CreateCmd,
		cmd.CompletionCmd,
		cmd.InfoCmd,
		cmd.RenderCmd,
		cmd.VersionCmd,
		cmd.LintCmd,
	)
}
