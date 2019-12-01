package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func openFile(path string) *os.File {
	f, err := os.Open(path)
	check(err)
	return f
}

type fuelComputation struct {
	totalFuel    int
	computations []int
}

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f := openFile(*fptr)
	defer closeFile(f)

	s := bufio.NewScanner(f)

	fuelNeeded := 0
	for s.Scan() {
		moduleMass, _ := strconv.Atoi(s.Text())
		fc := computeFuelPart2(moduleMass, fuelComputation{})
		fuelNeeded += fc.totalFuel
	}
	err := s.Err()
	check(err)

	fmt.Printf("fuel needed %+v\n", fuelNeeded)
}

func computeFuelPart1(quantity int, actual fuelComputation) fuelComputation {
	fuel := quantity/3 - 2
	return fuelComputation{
		totalFuel:    fuel,
		computations: append(actual.computations, fuel)}
}

func computeFuelPart2(quantity int, actual fuelComputation) fuelComputation {
	fuel := quantity/3 - 2
	if fuel <= 0 {
		return fuelComputation{
			totalFuel:    actual.totalFuel,
			computations: append(actual.computations, 0)}
	}

	return computeFuelPart2(fuel,
		fuelComputation{
			totalFuel:    actual.totalFuel + fuel,
			computations: append(actual.computations, fuel)})
}
