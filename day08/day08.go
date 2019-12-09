package day08

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"strings"
)

// Run is the entrypoint of day08 exercice
func Run(fileName string, width, height int) (int, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)
	s.Scan()

	digits := s.Text()
	err := s.Err()
	if err != nil {
		return 0, err
	}

	result := part2(digits, width, height)

	return result, nil
}

func part2(digits string, width, height int) int {
	layers := layers(digits, width, height)

	var resultImage string

	for i := 0; i < len(layers); i++ {
		resultImage = reduceLayer(resultImage, layers[i])
	}

	fmt.Println("Final result:")
	showPicture(resultImage, width, height)
	return 0
}

func reduceLayer(topLayer, backLayer string) string {
	if len(topLayer) == 0 {
		return backLayer
	}

	result := make([]byte, len(backLayer))
	for idx, pixel := range topLayer {
		switch pixel {
		case '0':
			result[idx] = '0'
		case '1':
			result[idx] = '1'
		case '2':
			result[idx] = backLayer[idx]
		}
	}
	return string(result)
}

func showPicture(layer string, width, height int) {
	moreVisible := strings.ReplaceAll(layer, "0", " ")
	moreVisible = strings.ReplaceAll(moreVisible, "1", "@")
	for r := 0; r < len(moreVisible)/width; r++ {
		fmt.Println(moreVisible[r*width : (r+1)*width])
	}
}

func part1(digits string, width, height int) int {
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

	return nbOf1 * nbOf2
}

func layers(s string, width, height int) []string {
	result := make([]string, 0)
	for layer := 0; layer < len(s)/(width*height); layer++ {
		start := layer * width * height
		end := (layer + 1) * width * height
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[start:end])
	}
	return result
}
