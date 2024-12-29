package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "licensmith",
	Short: "Licensmith is a LICENSE generator",
	Long:  "Effortlessly craft the perfect LICENSE for your Git repo in seconds with a single command!",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
