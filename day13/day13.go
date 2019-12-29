package day13

import (
	"adventofcode2019/common"
	"adventofcode2019/intcode"
	"bufio"
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

	createProgram := intcode.ProgramCreator(seq)
	p := createProgram()
	// init quarters
	p.SetMemory(0, 2)

	g := game{grid: make(map[point]tile)}

	in := make(chan int)
	out := make(chan int)
	quit := make(chan int)
	go p.Run(in, out, quit)

	for run := true; run; {
		select {
		case x := <-out:
			y, info := <-out, <-out
			fmt.Printf("received: %v, %v, %v\n", x, y, tile(info))
			if x == -1 && y == 0 {
				// we receive the score once the entire board is loaded
				// so it starts the game, we can send the first joystick move
				// once loaded, don't send a move at score reception but when
				// receiving new position of ball
				if !g.loaded {
					go sendMove(g, in)
				}
				g.setScore(info)
			} else {
				t := tile(info)
				g.placeTile(x, y, t)
				// if we receive a ball position and the game is loaded...
				// send a move
				if t == ball && g.loaded {
					go sendMove(g, in)
				}
			}
		case <-quit:
			run = false
		}
	}

	g.printGrid()

	blocks := g.count(isBlock)

	return blocks, nil
}

func sendMove(g game, in chan int) {
	joystickMove := g.guessPaddleMove()
	fmt.Printf("Move: %v\n", joystickMove)
	in <- int(joystickMove)
}

var (
	bgW  = col.New(col.BgWhite, col.FgBlack)
	bgB  = col.New(col.BgBlack, col.FgWhite)
	bgR  = col.New(col.BgRed, col.FgBlack)
	blue = col.New(col.FgBlue)
)

type game struct {
	grid   map[point]tile
	score  int
	ball   point
	paddle point
	loaded bool
}

func (g *game) setScore(score int) {
	g.score = score
	g.loaded = true
}

func (g *game) placeTile(x, y int, t tile) {
	p := point{x, y}
	g.grid[p] = t
	if t == ball {
		g.ball = p
	}
	if t == hpaddle {
		g.paddle = p
	}
}

func (g *game) guessPaddleMove() joystick {
	if g.ball.x < g.paddle.x {
		return left
	} else if g.ball.x > g.paddle.x {
		return right
	}
	return neutral
}

func (g *game) count(tileFilter func(tile) bool) int {
	result := 0
	for _, t := range g.grid {
		if tileFilter(t) {
			result++
		}
	}
	return result
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

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := point{x, y}
			tile, found := g.grid[p]
			if !found {
				tile = empty
			}
			fmt.Print(tile)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Score: %v\n", blue.Sprint(g.score))
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

type joystick int

func (j joystick) String() string {
	switch j {
	case neutral:
		return "neutral"
	case left:
		return "left"
	case right:
		return "right"
	default:
		return bgR.Sprintf("unknown(%v)", int(j))
	}
}

const (
	neutral = joystick(0)
	left    = joystick(-1)
	right   = joystick(1)
)

type point struct {
	x int
	y int
}
