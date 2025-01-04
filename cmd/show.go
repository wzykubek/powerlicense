package cmd

import (
	"fmt"
	"os"
	"strings"

	t "go.wzykubek.xyz/licensmith/internal/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

func parseValueList(arr []string) string {
	titleCaser := cases.Title(language.English)
	var o string

	for k, v := range arr {
		v = strings.ReplaceAll(v, "-", " ")
		v = titleCaser.String(v)
		o = o + "- " + v
		if k < len(arr)-1 {
			o = o + "\n"
		}
	}

	return o
}

var showCmd = &cobra.Command{
	Use:   "show [license id]",
	Short: "Show license details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		licenseID := args[0]
		tmpl, err := t.Parse(licenseID + ".tmpl")
		if err != nil {
			fmt.Printf("Error: There is no '%s' license\n", licenseID)
			os.Exit(2)
		}

		permissions := parseValueList(tmpl.Permissions)
		conditions := parseValueList(tmpl.Conditions)
		limitations := parseValueList(tmpl.Limitations)
		fmt.Printf(
      "%s (%s)\n\n%s\n\nPermissions:\n%s\n\nConditions:\n%s\n\nLimitations:\n%s\n",
			tmpl.Title,
			tmpl.ID,
			tmpl.Description,
			permissions,
			conditions,
      limitations,
		)
	},
}
