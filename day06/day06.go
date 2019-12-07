package day06

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"strings"
)

// Run is the entrypoint of day06 exercice
func Run(filepath string) (int, error) {
	f := common.OpenFile(filepath)
	defer common.CloseFile(f)

	s := bufio.NewScanner(f)

	space := Space{}
	for s.Scan() {
		objects := strings.SplitN(s.Text(), ")", 2)
		first, second := objects[0], objects[1]

		// add second as a child of first
		obj1, found := space[first]
		if !found {
			obj1 = &SpaceObject{name: first, childrenNames: make([]string, 0)}
			space[first] = obj1
		}
		obj1.childrenNames = append(obj1.childrenNames, second)

		// add first as parent of obj2
		obj2, found := space[second]
		if !found {
			obj2 = &SpaceObject{name: second, childrenNames: make([]string, 0)}
			space[second] = obj2
		}
		obj2.parentName = first
	}
	err := s.Err()
	if err != nil {
		return 0, nil
	}

	result := part2(space)

	return result, nil
}

// compute number of orbital transfer we need to orbit directly around Santa is orbiting aroung
func part2(space Space) int {
	ours := space.parents("YOU")
	fmt.Printf("ours: %v\n", ours)
	his := space.parents("SAN")
	fmt.Printf("santa's: %v\n", his)

	var innerLoopSlice []string
	var outerLoopSlice []string
	if len(ours) < len(his) {
		innerLoopSlice = ours
		outerLoopSlice = his
	} else {
		innerLoopSlice = his
		outerLoopSlice = ours
	}

	for i := 0; i < len(outerLoopSlice); i++ {
		revIdx := len(outerLoopSlice) - 1 - i
		for j := 0; j < len(innerLoopSlice); j++ {
			revJdx := len(innerLoopSlice) - 1 - j
			if outerLoopSlice[revIdx] == innerLoopSlice[revJdx] {
				return i + j
			}
		}
	}

	return -1
}

func part1(space Space) int {
	direct := space.DirectOrbits()
	fmt.Println("direct: ", direct)

	indirect := space.IndirectOrbits()
	fmt.Println("indirect: ", indirect)

	return direct + indirect
}

// SpaceObject is information about an object
type SpaceObject struct {
	name          string
	parentName    string
	childrenNames []string
}

// Space contains all space objects indexed by name
type Space map[string]*SpaceObject

// DirectOrbits returns the number of direct orbits in the space
func (space Space) DirectOrbits() int {
	// no direct orbit for COM
	result := len(space) - 1
	if result < 0 {
		return 0
	}
	return result
}

// IndirectOrbits returns the number of indirect orbits in the space
func (space Space) IndirectOrbits() int {
	result := 0
	cache := make(map[string]int)
	for o := range space {
		temp := space.indirectOrbits(o, &cache)
		cache[o] = temp
		result += temp
	}
	return result
}

func (space Space) indirectOrbits(o string, cache *map[string]int) int {
	result, found := (*cache)[o]
	if found {
		return result
	}
	infos := space[o]
	if infos.parentName == "COM" || infos.parentName == "" {
		return 0
	}

	indirectOfParent := space.indirectOrbits(infos.parentName, cache)
	(*cache)[infos.parentName] = indirectOfParent
	return 1 + indirectOfParent
}

func (space Space) parents(o string) []string {
	infos := space[o]
	if infos.parentName == "" {
		return make([]string, 0)
	}

	return append(space.parents(infos.parentName), infos.parentName)
}
