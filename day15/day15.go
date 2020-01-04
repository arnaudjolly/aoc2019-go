package day15

import (
	"adventofcode2019/common"
	"adventofcode2019/intcode"
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/RyanCarrier/dijkstra"
	"github.com/fatih/color"
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
	g.walkTheMap(in, out)
	g.printGrid()

	graph := g.buildGraph()

	// part1
	originID := graph.AddMappedVertex(origin.String())
	oxygenID := graph.AddMappedVertex(g.oxygenSystemPosition.String())
	best, err := graph.Shortest(originID, oxygenID)
	if err != nil {
		return 0, err
	}
	part1 := int(best.Distance)
	fmt.Printf("part1 result: %v\n", part1)

	// part2
	maxMinDistance, err := g.maxMinDistanceFrom(g.oxygenSystemPosition)
	if err != nil {
		return 0, err
	}
	return maxMinDistance, nil
}

var (
	blackOnCyan   = color.New(color.BgCyan, color.FgBlack)
	blackOnRed    = color.New(color.BgRed, color.FgBlack)
	blackOnYellow = color.New(color.BgYellow, color.FgBlack)
	blackOnWhite  = color.New(color.BgWhite, color.FgBlack)
	whiteOnBlack  = color.New(color.BgBlack, color.FgWhite)
	whiteOnGreen  = color.New(color.BgGreen, color.FgWhite)
)

type directionStack []direction

func (s directionStack) Empty() bool      { return len(s) == 0 }
func (s directionStack) Peek() direction  { return s[len(s)-1] }
func (s *directionStack) Put(i direction) { (*s) = append((*s), i) }
func (s *directionStack) Pop() direction {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}

type game struct {
	grid                 map[point]tile
	droid                point
	oxygenSystemPosition point
	commandSent          direction
}

func (g *game) buildGraph() *dijkstra.Graph {
	graph := dijkstra.NewGraph()

	minX, maxX := 1, -1
	minY, maxY := 1, -1
	for p := range g.grid {
		minX, maxX = common.MinInt(minX, p.x), common.MaxInt(maxX, p.x)
		minY, maxY = common.MinInt(minY, p.y), common.MaxInt(maxY, p.y)
	}

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}
			t := g.tileAt(p)

			if t.isInternalCell() {
				// add arcs for each connected cells
				for _, cell := range [4]point{
					p.northTile(),
					p.eastTile(),
					p.southTile(),
					p.westTile(),
				} {
					if g.tileAt(cell).isInternalCell() {
						graph.AddMappedArc(p.String(), cell.String(), 1)
					}
				}
			}
		}
	}

	return graph
}

func (g *game) maxMinDistanceFrom(pt point) (int, error) {
	graph := g.buildGraph()
	// for each internal cell, find the min distance to oxygenSystem
	// keep the max value for that
	max := 0
	ptID := graph.AddMappedVertex(pt.String())
	for cell, t := range g.grid {
		if t.isInternalCell() && cell != g.oxygenSystemPosition {
			cellID := graph.AddMappedVertex(cell.String())
			dist, err := graph.Shortest(cellID, ptID)
			if err != nil {
				g.markPointAs(pt, tile(9))
				g.printGrid()
				fmt.Println(graph)
				return 0, fmt.Errorf("from point %v: %v", pt, err)
			}
			max = common.MaxInt(max, int(dist.Distance))
		}
	}

	return max, nil
}

func (g *game) tileAt(p point) tile {
	t, ok := g.grid[p]
	if !ok {
		t = unvisited
	}
	return t
}

func (g *game) walkTheMap(in, out chan int) {
	dirs := g.directions()
	for !dirs.Empty() {
		dir := dirs.Pop()
		t := g.handleDirection(dir, in, out)
		if t != wall {
			g.walkTheMap(in, out)
			g.handleDirection(dir.revert(), in, out)
		}
	}
}

func (g *game) handleDirection(d direction, in, out chan int) tile {
	g.commandSent = d
	in <- int(g.commandSent)
	t := tile(<-out)
	g.handle(t)
	return t
}

func (g *game) directions() directionStack {
	var result directionStack
	if g.tileAt(g.droid.southTile()) == unvisited {
		result.Put(south)
	}
	if g.tileAt(g.droid.eastTile()) == unvisited {
		result.Put(east)
	}
	if g.tileAt(g.droid.westTile()) == unvisited {
		result.Put(west)
	}
	if g.tileAt(g.droid.northTile()) == unvisited {
		result.Put(north)
	}
	return result
}

func (g *game) handle(t tile) {
	tileToMark := g.droid.destinationOfDirection(g.commandSent)
	switch t {
	case wall:
		g.markPointAs(tileToMark, wall)
	case oxygenSystem:
		g.moveDroidTo(tileToMark)
		g.markPointAs(tileToMark, oxygenSystem)
	case visited:
		g.moveDroidTo(tileToMark)
		g.markPointAs(tileToMark, visited)
	}
}

func (g *game) moveDroidTo(dest point) {
	g.droid = dest
}

func (g *game) markPointAs(p point, t tile) {
	g.grid[p] = t
	if t == oxygenSystem {
		g.oxygenSystemPosition = p
	}
}

func (g *game) printGrid() {
	fmt.Println("Grid:")
	minX, maxX := 1, -1
	minY, maxY := 1, -1
	for p := range g.grid {
		minX, maxX = common.MinInt(minX, p.x), common.MaxInt(maxX, p.x)
		minY, maxY = common.MinInt(minY, p.y), common.MaxInt(maxY, p.y)
	}

	// ensure droid is on the map
	minX, maxX = common.MinInt(minX, g.droid.x), common.MaxInt(maxX, g.droid.x)
	minY, maxY = common.MinInt(minY, g.droid.y), common.MaxInt(maxY, g.droid.y)

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}
			tile, found := g.grid[p]
			if !found {
				tile = unvisited
			}
			col := g.colorOf(p)
			if p == g.droid {
				col.Print("D")
			} else if p == origin {
				col.Print("S")
			} else {
				col.Print(tile)
			}
		}
		fmt.Printf("\n")
	}
}

func (g *game) colorOf(p point) *color.Color {
	t := g.tileAt(p)
	switch {
	case t == oxygenSystem:
		return whiteOnGreen
	case p == origin || p == g.droid:
		return blackOnCyan
	case t == visited:
		return blackOnWhite
	case t == wall:
		return blackOnYellow
	default:
		return blackOnRed
	}
}

type tile int

func (t tile) String() string {
	switch t {
	case wall:
		return blackOnYellow.Sprint("#")
	case visited:
		return blackOnWhite.Sprint(" ")
	case unvisited:
		return whiteOnBlack.Sprint(" ")
	case oxygenSystem:
		return whiteOnGreen.Sprint("X")
	default:
		return blackOnRed.Sprintf("%1v", int(t))
	}
}

func (t tile) isInternalCell() bool {
	return t == visited || t == oxygenSystem
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
		return blackOnRed.Sprintf("unknown(%v)", int(d))
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

func (p point) String() string {
	return fmt.Sprintf("(%v %v)", p.x, p.y)
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
