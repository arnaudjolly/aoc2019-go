package day13

import (
	"adventofcode2019/common"
	"bufio"
	"errors"
	"fmt"
	col "github.com/fatih/color"
	"strconv"
	"strings"
)

// Run is the entrypoint of day13 exercice
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

	createProgram := programCreator(seq)
	p := createProgram()
	out := make(chan int)

	c := challenge{grid: make(map[point]tile)}

	go p.Run(nil, out)

	for !p.halted {
		x, y, t := <-out, <-out, tile(<-out)
		c.placeTile(x, y, t)
	}

	c.printGrid()

	blocks := c.count(isBlock)

	return blocks, nil
}

var (
	bgW = col.New(col.BgWhite, col.FgBlack)
	bgB = col.New(col.BgBlack, col.FgWhite)
	bgR = col.New(col.BgRed, col.FgBlack)
)

type challenge struct {
	grid map[point]tile
}

func (c *challenge) placeTile(x, y int, t tile) {
	p := point{x, y}
	c.grid[p] = t
}

func (c *challenge) count(tileFilter func(tile) bool) int {
	result := 0
	for _, t := range c.grid {
		if tileFilter(t) {
			result++
		}
	}
	return result
}

func (c *challenge) printGrid() {
	fmt.Println("Grid:")
	minX, maxX := 1, -1
	minY, maxY := 1, -1

	for p := range c.grid {
		if p.x > maxX {
			maxX = p.x
		}
		if p.x < minX {
			minX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
		if p.y < minY {
			minY = p.y
		}
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}
			tile, found := c.grid[p]
			if !found {
				tile = empty
			}
			fmt.Print(tile)
		}
		fmt.Printf("\n")
	}
}

type tile int

func (t tile) String() string {
	switch t {
	case empty:
		return bgW.Sprint(" ")
	case wall:
		return bgB.Sprint(" ")
	case block:
		return bgW.Sprint("X")
	case hpaddle:
		return bgW.Sprint("-")
	case ball:
		return bgW.Sprint("o")
	default:
		return bgR.Sprintf("%1v", int(t))
	}
}

const (
	empty   = tile(0)
	wall    = tile(1)
	block   = tile(2)
	hpaddle = tile(3)
	ball    = tile(4)
)

func isBlock(t tile) bool { return t == block }

type point struct {
	x int
	y int
}

func programCreator(state []int) func() *IntCodeProgram {
	// keep the initial sequence safe
	safeBackup := make([]int, len(state))
	copy(safeBackup, state)

	return func() *IntCodeProgram {
		attempt := make([]int, len(safeBackup))
		copy(attempt, safeBackup)
		return &IntCodeProgram{program: attempt}
	}
}

// IntCodeProgram contains the input data
type IntCodeProgram struct {
	program      []int
	instrPtr     int
	extraMemory  map[int]int
	input        chan int
	output       chan int
	halted       bool
	relativeBase int
}

// Run executes the program
func (p *IntCodeProgram) Run(in, out chan int) error {
	p.input = in
	p.output = out

	for !p.halted {
		err := p.ExecuteNextInstruction()
		if err != nil {
			return err
		}
	}
	return nil
}

// MemorySlice returns a slice of memory
// mixing program and extraMemory storage
func (p *IntCodeProgram) MemorySlice(start, end int) []int {
	result := make([]int, end-start)
	for idx := range result {
		result[idx] = p.MemoryAt(start + idx)
	}
	return result
}

// MemoryAt returns the value at a specific address
// this allows to retrieve values outside program memory space
func (p *IntCodeProgram) MemoryAt(address int) int {
	if address < len(p.program) {
		return p.program[address]
	}
	return p.extraMemory[address-len(p.program)]
}

// SetMemory allows to set a value at address mixing program and extraMemory
func (p *IntCodeProgram) SetMemory(address int, v int) {
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
func (p *IntCodeProgram) IsCompleted() bool {
	return p.MemoryAt(p.instrPtr) == 99
}

// ExecuteNextInstruction identifies instruction to execute and do it
func (p *IntCodeProgram) ExecuteNextInstruction() error {
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
func (p *IntCodeProgram) ExecuteRelativeBaseOffset() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+2)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)

	p.relativeBase += firstParam

	p.instrPtr += 2
}

// ExecuteEquals stores 1 in third if first == second else 0
func (p *IntCodeProgram) ExecuteEquals() {
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
func (p *IntCodeProgram) ExecuteLessThan() {
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
func (p *IntCodeProgram) ExecuteJumpIfTrue() {
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
func (p *IntCodeProgram) ExecuteJumpIfFalse() {
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
func (p *IntCodeProgram) ExecuteInput() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+2)
	paramModes := getParamModes(inst[0])

	dest := p.resolveDestination(0, inst, paramModes)

	p.SetMemory(dest, <-p.input)
	p.instrPtr += 2
}

// ExecuteOutput simulate a print
func (p *IntCodeProgram) ExecuteOutput() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+2)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)

	p.output <- firstParam
	p.instrPtr += 2
}

// ExecuteAdd handles addition opcode
func (p *IntCodeProgram) ExecuteAdd() {
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
func (p *IntCodeProgram) ExecuteMultiply() {
	inst := p.MemorySlice(p.instrPtr, p.instrPtr+4)
	paramModes := getParamModes(inst[0])

	firstParam := p.resolveParam(0, inst, paramModes)
	secondParam := p.resolveParam(1, inst, paramModes)
	dest := p.resolveDestination(2, inst, paramModes)

	p.SetMemory(dest, firstParam*secondParam)
	p.instrPtr += 4
}

func (p *IntCodeProgram) resolveParam(i int, instruction []int, modes map[int]int) int {
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

func (p *IntCodeProgram) resolveDestination(i int, instruction []int, modes map[int]int) int {
	dest := instruction[i+1]
	if modes[i] == 2 {
		dest += p.relativeBase
	}
	return dest
}
