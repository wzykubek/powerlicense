package cmd

import (
	"fmt"
	"os"
	"strings"

	t "go.wzykubek.xyz/licensmith/internal/template"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available licenses",
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := t.List()
    if err != nil {
      fmt.Println("Internal Error:", err)
      os.Exit(127)
    }

		fmt.Println(strings.Join(templates, ", "))
	},
}
