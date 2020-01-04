package day16

import (
	"adventofcode2019/common"
	"adventofcode2019/fft"
	"bufio"
)

// Run is the entryPoint of this day16 module
func Run(fileName string, steps int) (string, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)
	s.Scan()
	input := s.Text()
	err := s.Err()
	if err != nil {
		return "", err
	}

	myFft, err := fft.New([]int{0, 1, 0, -1}, input)
	if err != nil {
		return "", err
	}

	myFft.ProcessNSteps(steps)
	result := myFft.StringOfFirstN(8)

	return result, nil
}
