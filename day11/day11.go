package day11

import (
	"adventofcode2019/common"
	"bufio"
	"errors"
	"fmt"
	col "github.com/fatih/color"
	"strconv"
	"strings"
)

// Run is the entrypoint of day11 exercice
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

	runSampleMovements()

	createProgram := programCreator(seq)
	p := createProgram()
	in := make(chan int)
	out := make(chan int)

	c := challenge{grid: make(map[point]color)}

	go p.Run(in, out)

	for !p.halted {
		// fixme: sometimes it works, sometimes it panics with "all goroutines are asleep: deadlock!"
		// don't know why yet...
		// 1. send the color under the robot to the program
		in <- int(c.robotColor())
		// 2. read the first value output by the program as a color to paint the position
		newColor := color(<-out)
		// 3. read the second value output by the program as a direction to turn the robot
		direction := direction(<-out)

		// 4. act on the grid and robot
		c.paint(newColor)
		c.move(direction)
	}

	c.printGrid()
	return len(c.grid), nil
}

func runSampleMovements() {
	fmt.Println("running sample movements to see if this part is ok")
	c := challenge{grid: make(map[point]color)}
	c.paint(white)
	c.move(turnLeft)
	c.printGrid()

	c.paint(black)
	c.move(turnLeft)
	c.printGrid()

	c.paint(white)
	c.move(turnLeft)
	c.paint(white)
	c.move(turnLeft)
	c.printGrid()

	c.paint(black)
	c.move(turnRight)
	c.paint(white)
	c.move(turnLeft)
	c.paint(white)
	c.move(turnLeft)
	c.printGrid()
	fmt.Println("sample tests done.")
}

type challenge struct {
	robot robot
	grid  map[point]color
}

func (c *challenge) robotColor() color {
	result, found := c.grid[c.robot.position]
	if !found {
		result = black
	}
	return result
}

func (c *challenge) paint(aColor color) {
	c.grid[c.robot.position] = aColor
}
func (c *challenge) move(aDirection direction) {
	c.robot.turn(aDirection)
	c.robot.move()
}

var (
	bgW = col.New(col.BgWhite)
	bgB = col.New(col.BgBlack)
	bgR = col.New(col.BgRed)
)

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

	// ensure robot is on the grid to view it
	if c.robot.position.x > maxX {
		maxX = c.robot.position.x
	}
	if c.robot.position.x < minX {
		minX = c.robot.position.x
	}
	if c.robot.position.y > maxY {
		maxY = c.robot.position.y
	}
	if c.robot.position.y < minY {
		minY = c.robot.position.y
	}

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}
			colorIndicator, found := c.grid[p]
			if !found {
				colorIndicator = black
			}
			if c.robot.position == p {
				bgR.Print(c.robot.headingTo.String())
			} else if colorIndicator == black {
				bgB.Print(" ")
			} else if colorIndicator == white {
				bgW.Print(" ")
			}
		}
		fmt.Printf("\n")
	}
}

type robot struct {
	position  point
	headingTo orientation
}

func (r *robot) turn(direction direction) {
	if direction == turnLeft {
		r.headingTo = orientation((r.headingTo + 3) % 4)
	} else if direction == turnRight {
		r.headingTo = orientation((r.headingTo + 1) % 4)
	} else {
		panic(fmt.Sprintf("unknown direction: %v\n", direction))
	}
}

func (r *robot) move() {
	x, y := r.position.x, r.position.y
	switch r.headingTo {
	case up:
		r.position = point{x: x, y: y + 1}
	case down:
		r.position = point{x: x, y: y - 1}
	case left:
		r.position = point{x: x - 1, y: y}
	case right:
		r.position = point{x: x + 1, y: y}
	}
}

type color int

const (
	black = color(0)
	white = color(1)
)

type orientation int

const orientations = "^>v<"

func (o orientation) String() string {
	v := int(o)
	return orientations[v : v+1]
}

const (
	up    = orientation(0)
	right = orientation(1)
	down  = orientation(2)
	left  = orientation(3)
)

type direction int

const (
	turnLeft  = direction(0)
	turnRight = direction(1)
)

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
