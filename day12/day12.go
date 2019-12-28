package day12

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"regexp"
	"strconv"

	"gonum.org/v1/gonum/stat/combin"
)

// Run is the entryPoint of this day12 module
func Run(fileName string, steps, describeEvery int) (int, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)

	ch := challenge{make([]*positionAndVelocity, 0)}
	for s.Scan() {
		line := s.Text()
		p, err := parsePoint3d(line)
		if err != nil {
			return 0, err
		}
		ch.moons = append(ch.moons, &positionAndVelocity{pos: p})
	}
	err := s.Err()
	if err != nil {
		return 0, nil
	}

	describe(ch, 0)
	part1(ch, steps, describeEvery)

	totalEnergy := ch.totalEnergy()
	fmt.Printf("Total energy after %v steps: %v\n", steps, totalEnergy)

	return totalEnergy, nil
}

func part1(ch challenge, steps, describeEvery int) {
	for i := 1; i <= steps; i++ {
		ch.computeVelocities()
		ch.applyVelocities()
		if i%describeEvery == 0 {
			describe(ch, i)
		}
	}
}

func describe(chal challenge, step int) {
	fmt.Printf("After %v steps:\n", step)
	for _, pv := range chal.moons {
		fmt.Println(pv)
	}
}

type challenge struct {
	moons []*positionAndVelocity
}

func (ch *challenge) computeVelocities() {
	gen := combin.NewCombinationGenerator(4, 2)
	for gen.Next() {
		pair := gen.Combination(nil)

		moon1 := ch.moons[pair[0]]
		moon2 := ch.moons[pair[1]]

		if moon1.pos.x > moon2.pos.x {
			moon1.vel.x--
			moon2.vel.x++
		} else if moon1.pos.x < moon2.pos.x {
			moon1.vel.x++
			moon2.vel.x--
		}

		if moon1.pos.y > moon2.pos.y {
			moon1.vel.y--
			moon2.vel.y++
		} else if moon1.pos.y < moon2.pos.y {
			moon1.vel.y++
			moon2.vel.y--
		}

		if moon1.pos.z > moon2.pos.z {
			moon1.vel.z--
			moon2.vel.z++
		} else if moon1.pos.z < moon2.pos.z {
			moon1.vel.z++
			moon2.vel.z--
		}
	}
}

func (ch *challenge) applyVelocities() {
	for _, moon := range ch.moons {
		moon.pos = point3d{x: moon.pos.x + moon.vel.x, y: moon.pos.y + moon.vel.y, z: moon.pos.z + moon.vel.z}
	}
}

func (ch *challenge) totalEnergy() int {
	fmt.Println("computing total energy")
	result := 0
	for _, pv := range ch.moons {
		result += pv.totalEnergy()
	}
	return result
}

type positionAndVelocity struct {
	pos point3d
	vel point3d
}

func (pv positionAndVelocity) String() string {
	return fmt.Sprintf("pos=%v, vel=%v", pv.pos, pv.vel)
}

func (pv positionAndVelocity) totalEnergy() int {
	potx := common.AbsInt(pv.pos.x)
	poty := common.AbsInt(pv.pos.y)
	potz := common.AbsInt(pv.pos.z)
	pot := potx + poty + potz

	kinx := common.AbsInt(pv.vel.x)
	kiny := common.AbsInt(pv.vel.y)
	kinz := common.AbsInt(pv.vel.z)
	kin := kinx + kiny + kinz

	result := pot * kin

	fmt.Printf("pot: %2v + %2v + %2v = %3v;  kin: %2v + %2v + %2v = %3v;  total: %3v * %3v = %6v\n",
		potx, poty, potz, pot,
		kinx, kiny, kinz, kin,
		pot, kin, result)

	return result
}

type point3d struct {
	x int
	y int
	z int
}

func (p point3d) String() string {
	return fmt.Sprintf("<x=%v, y=%v, z=%v>", p.x, p.y, p.z)
}

var point3dRegexp = regexp.MustCompile(`<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`)

func parsePoint3d(text string) (point3d, error) {
	matches := point3dRegexp.FindStringSubmatch(text)
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		return point3d{}, err
	}
	y, err := strconv.Atoi(matches[2])
	if err != nil {
		return point3d{}, err
	}
	z, err := strconv.Atoi(matches[3])
	if err != nil {
		return point3d{}, err
	}
	return point3d{x, y, z}, nil
}
