package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	OutputFile := flag.String("output", "LICENSE", "Specify different output file")
	LicenseID := flag.String("license", "", "Specify license by SPDX ID (e.g. BSD-3-Clause)")
	AuthorName := flag.String("name", "", "Set the author name (read from Git by default)")
	AuthorEmail := flag.String("email", "", "Set the author email (read from Git by default)")
	ListLicenses := flag.Bool("list", false, "List available licenses")
	flag.Parse()

	if *ListLicenses {
		listLicenses()
		os.Exit(0)
	}

	if *LicenseID == "" {
		fmt.Printf("Error: No license specified\n\nUse --license LICENSE\n\nAvailable licenses:\n")
		listLicenses()
		os.Exit(1)
	}

	licenseCtx, err := NewLicenseContext(*AuthorName, *AuthorEmail)
	if err != nil && err.Error() == "Can't read Git config" {
		fmt.Printf(
			"Error: Can't read Git config.\n\nUse --name \"NAME\" and --email EMAIL instead.\n",
		)
		os.Exit(3)
	}

	licenser := Licenser{
		LicenseID:      *LicenseID,
		LicenseContext: licenseCtx,
		OutputFile:     *OutputFile,
	}

	err = licenser.Generate()
	if err != nil && err.Error() == "Not supported license" {
		fmt.Printf("Error: There is no '%s' license\n\nAvailable licenses:\n", *LicenseID)
		listLicenses()
		os.Exit(2)
	}

	if err = licenser.WriteFile(); err != nil {
		panic(err)
	}
}
