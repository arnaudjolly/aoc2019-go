package day08

import (
	"adventofcode2019/common"
	"bufio"
)

// Run is the entrypoint of day08 exercice
func Run(fileName string) (int, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)
	s.Scan()

	digits := s.Text()
	err := s.Err()
	if err != nil {
		return 0, err
	}

	height := 6
	width := 25

	minZeroes := 25
	var resultLayer string

	for _, layer := range layers(digits, width, height) {

		nbZero := 0
		for _, d := range layer {
			if d == '0' {
				nbZero++
			}
		}

		if nbZero < minZeroes {
			resultLayer, minZeroes = layer, nbZero
		}
	}

	nbOf1 := 0
	nbOf2 := 0
	for _, d := range resultLayer {
		switch d {
		case '1':
			nbOf1++
		case '2':
			nbOf2++
		}
	}
	return nbOf1 * nbOf2, nil
}

func layers(s string, width, height int) []string {
	result := make([]string, 0)
	for layer := 0; layer < len(s)/(width*height); layer++ {
		result = append(result, s[layer*width*height:(layer+1)*width*height])
	}
	return result
}
