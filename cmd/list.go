package cmd

import (
	"fmt"
  "strings"

  "go.wzykubek.xyz/licensmith/internal"

	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available licenses",
	Run: func(cmd *cobra.Command, args []string) {
    tmplList := internal.ListTemplates()
    fmt.Println(strings.Join(tmplList, ", "))
	},
}
