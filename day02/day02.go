package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	s.Scan()

	parts := strings.Split(s.Text(), ",")
	check(s.Err())

	var program IntCodeProgram = make([]int, len(parts))
	for idx, elt := range parts {
		code, err := strconv.Atoi(elt)
		check(err)

		program[idx] = code
	}

	// init program
	program[1] = 12
	program[2] = 2

	// do the computation
out:
	for current := 0; current < len(program); current += 4 {
		instruction := program[current : current+4]

		switch instruction[0] {
		case 1:
			// opcode 1: addition
			program.Add(instruction)
		case 2:
			// opcode 2: multiplication
			program.Multiply(instruction)
		case 99:
			// halt
			break out
		default:
			log.Fatal(errors.New("bad opcode " + string(instruction[0])))
		}
	}

	fmt.Printf("result is %v\n", program[0])
}

// IntCodeProgram contains the input data
type IntCodeProgram []int

// Add handles addition opcode
func (p *IntCodeProgram) Add(inst []int) {
	(*p)[inst[3]] = (*p)[inst[1]] + (*p)[inst[2]]
}

// Multiply handles multiplication opcode
func (p *IntCodeProgram) Multiply(inst []int) {
	(*p)[inst[3]] = (*p)[inst[1]] * (*p)[inst[2]]
}

