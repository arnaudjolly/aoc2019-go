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

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	check(err)
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(f)

	fuelNeeded := int64(0)
	for s.Scan() {
		moduleMass, _ := strconv.ParseInt(s.Text(), 10, 0)
		fuelNeeded += computeFuelPart2(moduleMass)
	}
	err = s.Err()
	check(err)

	fmt.Printf("fuel needed %+v\n", fuelNeeded)
}

func computeFuelPart1(quantity int64) int64 {
	return quantity/3 - 2
}

func computeFuelPart2(quantity int64) int64 {
	fuelNeeded := quantity/3 - 2
	if fuelNeeded <= 0 {
		return 0
	}
	return fuelNeeded + computeFuelPart2(fuelNeeded)
}
