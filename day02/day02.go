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

	seq := make([]int, 0)
	for _, elt := range parts {
		code, err := strconv.Atoi(elt)
		check(err)

		seq = append(seq, code)
	}

	program := IntCodeProgram{memory: seq}

	// init program
	program.Init(12, 2)

	// do the computation
	result, err := program.Run()
	check(err)

	// print the result
	fmt.Printf("result is %v\n", result)
}

// IntCodeProgram contains the input data
type IntCodeProgram struct {
	memory   []int
	instrPtr int
}

// Init the program
func (p *IntCodeProgram) Init(noun, verb int) {
	p.memory[1] = noun
	p.memory[2] = verb
}

// Run executes the program
func (p *IntCodeProgram) Run() (int, error) {
	for !p.IsCompleted() {
		err := p.ExecuteNextInstruction()
		if err != nil {
			return 0, err
		}
	}
	return p.Result(), nil
}

// IsCompleted informs about the completeness of the program
func (p *IntCodeProgram) IsCompleted() bool {
	return p.memory[p.instrPtr] == 99
}

// ExecuteNextInstruction identifies instruction to execute and do it
func (p *IntCodeProgram) ExecuteNextInstruction() error {
	switch p.memory[p.instrPtr] {
	case 1:
		p.ExecuteAdd()
	case 2:
		p.ExecuteMultiply()
	default:
		return errors.New("unknown opcode: " + string(p.memory[p.instrPtr]))
	}
	return nil
}

// ExecuteAdd handles addition opcode
func (p *IntCodeProgram) ExecuteAdd() {
	inst := p.memory[p.instrPtr : p.instrPtr+4]
	p.memory[inst[3]] = p.memory[inst[1]] + p.memory[inst[2]]
	p.instrPtr += 4
}

// ExecuteMultiply handles multiplication opcode
func (p *IntCodeProgram) ExecuteMultiply() {
	inst := p.memory[p.instrPtr : p.instrPtr+4]
	p.memory[inst[3]] = p.memory[inst[1]] * p.memory[inst[2]]
	p.instrPtr += 4
}

// Result is the temporary result of the program if not yet completed
// or the final result if it is
func (p *IntCodeProgram) Result() int {
	return p.memory[0]
}
