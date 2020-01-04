package day15

import (
	"adventofcode2019/common"
	"adventofcode2019/intcode"
	"bufio"
	"fmt"
	"github.com/fatih/color"
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

	g := game{
		grid:  map[point]tile{origin: visited},
		droid: origin,
	}

	in := make(chan int)
	out := make(chan int)
	quit := make(chan int)
	go p.Run(in, out, quit)
	g.findOxygenSystem(make([]direction, 0), in, out)

	g.printGrid()

	return 0, nil
}

var (
	bgW = color.New(color.BgWhite, color.FgBlack)
	bgB = color.New(color.BgBlack, color.FgWhite)
	bgR = color.New(color.BgRed, color.FgBlack)
	bgG = color.New(color.BgGreen, color.FgWhite)
)

type game struct {
	grid        map[point]tile
	droid       point
	commandSent direction
}

func (g *game) tileAt(p point) tile {
	t, ok := g.grid[p]
	if !ok {
		t = unvisited
	}
	return t
}

func (g *game) findOxygenSystem(path []direction, in, out chan int) ([]direction, bool) {
	g.printGrid()
	if g.tileAt(g.droid.northTile()) == unvisited {
		result, found := g.tryInDirection(north, path, in, out)
		if found {
			return result, found
		}
	}
	if g.tileAt(g.droid.eastTile()) == unvisited {
		result, found := g.tryInDirection(east, path, in, out)
		if found {
			return result, found
		}
	}
	if g.tileAt(g.droid.westTile()) == unvisited {
		result, found := g.tryInDirection(west, path, in, out)
		if found {
			return result, found
		}
	}
	if g.tileAt(g.droid.southTile()) == unvisited {
		result, found := g.tryInDirection(south, path, in, out)
		if found {
			return result, found
		}
	}
	return path, false
}

func (g *game) tryInDirection(dir direction, path []direction, in, out chan int) ([]direction, bool) {
	fmt.Printf("trying to visit %v of %v\n", dir, g.droid)
	g.commandSent = dir
	in <- int(g.commandSent)
	t := tile(<-out)
	g.handle(t)
	fmt.Printf("sent:%v, received:%v, afterHandlePosition:%v\n", g.commandSent, t, g.droid)
	if t == wall {
		// finished, we hit a wall!
		return path, false
	}
	if t == oxygenSystem {
		// congrats, one path to the system has been discovered !
		return append(path, dir), true
	}

	// oh! a new undiscovered tile has been found, droid is already at good position
	// let's find the solution from this point
	dirSolution, dirFound := g.findOxygenSystem(make([]direction, 0), in, out)
	if dirFound {
		res := make([]direction, len(path)+len(dirSolution))
		copy(res, path)
		for idx, d := range dirSolution {
			res[len(path)+idx] = d
		}
		return res, true
	}
	// no solution found from the new point, droid should backtrack one step
	g.backtrack(dir, in, out)

	return path, false
}

func (g *game) handle(t tile) {
	fmt.Printf("received: %v for position:%v and command:%v\n", t, g.droid, g.commandSent)
	tileToMark := g.droid.destinationOfDirection(g.commandSent)
	switch t {
	case wall:
		g.grid[tileToMark] = wall
	case oxygenSystem:
		g.moveDroidTo(tileToMark)
		g.grid[tileToMark] = oxygenSystem
	case visited:
		g.moveDroidTo(tileToMark)
		g.grid[tileToMark] = visited
	}
}

func (g *game) backtrack(dir direction, in, out chan int) {
	fmt.Printf("backtracking %v...\n", dir)
	g.commandSent = dir.revert()
	in <- int(g.commandSent)
	g.handle(tile(<-out))
}

func (g *game) moveDroidTo(dest point) {
	g.droid = dest
}

func (g *game) markPointAs(p point, t tile) {
	g.grid[p] = t
}

func (g *game) printGrid() {
	fmt.Println("Grid:")
	minX, maxX := 1, -1
	minY, maxY := 1, -1

	for p := range g.grid {
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

	// ensure droid is on the map
	if g.droid.x < minX {
		minX = g.droid.x - 1
	}
	if g.droid.x > maxX {
		maxX = g.droid.x + 1
	}
	if g.droid.y < minY {
		minY = g.droid.y - 1
	}
	if g.droid.y > maxY {
		maxY = g.droid.y + 1
	}

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}
			tile, found := g.grid[p]
			if !found {
				tile = unvisited
			}
			if p == g.droid {
				bgR.Print("D")
			} else if p == origin {
				bgR.Print("S")
			} else {
				fmt.Print(tile)
			}
		}
		fmt.Printf("\n")
	}
}

type tile int

func (t tile) String() string {
	switch t {
	case wall:
		return bgW.Sprint("#")
	case visited:
		return bgW.Sprint(".")
	case unvisited:
		return bgB.Sprint("?")
	case oxygenSystem:
		return bgG.Sprint("X")
	default:
		return bgR.Sprintf("%1v", int(t))
	}
}

const (
	unvisited    = tile(-1)
	wall         = tile(0)
	visited      = tile(1)
	oxygenSystem = tile(2)
)

type direction int

func (d direction) String() string {
	switch d {
	case north:
		return "north"
	case west:
		return "west"
	case south:
		return "south"
	case east:
		return "east"
	default:
		return bgR.Sprintf("unknown(%v)", int(d))
	}
}

func (d direction) revert() direction {
	switch d {
	case north:
		return south
	case west:
		return east
	case south:
		return north
	case east:
		return west
	default:
		panic("reverting an unknown direction")
	}
}

func revertPath(path []direction) []direction {
	result := make([]direction, len(path))
	for i, d := range path {
		result[len(path)-1-i] = d.revert()
	}
	return result
}

const (
	north = direction(1)
	south = direction(2)
	west  = direction(3)
	east  = direction(4)
)

type point struct {
	x int
	y int
}

var origin = point{0, 0}

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

func (p point) destinationOfDirection(d direction) point {
	switch d {
	case north:
		return p.northTile()
	case south:
		return p.southTile()
	case east:
		return p.eastTile()
	case west:
		return p.westTile()
	default:
		panic("unknown direction " + d.String())
	}
}
