package util

import (
	"errors"
	"fmt"
	"go-mirayway/model"
	"strings"
)

var (
	LimitLengthPassword       = 60
	RequiredLengthPassword    = 6
	LimitLengthUsername       = 20
	RequiredLengthUsername    = 3
	LimitLengthTitlePost      = 200
	RequiredLengthTitlePost   = 1
	LimitLengthContentPost    = 500
	RequiredLengthContentPost = 5
)

func Limit(str string, limit, required int) bool {
	return len(str) >= required && len(str) <= limit
}

func ValidatePassword(password string) error {
	if !Limit(password, LimitLengthPassword, RequiredLengthPassword) {
		return errors.New(fmt.Sprintf("password must be >= %v and <= %v characters", RequiredLengthPassword, LimitLengthPassword))
	}

	if strings.Contains(password, " ") {
		return errors.New("password must not contain space character")
	}

	return nil
}

func ValidateUsername(username string) error {
	if !Limit(username, LimitLengthUsername, RequiredLengthUsername) {
		return errors.New(fmt.Sprintf("username must be >= %v and <= %v characters", RequiredLengthUsername, LimitLengthUsername))
	}

	if strings.Contains(username, " ") {
		return errors.New("username must not contain space character")
	}

	return nil
}

func ValidatePost(post model.PostRequest) error {
	if !Limit(post.Title, LimitLengthTitlePost, RequiredLengthTitlePost) {
		return errors.New(fmt.Sprintf("Title post must be >= %v and <= %v characters", RequiredLengthTitlePost, LimitLengthTitlePost))
	}

	if !Limit(post.Content, LimitLengthContentPost, RequiredLengthContentPost) {
		return errors.New(fmt.Sprintf("Title post must be >= %v and <= %v characters", RequiredLengthContentPost, LimitLengthContentPost))
	}

	return nil
}
