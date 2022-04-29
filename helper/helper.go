package helper

import "regexp"

func CheckPhone(phone string) bool {
	validPhone := regexp.MustCompile("^[0-9]{10}$")
	return validPhone.MatchString(phone)
}

func CheckName(name string) bool {
	validName := regexp.MustCompile("^[a-zA-Z]{3,}$")
	return validName.MatchString(name)
}
