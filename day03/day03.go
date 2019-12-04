package day03

import (
	"adventofcode2019/common"
	"bufio"
	"strconv"
	"strings"
)

// Run is the entrypoint of day03 exercice
func Run(filepath string) (int, error) {
	f := common.OpenFile(filepath)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)

	// scan the first line
	s.Scan()
	path1 := strings.Split(s.Text(), ",")
	err := s.Err()
	if err != nil {
		return 0, err
	}

	// scan the second line
	s.Scan()
	path2 := strings.Split(s.Text(), ",")
	err = s.Err()
	if err != nil {
		return 0, err
	}

	// transform strings to move type
	moves1, err := toMoves(path1)
	if err != nil {
		return 0, err
	}
	moves2, err := toMoves(path2)
	if err != nil {
		return 0, err
	}

	// create the 2 wires
	wire1 := Wire{}
	wire2 := Wire{}

	// apply moves to each wire
	wire1.UnwindFollowing(moves1)
	wire2.UnwindFollowing(moves2)

	// get the intersection
	crosspoints := wire1.Intersect(wire2)

	// distance is always positive so use a negative initial value
	minDistance := -1
	for p := range crosspoints.items {

		//distance := p.Part1DistanceComputation()
		distance := p.Part2DistanceComputation(wire1, wire2)

		if distance < minDistance || minDistance == -1 {
			minDistance = distance
		}
	}
	return minDistance, nil
}

//Part1DistanceComputation handles distance computation for Part 1
func (p *Point) Part1DistanceComputation() int {
	return common.AbsInt(p.x) + common.AbsInt(p.y)
}

//Part2DistanceComputation handles distance computation for Part 2
func (p *Point) Part2DistanceComputation(w1, w2 Wire) int {
	return w1.stepsToPoint[*p] + w2.stepsToPoint[*p]
}

type move struct {
	direction byte
	length    int
}

func toMoves(path []string) ([]move, error) {
	moves := make([]move, len(path))
	for i, e := range path {
		m, err := parseMove(e)
		if err != nil {
			return nil, err
		}
		moves[i] = m
	}
	return moves, nil
}

func parseMove(code string) (move, error) {
	direction := code[0]
	length, err := strconv.Atoi(code[1:])
	if err != nil {
		return move{}, err
	}
	return move{direction, length}, nil
}

// Point represents a point ;-)
type Point struct {
	x int
	y int
}

// PointSet represents a set of points
type PointSet struct {
	items map[Point]bool
}

// Add a Point to the set
func (s *PointSet) Add(p Point) {
	if s.items == nil {
		s.items = make(map[Point]bool)
	}

	_, found := s.items[p]
	if !found {
		s.items[p] = true
	}
}

// Contains returns true if p is in the set
func (s *PointSet) Contains(p Point) bool {
	_, found := s.items[p]
	return found
}

// Intersect returns intersection between two pointSets
func (s *PointSet) Intersect(other PointSet) PointSet {
	result := PointSet{}
	for item := range other.items {
		if s.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Wire represents a wire :-P
type Wire struct {
	actualPosition Point
	points         PointSet
	stepsToPoint   map[Point]int
	currentStep    int
}

// UnwindFollowing makes the wire follow a sequence of moves
func (w *Wire) UnwindFollowing(moves []move) {
	if w.stepsToPoint == nil {
		w.stepsToPoint = make(map[Point]int)
	}
	for _, m := range moves {
		unitMove := unitMoveActionToDirection(m.direction)
		w.Apply(unitMove, m.length)
	}
}

func unitMoveActionToDirection(directionCode byte) func(*Wire) Point {
	switch directionCode {
	case 'U':
		return U1
	case 'D':
		return D1
	case 'L':
		return L1
	case 'R':
		return R1
	}
	return func(w *Wire) Point { return w.actualPosition }
}

// R1 return the destination point if the wire do R1
func R1(w *Wire) Point {
	return Point{x: w.actualPosition.x + 1, y: w.actualPosition.y}
}

// L1 return the destination point if the wire do L1
func L1(w *Wire) Point {
	return Point{x: w.actualPosition.x - 1, y: w.actualPosition.y}
}

// U1 return the destination point if the wire do U1
func U1(w *Wire) Point {
	return Point{x: w.actualPosition.x, y: w.actualPosition.y + 1}
}

// D1 return the destination point if the wire do D1
func D1(w *Wire) Point {
	return Point{x: w.actualPosition.x, y: w.actualPosition.y - 1}
}

// Apply n times a move represented by f
func (w *Wire) Apply(f func(*Wire) Point, n int) {
	for i := 0; i < n; i++ {
		nextPosition := f(w)
		w.currentStep++
		w.actualPosition = nextPosition
		if !w.points.Contains(nextPosition) {
			w.points.Add(nextPosition)
			w.stepsToPoint[nextPosition] = w.currentStep
		}
	}
}

// Intersect is for intersection of wires
func (w *Wire) Intersect(other Wire) PointSet {
	return w.points.Intersect(other.points)
}
