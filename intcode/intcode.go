package intcode

import (
	"errors"
	"strconv"
)

// ProgramCreator allows you to create an instance of Program
func ProgramCreator(state []int) func() *Program {
	// keep the initial sequence safe
	safeBackup := make([]int, len(state))
	copy(safeBackup, state)

	return func() *Program {
		attempt := make([]int, len(safeBackup))
		copy(attempt, safeBackup)
		return &Program{program: attempt}
	}
}

// Program contains the input data
type Program struct {
	program      []int
	instrPtr     int
	extraMemory  map[int]int
	input        chan int
	output       chan int
	halted       bool
	relativeBase int
}

// Run executes the program
func (p *Program) Run(in, out, quit chan int) error {
	p.input = in
	p.output = out

	for !p.halted {
		err := p.ExecuteNextInstruction()
		if err != nil {
			return err
		}
	}
	quit <- 0
	return nil
}

// MemorySlice returns a slice of memory
// mixing program and extraMemory storage
func (p *Program) MemorySlice(start, end int) []int {
	result := make([]int, end-start)
	for idx := range result {
		result[idx] = p.MemoryAt(start + idx)
	}
	return result
}

// MemoryAt returns the value at a specific address
// this allows to retrieve values outside program memory space
func (p *Program) MemoryAt(address int) int {
	if address < len(p.program) {
		return p.program[address]
	}
	return p.extraMemory[address-len(p.program)]
}

// SetMemory allows to set a value at address mixing program and extraMemory
func (p *Program) SetMemory(address int, v int) {
	if address < len(p.program) {
		p.program[address] = v
	} else {
		if p.extraMemory == nil {
			p.extraMemory = make(map[int]int)
		}
		p.extraMemory[address-len(p.program)] = v
	}
}

// IsCompleted informs about the completeness of the program
func (p *Program) IsCompleted() bool {
	return p.MemoryAt(p.instrPtr) == 99
}

// ExecuteNextInstruction identifies instruction to execute and do it
func (p *Program) ExecuteNextInstruction() error {
	instrCode := p.MemoryAt(p.instrPtr)
	opcode := instrCode % 100
	switch opcode {
	case 1:
		p.ExecuteAdd()
	case 2:
		p.ExecuteMultiply()
	case 3:
		p.ExecuteInput()
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
	case 9:
		p.ExecuteRelativeBaseOffset()
	case 99:
		close(p.output)
		p.halted = true
	default:
		return errors.New("unknown opcode: " + strconv.Itoa(opcode))
	}
	return nil
}

// ExecuteRelativeBaseOffset adjusts the relative base
func (p *Program) ExecuteRelativeBaseOffset() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+2)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)

	p.relativeBase += firstParam

	p.instrPtr += 2
}

// ExecuteEquals stores 1 in third if first == second else 0
func (p *Program) ExecuteEquals() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+4)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)
	dest := p.resolveDestination(2, inst, paramModes)

	if firstParam == secondParam {
		p.SetMemory(dest, 1)
	} else {
		p.SetMemory(dest, 0)
	}

	p.instrPtr += 4
}

// ExecuteLessThan stores 1 in third if first < second else 0
func (p *Program) ExecuteLessThan() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+4)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)
	dest := p.resolveDestination(2, inst, paramModes)

	if firstParam < secondParam {
		p.SetMemory(dest, 1)
	} else {
		p.SetMemory(dest, 0)
	}

	p.instrPtr += 4
}

// ExecuteJumpIfTrue jump to firstParam if non-zero
func (p *Program) ExecuteJumpIfTrue() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+3)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)

	if firstParam != 0 {
		p.instrPtr = secondParam
	} else {
		p.instrPtr += 3
	}
}

// ExecuteJumpIfFalse jump to firstParam if zero
func (p *Program) ExecuteJumpIfFalse() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+3)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)

	if firstParam == 0 {
		p.instrPtr = secondParam
	} else {
		p.instrPtr += 3
	}
}

// ExecuteInput simulate a "read" and insert input at the address coming next
func (p *Program) ExecuteInput() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+2)
	paramModes := getParamModes(inst[0])

	dest := p.resolveDestination(0, inst, paramModes)

	p.SetMemory(dest, <-p.input)
	p.instrPtr += 2
}

// ExecuteOutput simulate a print
func (p *Program) ExecuteOutput() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+2)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)

	p.output <- firstParam
	p.instrPtr += 2
}

// ExecuteAdd handles addition opcode
func (p *Program) ExecuteAdd() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+4)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)
	dest := p.resolveDestination(2, inst, paramModes)

	p.SetMemory(dest, firstParam+secondParam)
	p.instrPtr += 4
}

func getParamModes(opCodeInstr int) map[int]int {
	result := make(map[int]int)
	modeOpCode := strconv.Itoa(opCodeInstr)
	for i := len(modeOpCode) - 3; i >= 0; i-- {
		switch modeOpCode[i] {
		case '0':
			result[len(modeOpCode)-3-i] = 0
		case '1':
			result[len(modeOpCode)-3-i] = 1
		case '2':
			result[len(modeOpCode)-3-i] = 2
		}
	}
	return result
}

// ExecuteMultiply handles multiplication opcode
func (p *Program) ExecuteMultiply() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+4)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)
	dest := p.resolveDestination(2, inst, paramModes)

	p.SetMemory(dest, firstParam*secondParam)
	p.instrPtr += 4
}

func (p *Program) resolveParam(i int, instruction []int, modes map[int]int) int {
	mode, found := modes[i]
	if !found {
		mode = 0
	}

	param := instruction[i+1]
	switch mode {
	case 0:
		// address mode
		return p.MemoryAt(param)
	case 1:
		// immediate mode
		return param
	case 2:
		// relative mode
		return p.MemoryAt(p.relativeBase + param)
	default:
		// should never happens
		return 0
	}
}

func (p *Program) resolveDestination(i int, instruction []int, modes map[int]int) int {
	dest := instruction[i+1]
	if modes[i] == 2 {
		dest += p.relativeBase
	}
	return dest
}
