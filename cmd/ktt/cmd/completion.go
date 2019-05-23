package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate auto completion",
}

func init() {
	CompletionCmd.AddCommand(
		&cobra.Command{
			Use:   "bash",
			Short: "generate bash auto completion",
			RunE:  genBashCompletion,
		},
		&cobra.Command{
			Use:   "zsh",
			Short: "generate zsh auto completion",
			RunE:  genZshCompletion,
		},
	)
}

func genCompletion(args []string, genFunc func(w io.Writer) error) error {
	var out io.Writer = os.Stdout
	if len(args) > 0 {
		if args[0] != "-" {
			f, err := os.Open(args[1])
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}
	}
	return genFunc(out)
}

func genBashCompletion(cmd *cobra.Command, args []string) error {
	return genCompletion(args, cmd.Root().GenBashCompletion)
}

func genZshCompletion(cmd *cobra.Command, args []string) error {
	return genCompletion(args, cmd.Root().GenZshCompletion)
}
