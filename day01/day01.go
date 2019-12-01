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

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f := openFile(*fptr)
	defer closeFile(f)

	s := bufio.NewScanner(f)

	fuelNeeded := 0
	for s.Scan() {
		moduleMass, _ := strconv.Atoi(s.Text())

		fc := FuelComputation{}
		fuelNeeded += fc.ComputeFuelPart2(moduleMass)
	}
	err := s.Err()
	check(err)

	fmt.Printf("fuel needed %+v\n", fuelNeeded)
}

// FuelComputation is a type to contain fuel computation state
type FuelComputation struct {
	totalFuel    int
	computations []int
}

// ComputeFuelPart1 is here to solve part 1
func (fc *FuelComputation) ComputeFuelPart1(quantity int) int {
	fuel := quantity/3 - 2
	fc.totalFuel = fuel
	fc.computations = append(fc.computations, fuel)
	return fuel
}

// ComputeFuelPart2 is here to solve part 2
func (fc *FuelComputation) ComputeFuelPart2(quantity int) int {
	fuel := quantity/3 - 2
	if fuel <= 0 {
		fc.computations = append(fc.computations, 0)
		return fc.totalFuel
	}

	fc.totalFuel += fuel
	fc.computations = append(fc.computations, fuel)
	return fc.ComputeFuelPart2(fuel)
}
