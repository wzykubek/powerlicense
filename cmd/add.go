package cmd

import (
	"fmt"
	"os"

	"go.wzykubek.xyz/licensmith/internal"

	"github.com/spf13/cobra"
)

var AuthorName string
var AuthorEmail string
var OutputFile string

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&AuthorName, "name", "", "Author name (read from Git by default)")
	addCmd.Flags().StringVar(&AuthorEmail, "email", "", "Author email (read from Git by default)")
	addCmd.Flags().StringVarP(&OutputFile, "output", "o", "LICENSE", "Output file")
}

var addCmd = &cobra.Command{
	Use:   "add [license]",
	Short: "Add LICENSE based on SPDX ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		licenseID := args[0]

		licenseCtx, err := internal.NewLicenseContext(AuthorName, AuthorEmail)
		if err != nil && err.Error() == "can't read Git config" {
			fmt.Println("Error: Can't read Git config")
			os.Exit(3)
		}

		licenser := internal.Licenser{
			LicenseID:      licenseID,
			LicenseContext: licenseCtx,
			OutputFile:     OutputFile,
		}

		err = licenser.Generate()
		if err != nil && err.Error() == "usupported license" {
			fmt.Printf("Error: There is no '%s' license\n", licenseID)
			os.Exit(2)
		}

		if err = licenser.WriteFile(); err != nil {
			panic(err)
		}
	},
}
