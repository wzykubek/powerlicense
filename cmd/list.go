package cmd

import (
	"fmt"
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
		templates := t.List()
		fmt.Println(strings.Join(templates, ", "))
	},
}
