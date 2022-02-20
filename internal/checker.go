package internal

import (
	"forum/pkg/models"
	"regexp"
	"strings"
)

func IsEmptyPost(p *models.Post) bool {
	if len(removeSpace(p.Title)) == 0 || len(p.Content) == 0 {
		return true
	}
	return false
}

func CheckName(s string) bool {
	if len(removeSpace(s)) == 0 {
		return false
	}
	loginConvention := "^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$"
	if re, _ := regexp.Compile(loginConvention); !re.MatchString(s) {
		return false
	}
	return true
}
func CheckEmail(email string) bool {
	if email == "" {
		return false
	}

	// Valid characters for the email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return false
	}
	return true
}

func removeSpace(s string) string {
	return strings.Trim(s, " ")
}
