package day17

import (
	"adventofcode2019/common"
	"adventofcode2019/intcode"
	"bufio"
	"fmt"
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
	out := make(chan int)
	quit := make(chan int)
	go p.Run(nil, out, quit)

	spacemap := SpaceMap{}
	spacemap.PopulateFrom(out)
	spacemap.Print()

	return spacemap.SumAlignmentParams(), nil
}

// SpaceMap contains all information on the map
type SpaceMap struct {
	grid   map[point]tile
	width  int
	height int
}

// PopulateFrom populates the space map following information given by out values
func (sm *SpaceMap) PopulateFrom(out chan int) {
	grid := make(map[point]tile)
	x, y := 0, 0
	maxX, maxY := -1, -1
	for d := range out {
		switch d {
		case '\n':
			x, y = 0, y+1
		default:
			maxX = common.MaxInt(maxX, x)
			maxY = common.MaxInt(maxY, y)
			p := point{x, y}
			grid[p] = tile(d)
			x++
		}
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

type tile byte

const (
	space    = tile('.')
	scaffold = tile('#')
)

type point struct {
	x int
	y int
}

func (p point) String() string {
	return fmt.Sprintf("(%v %v)", p.x, p.y)
}

func (p point) northTile() point {
	return point{x: p.x, y: p.y + 1}
}
func (p point) southTile() point {
	return point{x: p.x, y: p.y - 1}
}
func (p point) eastTile() point {
	return point{x: p.x + 1, y: p.y}
}
func (p point) westTile() point {
	return point{x: p.x - 1, y: p.y}
}
