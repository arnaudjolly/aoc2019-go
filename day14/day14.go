package day14

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Run is the entryPoint of this day14 module
func Run(fileName string) (int, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	book := make([]formulae, 0)

	s := bufio.NewScanner(f)
	for s.Scan() {
		f, err := parseFormulae(s.Text())
		if err != nil {
			return 0, err
		}
		book = append(book, f)
	}
	err := s.Err()
	if err != nil {
		return 0, nil
	}

	fmt.Println("Book:")
	fmt.Println(book)

	return 0, nil
}

type formulae struct {
	inputs []dose
	output dose
}

var formulaeRegexp = regexp.MustCompile(`(.*) => (.*)`)

func parseFormulae(text string) (formulae, error) {
	matches := formulaeRegexp.FindStringSubmatch(text)

	input, err := parseDoseList(matches[1])
	if err != nil {
		return formulae{}, err
	}

	output, err := parseDose(matches[2])
	if err != nil {
		return formulae{}, err
	}
	f := formulae{input, output}

	return f, nil
}

func parseDoseList(line string) ([]dose, error) {
	result := make([]dose, 0)

	parts := strings.Split(line, ",")

	for _, p := range parts {
		d, err := parseDose(p)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

type dose struct {
	quantity  int
	component component
}

var doseRegexp = regexp.MustCompile(`(\d+) (\w+)`)

func parseDose(part string) (dose, error) {
	matches := doseRegexp.FindStringSubmatch(part)
	q, err := strconv.Atoi(matches[1])
	if err != nil {
		return dose{}, err
	}
	compo := component(matches[2])

	return dose{q, compo}, nil
}

type component string

const (
	fuel = component("FUEL")
	ore  = component("ORE")
)
