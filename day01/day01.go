package day01

import (
	"adventofcode2019/common"
	"bufio"
	"strconv"
)

// Run is the entryPoint of this day02 module
func Run(fileName string) (int, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)

	fuelNeeded := 0
	for s.Scan() {
		moduleMass, _ := strconv.Atoi(s.Text())

		fc := FuelComputation{}
		fuelNeeded += fc.ComputeFuelPart2(moduleMass)
	}
	err := s.Err()
	if err != nil {
		return 0, nil
	}

	return fuelNeeded, nil
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
