package day04

import (
	"strconv"
)

// Run is the entrypoint of day04 exercice
func Run(start, end int) int {
	count := 0
	for i := start; i <= end; i++ {
		ok := CheckPasswordPart1(i)
		if ok {
			count = count + 1
		}
	}
	return count
}

// CheckPasswordPart1 returns true if the number is a valid password
func CheckPasswordPart1(i int) bool {
	s := strconv.Itoa(i)
	previousChar := byte('0') - 1
	var twoConsecutiveFound bool
	for idx := 0; idx < len(s); idx++ {
		c := s[idx]
		if previousChar > c {
			return false
		}
		twoConsecutiveFound = twoConsecutiveFound || c == previousChar
		previousChar = c
	}
	return twoConsecutiveFound
}
