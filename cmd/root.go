package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "licensmith",
	Short: "Crafting the ideal license for your Git repository in seconds!",
	Long:  "Licensmith, a streamlined tool, allows you to create an LICENSE file for your Git repository with ease, using just one command. This tool is designed to save you time and effort.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
