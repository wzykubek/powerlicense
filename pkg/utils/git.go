package utils

import (
	"errors"
	"os/exec"
	"strings"
)

func GitUserData(key string) (string, error) {
	cmd := exec.Command("git", "config", "--get", key)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("can't read Git config")
	}

	value := strings.TrimSpace(string(out))
	return value, nil
}
