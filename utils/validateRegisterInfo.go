package utils

import (
	"errors"
	"regexp"
)

func ValidateEmail (email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func ValidatePwd (pwd string) error {
	if len(pwd) < 9 {
		return errors.New("Password too short.")
	}
	
	var atLeastOneCapital bool
	var atLeastOneLowerCase bool
	var atLeastOneSpecialChar bool

	strAsRunes := []rune(pwd)

	for i := range strAsRunes {
		runeVal := int(strAsRunes[i])

		switch {
		case runeVal > 0x20 && runeVal <= 0x2F:
			atLeastOneSpecialChar = true
		case runeVal > 0x40 && runeVal <= 0x5A:
			atLeastOneCapital = true
		case runeVal > 0x5A && runeVal <= 0x60:
			atLeastOneSpecialChar = true
		case runeVal > 0x60 && runeVal <= 0x7A:
			atLeastOneLowerCase = true
		case runeVal > 0x7A:
			atLeastOneSpecialChar = true
		}

		if atLeastOneCapital && atLeastOneLowerCase && atLeastOneSpecialChar {
			return nil
		}
	}

	if !atLeastOneCapital {
		return errors.New("Password requires at least one capital letter")
	}

	if !atLeastOneLowerCase {
		return errors.New("Password requires at least one lower case letter")
	}

	if !atLeastOneSpecialChar {
		return errors.New("Password requires at least one special character")
	}

	return nil
}
