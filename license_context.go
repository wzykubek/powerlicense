package main

import (
	"time"
)

type LicenseContext struct {
	AuthorName  string
	AuthorEmail string
	Year        int
}

func NewLicenseContext(authorName string, authorEmail string) (LicenseContext, error) {
	var err error
	if authorName == "" {
		authorName, err = gitUserData("user.name")
	}
	if authorEmail == "" {
		authorEmail, err = gitUserData("user.email")
	}
	if err != nil {
		return LicenseContext{}, err
	}

	return LicenseContext{
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		Year:        time.Now().Year(),
	}, nil
}
