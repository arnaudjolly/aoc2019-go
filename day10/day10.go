package day10

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"os"
)

// Run is the entrypoint of day10 exercice
func Run(filepath string) (int, int, error) {
	f := common.OpenFile(filepath)
	defer common.CloseFile(f)

	game := challenge{}
	err := game.loadFile(f)
	if err != nil {
		return 0, 0, nil
	}

	result := game.findBestAsteroid()
	return result.x, result.y, nil
}

type challenge struct {
	asteroids []Coord
}

func (c *challenge) loadFile(f *os.File) error {
	s := bufio.NewScanner(f)

	for j := 0; s.Scan(); j++ {
		line := s.Text()
		for i := 0; i < len(line); i++ {
			if line[i] == '#' {
				c.asteroids = append(c.asteroids, Coord{i, j})
			}
		}
	}
	return s.Err()
}

func (c *challenge) findBestAsteroid() Coord {
	result := Coord{-1, -1}
	maxSeen := -1

	for idx, asteroid := range c.asteroids {
		nb := countAsteroids(idx, c.asteroids)
		if nb > maxSeen {
			maxSeen = nb
			result = asteroid
		}
	}

	fmt.Printf("seen asteroids: %v\n", maxSeen)
	return result
}

func countAsteroids(from int, asteroids []Coord) int {
	seenAsteroids := make(map[Coord][]Coord)
	point := asteroids[from]

	for idx, asteroid := range asteroids {
		if from == idx {
			continue
		}

		direction := Coord{asteroid.x - point.x, asteroid.y - point.y}.normalize()
		s, found := seenAsteroids[direction]
		if !found {
			s = make([]Coord, 0)
		}
		seenAsteroids[direction] = append(s, asteroid)
	}

	return len(seenAsteroids)
}

// Coord represents coordinates
type Coord struct {
	x int
	y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%v, %v)", c.x, c.y)
}

func (c Coord) normalize() Coord {
	var result Coord
	if c.x == 0 && c.y != 0 {
		result = Coord{0, c.y / common.AbsInt(c.y)}
	} else if c.x != 0 && c.y == 0 {
		result = Coord{c.x / common.AbsInt(c.x), c.y}
	} else {
		gcd := gcd(common.AbsInt(c.x), common.AbsInt(c.y))
		result = Coord{c.x / gcd, c.y / gcd}
	}
	return result
}

func gcd(nb1 int, nb2 int) int {
	gcdnum := 1

	for i := 1; i <= nb1 && i <= nb2; i++ {
		if nb1%i == 0 && nb2%i == 0 {
			gcdnum = i
		}
	}
	return gcdnum
}
