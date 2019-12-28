package day10

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
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

	asteroid := game.guess200thDestroyedAsteroid()

	return asteroid.x, asteroid.y, nil
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

func (c *challenge) findBestAsteroidForStation() (Coord, int, int) {
	result := Coord{-1, -1}
	maxSeen := -1
	stationIdx := -1

	for idx, asteroid := range c.asteroids {
		nb := countAsteroids(idx, c.asteroids)
		if nb > maxSeen {
			maxSeen = nb
			result = asteroid
			stationIdx = idx
		}
	}

	fmt.Printf("seen asteroids: %v\n", maxSeen)
	return result, stationIdx, maxSeen
}

func countAsteroids(from int, asteroids []Coord) int {
	seenAsteroids := getAsteroidMapByAngle(from, asteroids)
	return len(seenAsteroids)
}

func getAsteroidMapByAngle(stationIdx int, asteroids []Coord) map[Coord][]Coord {
	result := make(map[Coord][]Coord)
	station := asteroids[stationIdx]

	for idx, asteroid := range asteroids {
		if stationIdx == idx {
			continue
		}

		direction := Coord{asteroid.x - station.x, asteroid.y - station.y}.reduce()
		s, found := result[direction]
		if !found {
			s = make([]Coord, 0)
		}
		result[direction] = append(s, asteroid)
	}

	return result
}

func (c *challenge) guess200thDestroyedAsteroid() Coord {
	station, stationIdx, _ := c.findBestAsteroidForStation()
	asteroidMap := getAsteroidMapByAngle(stationIdx, c.asteroids)

	keys := make([]Coord, len(asteroidMap))
	i := 0
	for k := range asteroidMap {
		keys[i] = k
		i++
	}
	fmt.Printf("station at %v\n", station)
	// sort keys of the map by angle with the laser direction
	// laser direction at start is up: (0, -1) as y points downward
	// when normalizing vectors of directions, angle with the laser direction is
	// = arccos( u-> . v->) = arccos( -u.y )
	sort.SliceStable(keys, func(i, j int) bool {
		key1 := normalize(keys[i])
		key2 := normalize(keys[j])

		var angleToReachAst1, angleToReachAst2 float64
		angleToReachAst1 = math.Acos(-float64(key1.y))
		if key1.x < 0 {
			angleToReachAst1 = 2*math.Pi - angleToReachAst1
		}
		angleToReachAst2 = math.Acos(-float64(key2.y))
		if key2.x < 0 {
			angleToReachAst2 = 2*math.Pi - angleToReachAst2
		}
		return angleToReachAst1 < angleToReachAst2
	})

	// if length of different angle is greater than 200, then only one pass should give you the 200th destroyed asteroid
	if len(asteroidMap) > 200 {
		alignedAsteroidsFor200thAngle := asteroidMap[keys[199]]
		// sort them by distance from the station
		sort.SliceStable(alignedAsteroidsFor200thAngle, func(i, j int) bool {
			asteroid1 := alignedAsteroidsFor200thAngle[i]
			asteroid2 := alignedAsteroidsFor200thAngle[j]
			ast1FromStation := math.Hypot(float64(asteroid1.x-station.x), float64(asteroid1.y-station.y))
			ast2FromStation := math.Hypot(float64(asteroid2.x-station.x), float64(asteroid2.y-station.y))
			return ast1FromStation < ast2FromStation
		})
		// and take the first one
		return alignedAsteroidsFor200thAngle[0]

	}
	// consume the closest asteroid from each angle, rince and repeat until 200 is here or no asteroid left
	// let's that computation for next time
	return Coord{}
}

// Coord represents coordinates
type Coord struct {
	x int
	y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%v, %v)", c.x, c.y)
}

func (c Coord) reduce() Coord {
	var result Coord
	if c.x == 0 && c.y != 0 {
		result = Coord{0, c.y / common.AbsInt(c.y)}
	} else if c.x != 0 && c.y == 0 {
		result = Coord{c.x / common.AbsInt(c.x), c.y}
	} else {
		gcd := common.Gcd(common.AbsInt(c.x), common.AbsInt(c.y))
		result = Coord{c.x / gcd, c.y / gcd}
	}
	return result
}

// NVector is a vector with a norm of 1
type NVector struct {
	x float64
	y float64
}

func normalize(c Coord) NVector {
	x := float64(c.x)
	y := float64(c.y)
	norm := math.Hypot(x, y)
	return NVector{x / norm, y / norm}
}
