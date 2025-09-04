package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate aws-sso command docs",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		if dir == "" {
			dir = os.TempDir()
		}

		return docsAction(os.Stdout, dir)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringP("dir", "d", "", "Destination directory for docs")
}

func docsAction(out io.Writer, dir string) error {
	if err := doc.GenMarkdownTreeCustom(rootCmd, dir,
		func(s string) string { return s },
		func(s string) string {
			b, found := strings.CutSuffix(s, ".md")
			if !found {
				panic("could not find .md suffix")
			}
			return "../" + b
		}); err != nil {
		return err
	}

	_, err := fmt.Fprintf(out, "Documentation successfully created in %s\n", dir)
	return err
}
