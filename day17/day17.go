package day17

import (
	"adventofcode2019/common"
	"adventofcode2019/intcode"
	"bufio"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Run is the entrypoint of day15 exercice
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

	// create a program instance to get the map and compute commands from it
	createProgram := intcode.ProgramCreator(seq)
	p := createProgram()
	out := make(chan int)
	quit := make(chan int)
	go p.Run(nil, out, quit)

	spacemap := SpaceMap{}
	spacemap.PopulateFrom(out)
	<-quit

	// prepare functions A B C
	cds := spacemap.robotCommands()
	fmt.Printf("Commands(%v):\n %v\n", len(cds.String()), cds)
	res := splitCommands(cds)
	fmt.Printf("Splitted:\nA:%v\nB:%v\nC:%v\nRoutine:%v\n", res.A, res.B, res.C, res.mainRoutine)

	// now create the real instance to send
	manual := createProgram()
	// override movement logic
	manual.SetMemory(0, 2)

	in2 := make(chan int)
	out2 := make(chan int)
	quit2 := make(chan int)
	go manual.Run(in2, out2, quit2)

	// stack commands to send to the program
	go send([]string{
		// Main:
		strings.Join(res.mainRoutine, ","),
		// Function A:
		res.A.String(),
		// Function B:
		res.B.String(),
		// Function C:
		res.C.String(),
		// Continuous video feed?
		"n",
	}, in2)

	// output everything the program send to us
	go output(out2)
	// make the main goroutine wait for termination of second program to stop
	<-quit2

	return 0, nil
}

func output(out chan int) {
	for c := range out {
		if c > 0xff {
			// score is greater than a byte
			fmt.Printf("Stardust: %v\n", c)
		} else {
			fmt.Print(string(c))
		}
	}
}

func send(strings []string, ch chan int) {
	// wait a bit to make printed line be readable ;)
	time.Sleep(50 * time.Millisecond)
	for _, str := range strings {
		fmt.Printf("%v\n", str)
		for _, c := range str {
			ch <- int(c)
		}
		ch <- int('\n')
		// wait a bit to make printed line be readable ;)
		time.Sleep(50 * time.Millisecond)
	}
}

type splitterResult struct {
	cds         commands
	mainRoutine []string
	A           commands
	B           commands
	C           commands
}

// brute force B and C based on length of A
func splitCommands(cds commands) splitterResult {
	found := false
	// begin with A and B of size 1 and 0
	sizeA, sizeB := 1, 1

	var result splitterResult

	for !found {
		// reset routine for this try
		result.mainRoutine = make([]string, 0)

		a := cds[:sizeA]
		result.mainRoutine = append(result.mainRoutine, "A")
		rest := cds[sizeA:]

		b := rest[:sizeB]
		if len(b) != 0 {
			result.mainRoutine = append(result.mainRoutine, "B")
			rest = rest[sizeB:]
		}

		// if b is too heavy to be stored in a function: try with a bigger A
		if len(b.String()) > 20 {
			sizeA++
			sizeB = 1
			continue
		}

		// but if a is too heavy... panic because you simply miss the solution or algo is not ok
		if len(a.String()) > 20 {
			panic("not found!!!!!!")
		}

		// repeat consuming each A and B until what remains has same length as before
		for previous := 0; len(rest) != previous; previous = len(rest) {
			// consume all A
			for reflect.DeepEqual(a, rest[:sizeA]) {
				result.mainRoutine = append(result.mainRoutine, "A")
				rest = rest[sizeA:]
			}
			// consume all B
			for reflect.DeepEqual(b, rest[:sizeB]) {
				result.mainRoutine = append(result.mainRoutine, "B")
				rest = rest[sizeB:]
			}
		}

		// with this A and B set
		// try to add successive elements to build C
		var c commands
		for _, item := range rest {
			c = append(c, item)
			// stop here if c is too heavy to be stored in a function
			if len(c.String()) > 20 {
				break
			}

			// check if the rest can be written in sequences of A, B and C
			if routine, ok := rest.composedOf(a, b, c); ok {
				result.mainRoutine = append(result.mainRoutine, routine...)
				result.A = a
				result.B = b
				result.C = c
				found = true
				break
			}
		}
		// next run, try with a bigger B
		sizeB++
	}

	return result
}

type orientation byte

func (o orientation) String() string {
	return string(byte(o))
}

const (
	left  = orientation('L')
	right = orientation('R')
)

type instruction struct {
	turn   orientation
	length int
}

func (i instruction) String() string {
	return fmt.Sprintf("%v,%v", i.turn, i.length)
}

// SpaceMap contains all information on the map
type SpaceMap struct {
	grid          map[point]tile
	width         int
	height        int
	robotPosition point
}

// PopulateFrom populates the space map following information given by out values
func (sm *SpaceMap) PopulateFrom(out chan int) {
	grid := make(map[point]tile)
	x, y := 0, 0
	maxX, maxY := -1, -1
	lastSeen := 0
	for d := range out {
		switch d {
		case '\n':
			if lastSeen == '\n' {
				break
			}
			x, y = 0, y+1
			break
		default:
			maxX = common.MaxInt(maxX, x)
			maxY = common.MaxInt(maxY, y)
			p := point{x, y}
			t := tile(d)
			grid[p] = t
			if t.isRobot() {
				sm.robotPosition = p
			}
			x++
		}
		lastSeen = d
	}
	sm.grid = grid
	sm.width = maxX + 1
	sm.height = maxY + 1
}

// Print allows to print the map
func (sm *SpaceMap) Print() {
	var strb strings.Builder
	strb.WriteString(fmt.Sprintf("Grid (width:%v, height:%v):\n", sm.width, sm.height))
	for y := 0; y < sm.height; y++ {
		for x := 0; x < sm.width; x++ {
			t := sm.grid[point{x, y}]
			strb.WriteByte(byte(t))
		}
		strb.WriteByte('\n')
	}
	fmt.Print(strb.String())
}

// SumAlignmentParams answers part1 of problem
func (sm *SpaceMap) SumAlignmentParams() int {
	sum := 0
	for pt, t := range sm.grid {
		// intersection is a scaffold
		if t != scaffold {
			continue
		}

		// intersection can't be at edges
		if pt.x == 0 || pt.x == sm.width-1 || pt.y == 0 || pt.y == sm.height-1 {
			continue
		}

		if sm.grid[pt.northTile()] == scaffold &&
			sm.grid[pt.southTile()] == scaffold &&
			sm.grid[pt.eastTile()] == scaffold &&
			sm.grid[pt.westTile()] == scaffold {
			sum += pt.x * pt.y
		}
	}
	return sum
}

func (sm *SpaceMap) tileInDirection(f func(point, tile) point) func(point, tile) tile {
	return func(p point, t tile) tile {
		return sm.grid[f(p, t)]
	}
}

func (sm *SpaceMap) robotCommands() commands {
	robot := sm.robotPosition
	robotTile := sm.grid[robot]

	f := sm.tileInDirection(frontPoint)
	l := sm.tileInDirection(leftPoint)
	r := sm.tileInDirection(rightPoint)

	frontTile := f(robot, robotTile)
	leftTile := l(robot, robotTile)
	rightTile := r(robot, robotTile)

	if frontTile == scaffold || leftTile != scaffold && rightTile != scaffold {
		panic("map problem: can't find an instruction beginning with turn left or right")
	}

	result := commands(make([]instruction, 0))
	for !(leftTile == space && rightTile == space) {

		// find orientation
		instr := instruction{}
		if leftTile == scaffold {
			instr.turn = left
		} else {
			instr.turn = right
		}

		// advance
		robotTile = robotTile.turn(instr.turn)
		sm.grid[robot] = scaffold
		for f(robot, robotTile) == scaffold {
			instr.length++
			robot = frontPoint(robot, robotTile)
		}
		leftTile = l(robot, robotTile)
		rightTile = r(robot, robotTile)
		result = append(result, instr)
	}

	return result
}

type commands []instruction

func (cds commands) String() string {
	var strb strings.Builder
	for i, instr := range cds {
		if i != 0 {
			strb.WriteByte(',')
		}
		strb.WriteString(instr.String())
	}
	return strb.String()
}

func (cds commands) composedOf(A, B, C commands) ([]string, bool) {
	result := make([]string, 0)
	if len(cds) == 0 {
		return result, true
	}
	matchA := len(A) > 0 && reflect.DeepEqual(A, cds[:len(A)])
	if matchA {
		if solution, ok := cds[len(A):].composedOf(A, B, C); ok {
			result = append(result, "A")
			if len(solution) > 0 {
				result = append(result, solution...)
			}
			return result, true
		}
		return nil, false
	}
	matchB := len(B) > 0 && reflect.DeepEqual(B, cds[:len(B)])
	if matchB {
		if solution, ok := cds[len(B):].composedOf(A, B, C); ok {
			result = append(result, "B")
			if len(solution) > 0 {
				result = append(result, solution...)
			}
			return result, true
		}
		return nil, false
	}
	matchC := len(C) > 0 && reflect.DeepEqual(C, cds[:len(C)])
	if matchC {
		if solution, ok := cds[len(C):].composedOf(A, B, C); ok {
			result = append(result, "C")
			if len(solution) > 0 {
				result = append(result, solution...)
			}
			return result, true
		}
		return nil, false
	}
	return nil, false
}

type tile byte

const (
	space     = tile('.')
	scaffold  = tile('#')
	faceUp    = tile('^')
	faceLeft  = tile('<')
	faceRight = tile('>')
	faceDown  = tile('v')
)

func (t tile) isRobot() bool {
	return t == faceDown || t == faceLeft || t == faceRight || t == faceUp
}

func (t tile) turn(o orientation) tile {
	switch t {
	case faceUp:
		if o == left {
			return faceLeft
		}
		return faceRight
	case faceDown:
		if o == left {
			return faceRight
		}
		return faceLeft
	case faceLeft:
		if o == left {
			return faceDown
		}
		return faceUp
	case faceRight:
		if o == left {
			return faceUp
		}
		return faceDown
	default:
		return t
	}
}

func frontPoint(pt point, t tile) point {
	if !t.isRobot() {
		panic("not a robot tile")
	}

	switch t {
	case faceUp:
		return pt.northTile()
	case faceDown:
		return pt.southTile()
	case faceLeft:
		return pt.westTile()
	case faceRight:
		return pt.eastTile()
	}

	panic("messed up with tiles !")
}

func leftPoint(pt point, t tile) point {
	return frontPoint(pt, t.turn(left))
}

func rightPoint(pt point, t tile) point {
	return frontPoint(pt, t.turn(right))
}

type point struct {
	x int
	y int
}

func (p point) String() string {
	return fmt.Sprintf("(%v %v)", p.x, p.y)
}

func (p point) northTile() point {
	return point{x: p.x, y: p.y - 1}
}
func (p point) southTile() point {
	return point{x: p.x, y: p.y + 1}
}
func (p point) eastTile() point {
	return point{x: p.x + 1, y: p.y}
}
func (p point) westTile() point {
	return point{x: p.x - 1, y: p.y}
}
