package main

import (
	"adventofcode2019/common"
	"adventofcode2019/day01"
	"adventofcode2019/day02"
	"adventofcode2019/day03"
	"adventofcode2019/day04"
	"flag"
	"fmt"
)

func main() {

	// specific for day02
	objptr := flag.Int("objective", 19690720, "objective to get")

	// specific for day04
	startptr := flag.Int("start", 123257, "start of day04 range")
	endptr := flag.Int("end", 647015, "end of day04 range")

	// common flags
	fptr := flag.String("file", "input.txt", "file path to read from")
	dayptr := flag.Int("day", 4, "run the solution for day XX")
	flag.Parse()

	switch *dayptr {
	case 1:
		result, err := day01.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 2:
		result, err := day02.Run(*objptr, *fptr)
		common.CheckError(err)
		shareResult(result)
	case 3:
		result, err := day03.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 4:
		result := day04.Run(*startptr, *endptr)
		shareResult(result)
	}
}

func shareResult(result interface{}) {
	fmt.Printf("Result: %v\n", result)
}
