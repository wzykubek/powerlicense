package license

import (
	"go.wzykubek.xyz/licensmith/pkg/utils"

	"time"
)

type Context struct {
	AuthorName  string
	AuthorEmail string
	Year        int
}

func NewContext(authorName string, authorEmail string) (Context, error) {
	var err error
	if authorName == "" {
		authorName, err = utils.GitUserData("user.name")
	}
	if authorEmail == "" {
		authorEmail, err = utils.GitUserData("user.email")
	}
	if err != nil {
		return Context{}, err
	}

	return Context{
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		Year:        time.Now().Year(),
	}, nil
}
