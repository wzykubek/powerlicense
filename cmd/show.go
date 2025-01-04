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

var printTitle bool
var printDescription bool
var printPermissions bool
var printConditions bool
var printLimitations bool

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&printTitle, "title", false, "Print title")
	showCmd.Flags().BoolVar(&printDescription, "description", false, "Print description")
	showCmd.Flags().BoolVar(&printPermissions, "permissions", false, "Print permissions")
	showCmd.Flags().BoolVar(&printConditions, "conditions", false, "Print conditions")
	showCmd.Flags().BoolVar(&printLimitations, "limitations", false, "Print limitations")
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
		// If none of the flags is passed then print all attributes
		anyFlagChanged := false
		for _, flag := range []string{"title", "description", "permissions", "conditions", "limitations"} {
			if cmd.Flags().Changed(flag) {
				anyFlagChanged = true
				break
			}
		}
		if !anyFlagChanged {
			printTitle, printDescription, printPermissions, printConditions, printLimitations = true, true, true, true, true
		}

		licenseID := args[0]
		tmpl, err := t.Parse(licenseID + ".tmpl")
		if err != nil {
			fmt.Printf("Error: There is no '%s' license\n", licenseID)
			os.Exit(2)
		}

		printable := []string{}

		if printTitle {
			printable = append(printable, "Title: "+tmpl.Title)
		}

		if printDescription {
			printable = append(printable, "Description:\n"+tmpl.Description)
		}

		if printPermissions {
			printable = append(printable, "Permissions:\n"+parseValueList(tmpl.Permissions))
		}

		if printConditions {
			printable = append(printable, "Conditions:\n"+parseValueList(tmpl.Conditions))
		}

		if printLimitations {
			printable = append(printable, "Limitations:\n"+parseValueList(tmpl.Limitations))
		}

		fmt.Println(strings.Join(printable, "\n\n"))
	},
}
