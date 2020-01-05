package main

import (
	"adventofcode2019/common"
	"adventofcode2019/day01"
	"adventofcode2019/day02"
	"adventofcode2019/day03"
	"adventofcode2019/day04"
	"adventofcode2019/day05"
	"adventofcode2019/day06"
	"adventofcode2019/day07"
	"adventofcode2019/day08"
	"adventofcode2019/day09"
	"adventofcode2019/day10"
	"adventofcode2019/day11"
	"adventofcode2019/day12"
	"adventofcode2019/day13"
	"adventofcode2019/day14"
	"adventofcode2019/day15"
	"adventofcode2019/day16"
	"adventofcode2019/day17"
	"flag"
	"fmt"
)

func main() {

	// specific for day02
	objptr := flag.Int("objective", 19690720, "objective to get")

	// specific for day04
	startptr := flag.Int("start", 123257, "start of day04 range")
	endptr := flag.Int("end", 647015, "end of day04 range")

	// specific for day08
	widthptr := flag.Int("width", 25, "width of layer")
	heightptr := flag.Int("height", 6, "height of layer")

	// specific for day12
	stepsptr := flag.Int("steps", 10, "nb of steps")
	intervalptr := flag.Int("interval", 1, "describe state every X step")

	// common flags
	fptr := flag.String("file", "input.txt", "file path to read from")
	dayptr := flag.Int("day", 17, "run the solution for day XX")
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
	case 5:
		result, err := day05.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 6:
		result, err := day06.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 7:
		result, err := day07.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 8:
		result, err := day08.Run(*fptr, *widthptr, *heightptr)
		common.CheckError(err)
		shareResult(result)
	case 9:
		result, err := day09.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 10:
		x, y, err := day10.Run(*fptr)
		common.CheckError(err)
		fmt.Printf("Coordinates: %v, %v\n", x, y)
	case 11:
		result, err := day11.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 12:
		result, err := day12.Run(*fptr, *stepsptr, *intervalptr)
		common.CheckError(err)
		shareResult(result)
	case 13:
		result, err := day13.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 14:
		result, err := day14.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 15:
		result, err := day15.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	case 16:
		result, err := day16.Run(*fptr, *stepsptr)
		common.CheckError(err)
		shareResult(result)
	case 17:
		result, err := day17.Run(*fptr)
		common.CheckError(err)
		shareResult(result)
	}
}

func shareResult(result interface{}) {
	fmt.Printf("Result: %v\n", result)
}
