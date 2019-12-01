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
		fuelNeeded += computeFuelPart2(moduleMass)
	}
	err := s.Err()
	check(err)

	fmt.Printf("fuel needed %+v\n", fuelNeeded)
}

func computeFuelPart1(quantity int) int {
	return quantity/3 - 2
}

func computeFuelPart2(quantity int) int {
	fuelNeeded := quantity/3 - 2
	if fuelNeeded <= 0 {
		return 0
	}
	return fuelNeeded + computeFuelPart2(fuelNeeded)
}
