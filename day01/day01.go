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
		fuelNeeded += moduleMass/3 - 2
	}
	err = s.Err()
	check(err)

	fmt.Printf("fuel needed %+v\n", fuelNeeded)
}
