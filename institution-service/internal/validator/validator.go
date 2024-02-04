package validator

import "regexp"

var (
	PhoneRX = regexp.MustCompile("^\\+\\d{11}$")
)
