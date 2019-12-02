package day02

import (
	"adventofcode2019/common"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Run is the entrypoint of day02 exercice
func Run(objective int, filepath string) (int, error) {
	f := common.OpenFile(filepath)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)
	s.Scan()

	parts := strings.Split(s.Text(), ",")
	err := s.Err()
	if err != nil {
		return 0, err
	}

	seq := make([]int, 0)
	for _, elt := range parts {
		code, err := strconv.Atoi(elt)
		if err != nil {
			return 0, err
		}

		seq = append(seq, code)
	}

	for nounAttempt := 0; nounAttempt < 100; nounAttempt++ {
		for verbAttempt := 0; verbAttempt < 100; verbAttempt++ {
			// keep the initial sequence safe
			attempt := make([]int, len(seq))
			copy(attempt, seq)

			// launch a program execution
			program := IntCodeProgram{memory: attempt}
			result, err := program.Run(nounAttempt, verbAttempt)

			if err != nil {
				return 0, err
			}

			if result == objective {
				return 100*nounAttempt + verbAttempt, nil
			}
		}
	}

	return 0, errors.New(fmt.Sprint("no combination found to reach the objective: ", objective))
}

// IntCodeProgram contains the input data
type IntCodeProgram struct {
	memory   []int
	instrPtr int
}

// Run executes the program
func (p *IntCodeProgram) Run(noun, verb int) (int, error) {
	p.memory[1] = noun
	p.memory[2] = verb

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
