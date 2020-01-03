package day14

import (
	"adventofcode2019/common"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"regexp"
	"strconv"
	"strings"
)

// Run is the entryPoint of this day14 module
func Run(fileName string) (int, error) {
	f := common.OpenFile(fileName)
	defer common.CloseFile(f)

	book := book{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		f, err := parseFormulae(s.Text())
		if err != nil {
			return 0, err
		}
		book.formulas = append(book.formulas, f)
	}
	err := s.Err()
	if err != nil {
		return 0, nil
	}

	fmt.Printf("Book:\n%v\n", book)

	return book.howMuchOfThatToGetIngredients(ore, ingredients{fuel: 1})
}

type ingredients map[component]int

func (ings ingredients) Equals(other ingredients) bool {
	if len(ings) != len(other) {
		return false
	}

	for k, v := range ings {
		otherV, found := other[k]
		if !found || v != otherV {
			return false
		}
	}
	return true
}
func (ings ingredients) copy() ingredients {
	result := ingredients{}
	for comp, quantity := range ings {
		result[comp] = quantity
	}
	return result
}

type book struct {
	formulas        []formulae
	componentWeight map[component]int
}

func (bk book) String() string {
	var str strings.Builder
	for _, f := range bk.formulas {
		str.WriteString(f.String())
		str.WriteString("\n")
	}
	return str.String()
}

func (bk book) howMuchOfThatToGetIngredients(cmp component, needed ingredients) (int, error) {

	neededCopy := needed.copy()

	for comp := range needed {
		quantity := neededCopy[comp]
		if comp == cmp || quantity == 0 {
			// skip this component as it is the one we want
			continue
		}
		fmt.Printf("handling (%v, %v) %v\n", comp, quantity, neededCopy)

		// find the receipe that produces this component
		f, err := bk.getFormulaeProducing(comp)
		if err != nil {
			return 0, err
		}
		fmt.Printf("found formulae: %v\n", f)
		ratio, rmd := common.Eucl(quantity, f.output.quantity)
		// replace this component by its inputs in proportions
		if ratio != 0 {
			fmt.Printf("ratio:%v, rmd:%v\n", ratio, rmd)
			for _, input := range f.inputs {
				neededCopy[input.component] += input.quantity * ratio
			}
			newValue := neededCopy[comp] - quantity + rmd
			if newValue == 0 {
				delete(neededCopy, comp)
			} else {
				neededCopy[comp] = newValue
			}
		}
	}

	fmt.Printf("=> %v\n", neededCopy)
	if !neededCopy.Equals(needed) {
		// make another pass
		fmt.Println("make another pass!")
		return bk.howMuchOfThatToGetIngredients(cmp, neededCopy)
	}

	color.Red("deadend ! have to make a deal")
	// here we can't reduce it more without making a deal.
	// deal one component at a time
	deal, receipe := bk.findTheDeal(neededCopy)
	fmt.Printf("deal:%v, receipe:%v\n", deal, receipe)
	delete(neededCopy, deal)
	for _, i := range receipe.inputs {
		neededCopy[i.component] += i.quantity
	}

	fmt.Printf("after deal: %v\n", neededCopy)
	if len(neededCopy) == 1 {
		return neededCopy[cmp], nil
	}

	fmt.Println("after deal: try another pass")
	// make another pass
	return bk.howMuchOfThatToGetIngredients(cmp, neededCopy)
}

func (bk book) getFormulaeProducing(comp component) (formulae, error) {
	for _, f := range bk.formulas {
		if f.output.component == comp {
			return f, nil
		}
	}
	return formulae{}, fmt.Errorf("couldn't find a receipe to produce %v in book: %v", comp, bk)
}

// the deal is to accept using more component than necessary.
// I think the one that will free up a lot more other component is good choice
// to avoid multiple deals for the same component.
func (bk book) findTheDeal(ings ingredients) (component, formulae) {
	var deal component = component("")
	var dealFormulae formulae = formulae{}
	maxInput := -1

	for c := range ings {
		score := bk.weightOf(c)
		if score > maxInput {
			maxInput = score
			deal = c
			dealFormulae, _ = bk.getFormulaeProducing(c)
		}
	}
	return deal, dealFormulae
}

// I call this score: the weight of a component
// weight(ORE) = 0
// if 1A, 2B, 3C => 4D then weight(D) = 1 + weight(A)+ weight(B)+ weight(C)
// there are surely better methods to do this, but this one worked for me ;-)
func (bk book) weightOf(cmp component) int {
	if bk.componentWeight == nil {
		bk.componentWeight = make(map[component]int)
	}

	weight, found := bk.componentWeight[cmp]
	if found {
		return weight
	}

	if cmp == ore {
		// no step to produce this, ORE is free
		return 0
	}

	// else weight(output) = sum(weight(inputs)) + 1
	sumInputWeight := 0
	f, _ := bk.getFormulaeProducing(cmp)
	for _, d := range f.inputs {
		sumInputWeight += bk.weightOf(d.component)
	}

	result := 1 + sumInputWeight
	bk.componentWeight[cmp] = result
	return result
}

type formulae struct {
	inputs []dose
	output dose
}

func (f formulae) String() string {
	inputStr := make([]string, len(f.inputs))
	for idx, d := range f.inputs {
		inputStr[idx] = d.String()
	}
	return fmt.Sprintf("%v => %v", strings.Join(inputStr, ", "), f.output.String())
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

func (d dose) String() string {
	return fmt.Sprintf("%v %v", d.quantity, d.component)
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
