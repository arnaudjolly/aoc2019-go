package day05

import (
	"adventofcode2019/common"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Run is the entrypoint of day05 exercice
func Run(filepath string) (int, error) {
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

	// keep the initial sequence safe
	attempt := make([]int, len(seq))
	copy(attempt, seq)

	// launch a program execution
	program := IntCodeProgram{memory: attempt}
	result, err := program.Run()
	if err != nil {
		return 0, err
	}

	fmt.Println("output: ", program.output)

	return result, nil
}

// IntCodeProgram contains the input data
type IntCodeProgram struct {
	memory   []int
	instrPtr int
	output   []int
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
	instrCode := p.memory[p.instrPtr]
	opcode := instrCode % 100
	switch opcode {
	case 1:
		p.ExecuteAdd()
	case 2:
		p.ExecuteMultiply()
	case 3:
		//part1 : p.ExecuteInput(1)
		//part2 :
		p.ExecuteInput(5)
	case 4:
		p.ExecuteOutput()
	case 5:
		p.ExecuteJumpIfTrue()
	case 6:
		p.ExecuteJumpIfFalse()
	case 7:
		p.ExecuteLessThan()
	case 8:
		p.ExecuteEquals()

	default:
		return errors.New("unknown opcode: " + strconv.Itoa(opcode))
	}
	return nil
}

// ExecuteEquals stores 1 in third if first == second else 0
func (p *IntCodeProgram) ExecuteEquals() {
	inst := p.memory[p.instrPtr : p.instrPtr+4]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)
	secondParam := p.resolveParam(1, inst, immediateParams)

	if firstParam == secondParam {
		p.memory[inst[3]] = 1
	} else {
		p.memory[inst[3]] = 0
	}

	p.instrPtr += 4
}

// ExecuteLessThan stores 1 in third if first < second else 0
func (p *IntCodeProgram) ExecuteLessThan() {
	inst := p.memory[p.instrPtr : p.instrPtr+4]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)
	secondParam := p.resolveParam(1, inst, immediateParams)

	if firstParam < secondParam {
		p.memory[inst[3]] = 1
	} else {
		p.memory[inst[3]] = 0
	}

	p.instrPtr += 4
}

// ExecuteJumpIfTrue jump to firstParam if non-zero
func (p *IntCodeProgram) ExecuteJumpIfTrue() {
	inst := p.memory[p.instrPtr : p.instrPtr+3]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)
	secondParam := p.resolveParam(1, inst, immediateParams)

	if firstParam != 0 {
		p.instrPtr = secondParam
	} else {
		p.instrPtr += 3
	}
}

// ExecuteJumpIfFalse jump to firstParam if zero
func (p *IntCodeProgram) ExecuteJumpIfFalse() {
	inst := p.memory[p.instrPtr : p.instrPtr+3]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)
	secondParam := p.resolveParam(1, inst, immediateParams)

	if firstParam == 0 {
		p.instrPtr = secondParam
	} else {
		p.instrPtr += 3
	}
}

// ExecuteInput simulate a "read" and insert i at the address coming next
func (p *IntCodeProgram) ExecuteInput(i int) {
	destIdx := p.memory[p.instrPtr+1]
	p.memory[destIdx] = i
	p.instrPtr += 2
}

// ExecuteOutput simulate a print
func (p *IntCodeProgram) ExecuteOutput() {
	inst := p.memory[p.instrPtr : p.instrPtr+2]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)

	p.output = append(p.output, firstParam)
	p.instrPtr += 2
}

// ExecuteAdd handles addition opcode
func (p *IntCodeProgram) ExecuteAdd() {
	inst := p.memory[p.instrPtr : p.instrPtr+4]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)
	secondParam := p.resolveParam(1, inst, immediateParams)

	p.memory[inst[3]] = firstParam + secondParam
	p.instrPtr += 4
}

func getImmediateParamMap(opCodeInstr int) map[int]bool {
	result := make(map[int]bool)
	modeOpCode := strconv.Itoa(opCodeInstr)
	for i := len(modeOpCode) - 3; i >= 0; i-- {
		if modeOpCode[i] == '1' {
			result[len(modeOpCode)-3-i] = true
		}
	}
	return result
}

// ExecuteMultiply handles multiplication opcode
func (p *IntCodeProgram) ExecuteMultiply() {
	inst := p.memory[p.instrPtr : p.instrPtr+4]
	immediateParams := getImmediateParamMap(inst[0])

	firstParam := p.resolveParam(0, inst, immediateParams)
	secondParam := p.resolveParam(1, inst, immediateParams)

	p.memory[inst[3]] = firstParam * secondParam
	p.instrPtr += 4
}

// Result is the temporary result of the program if not yet completed
// or the final result if it is
func (p *IntCodeProgram) Result() int {
	return p.memory[0]
}

func (p *IntCodeProgram) resolveParam(i int, instruction []int, modes map[int]bool) int {
	if modes[i] {
		return instruction[i+1]
	}
	return p.memory[instruction[i+1]]
}
