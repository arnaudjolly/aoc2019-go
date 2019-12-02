package main

import (
	"adventofcode2019/common"
	"adventofcode2019/day01"
	"adventofcode2019/day02"
	"flag"
	"fmt"
)

func main() {

	dayptr := flag.Int("day", 2, "run the solution for day XX")
	objptr := flag.Int("objective", 19690720, "objective to get")
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	var result int
	var err error

	switch *dayptr {
	case 1:
		result, err = day01.Run(*fptr)
	case 2:
		result, err = day02.Run(*objptr, *fptr)
	}
	common.CheckError(err)

	fmt.Printf("result: %v\n", result)
}
