package util

import (
	"errors"
	"strings"
)

var limitLengthPassword = 4

func Limit(str string, limit int) bool {
	return len(str) < limit
}

func ValidationPassword(password string) error {
	if Limit(password, limitLengthPassword) {
		return errors.New("password must be more than 4 character")
	}

	if strings.Contains(password, " ") {
		return errors.New("password must not contain space")
	}

	return nil
}
