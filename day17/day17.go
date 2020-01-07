package day17

import (
	"adventofcode2019/common"
	"adventofcode2019/intcode"
	"bufio"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

	createProgram := intcode.ProgramCreator(seq)
	p := createProgram()
	// override movement logic
	// p.SetMemory(0, 2)
	// in := make(chan int)
	out := make(chan int)
	quit := make(chan int)
	go p.Run(nil, out, quit)

	spacemap := SpaceMap{}
	spacemap.PopulateFrom(out)
	spacemap.Print()

	// prepare functions A B C
	cds := spacemap.robotCommands()
	fmt.Printf("Commands(%v):\n %v\n", len(cds.String()), cds)

	res := splitCommands(cds)
	fmt.Printf("Splitted:\nA:%v\nB:%v\nC:%v\nRoutine:%v\n", res.A, res.B, res.C, res.mainRoutine)

	return spacemap.SumAlignmentParams(), nil
}

type splitterResult struct {
	cds         commands
	mainRoutine []string
	A           commands
	B           commands
	C           commands
}

func splitCommands(cds commands) splitterResult {
	found := false
	sizeA, sizeB := 1, 1

	var result splitterResult

	for !found {
		result.mainRoutine = make([]string, 0)
		a := cds[:sizeA]
		result.mainRoutine = append(result.mainRoutine, "A")
		rest := cds[sizeA:]

		var b commands
		if sizeA == sizeB {
			// consume each A in front of rest
			// and start B just after
			for reflect.DeepEqual(a, rest[:sizeA]) {
				result.mainRoutine = append(result.mainRoutine, "A")
				rest = rest[sizeA:]
			}
		}
		b = rest[:sizeB]
		result.mainRoutine = append(result.mainRoutine, "B")
		rest = rest[sizeB:]

		if len(b.String()) > 20 {
			sizeA++
			sizeB = 1
			continue
		}

		if len(a.String()) > 20 {
			panic("not found!!!!!!")
		}

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

		var c commands
		for _, item := range rest {
			c = append(c, item)
			if len(c.String()) > 20 {
				break
			}
			routine, ok := rest.composedOf(a, b, c)
			if ok {
				result.mainRoutine = append(result.mainRoutine, routine...)
				result.A = a
				result.B = b
				result.C = c
				found = true
				break
			}
		}
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
		subCmd := cds[len(A):]
		solution, ok := subCmd.composedOf(A, B, C)
		if ok {
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
		subCmd := cds[len(B):]
		solution, ok := subCmd.composedOf(A, B, C)
		if ok {
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
		subCmd := cds[len(C):]
		solution, ok := subCmd.composedOf(A, B, C)
		if ok {
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
