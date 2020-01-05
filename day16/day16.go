package day16

import (
	"adventofcode2019/common"
	"adventofcode2019/fft"
	"bufio"
	"fmt"
	"strings"
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

	realSignal := strings.Repeat(input, 10000)
	myFft, err := fft.New(realSignal)
	if err != nil {
		return "", err
	}

	fmt.Printf("fft:%v\n", myFft)

	myFft.ProcessNSteps(steps)
	result := myFft.StringOfFirstNWithOffset(8)

	return result, nil
}
