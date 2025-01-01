package cmd

import (
	"fmt"
	"os"

	l "go.wzykubek.xyz/licensmith/internal/license"

	"github.com/spf13/cobra"
)

var authorName string
var authorEmail string
var outputFile string

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&authorName, "name", "", "Author name (read from Git by default)")
	addCmd.Flags().StringVar(&authorEmail, "email", "", "Author email (read from Git by default)")
	addCmd.Flags().StringVarP(&outputFile, "output", "o", "LICENSE", "Output file")
}

var addCmd = &cobra.Command{
	Use:   "add [license id]",
	Short: "Add LICENSE based on SPDX ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		licenseID := args[0]

		ctx, err := l.NewContext(authorName, authorEmail)
		if err != nil && err.Error() == "can't read Git config" {
			fmt.Println("Error: Can't read Git config")
			os.Exit(3)
		}

		license := l.License{
			ID:      licenseID,
			Context: ctx,
		}

		err = license.Gen()
		if err != nil {
			if err.Error() == "usupported license" {
				fmt.Printf("Error: There is no '%s' license\n", licenseID)
				os.Exit(2)
			} else {
        fmt.Println("Internal Error:", err)
        os.Exit(127)
      }
		}

		if err = license.Write(outputFile); err != nil {
      fmt.Println("Error: Can't write file:", err)
		}
	},
}
