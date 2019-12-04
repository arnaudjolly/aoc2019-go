package day04

import (
	"strconv"
)

// Run is the entrypoint of day04 exercice
func Run(start, end int) int {
	validator := PasswordValidator{}

	validCount := 0
	for i := start; i <= end; i++ {
		ok := validator.Validate(i)
		if ok {
			validCount = validCount + 1
		}
	}
	return validCount
}

// PasswordValidator the main type
type PasswordValidator struct {
	digit   byte
	count   int
	history []int
}

// Validate a password
func (v *PasswordValidator) Validate(password int) bool {
	s := strconv.Itoa(password)
	if len(s) != 6 {
		return false
	}

	v.history = make([]int, 0)
	v.digit = s[0]
	v.count = 0

	for idx := 0; idx < 6; idx++ {
		ok := v.handle(s[idx])
		if !ok {
			return false
		}
	}
	v.history = append(v.history, v.count)

	return v.checkResult()
}

func (v *PasswordValidator) checkResult() bool {
	for _, l := range v.history {
		if l >= 2 {
			return true
		}
	}
	return false
}

func (v *PasswordValidator) handle(c byte) bool {
	if c < v.digit {
		// if c is decreasing, stop there: it's not valid
		return false
	}

	if c != v.digit {
		// digit has changed
		v.changeDigit(c)
	}

	v.count++
	return true
}

func (v *PasswordValidator) changeDigit(c byte) {
	v.history = append(v.history, v.count)
	v.digit = c
	v.count = 0
}
