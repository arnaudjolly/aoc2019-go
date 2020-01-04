package common

import (
	"log"
	"os"
)

// CheckError log a Fatal if error not nil
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CloseFile allows to close an *os.File
func CloseFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// OpenFile opens a path and returns an *os.File
func OpenFile(path string) *os.File {
	f, err := os.Open(path)
	CheckError(err)
	return f
}

// AbsInt is the math.Abs for ints
func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// MinInt returns the minimum of a and b
func MinInt(a, b int) int {
	if a-b <= 0 {
		return a
	}
	return b
}

// MaxInt returns the minimum of a and b
func MaxInt(a, b int) int {
	if a-b >= 0 {
		return a
	}
	return b
}

// Lcm2 is Least Common Multiple
func Lcm2(a, b int) int {
	return a * b / Gcd(a, b)
}

// Lcm3 is Least Common Multiple variant with 3 numbers
func Lcm3(a, b, c int) int {
	lcm := Lcm2(a, b)
	return c * (lcm / Gcd(lcm, c))
}

// Gcd is Greatest common divisor
func Gcd(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	for {
		if a%b == 0 {
			return b
		}
		a, b = b, a%b
	}
}

// Eucl is Euclidean division:
// Eucl(a, b) = (q, r) with b = a * q + r
func Eucl(a, b int) (int, int) {
	return (a - a%b) / b, a % b
}

// EuclU64 is Euclidean division for larger numbers
func EuclU64(a, b uint64) (uint64, uint64) {
	return (a - a%b) / b, a % b
}

// SliceIntEquals informs about []int equality
func SliceIntEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
