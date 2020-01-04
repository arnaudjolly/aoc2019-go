package fft

import (
	"adventofcode2019/common"
	"fmt"
	"strconv"
	"strings"
)

// Impl contains all details for Flawed Frequency Transmission
type Impl struct {
	input              []int
	pattern            []int
	patternsForIndices [][]int
}

// New creates a new instance of FFT
func New(pattern []int, line string) (*Impl, error) {
	impl := &Impl{pattern: pattern}
	err := impl.Parse(line)
	if err != nil {
		return nil, err
	}
	return impl, nil
}

// Parse allows to initialize FFT with digits from line
func (impl *Impl) Parse(line string) error {
	impl.input = make([]int, len(line))
	for idx, str := range strings.Split(line, "") {
		i, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		impl.input[idx] = i
	}
	return nil
}

// StringOfFirstN returns a string containing the first n digits of input
// n is clamped to fit the max length of result
func (impl *Impl) StringOfFirstN(n int) string {
	size := common.MinInt(len(impl.input), n)
	var builder strings.Builder
	for i := 0; i < size; i++ {
		builder.WriteString(strconv.Itoa(impl.input[i]))
	}
	return builder.String()
}

// Size returns the input size
func (impl *Impl) Size() int {
	return len(impl.input)
}

// ProcessNSteps processes N step of phases
func (impl *Impl) ProcessNSteps(n int) {
	for i := 0; i < n; i++ {
		impl.processStep()
	}
}

func (impl *Impl) processStep() {
	if impl.patternsForIndices == nil {
		impl.computePatternsForIndices()
	}
	res := make([]int, impl.Size())
	for i := 0; i < impl.Size(); i++ {
		pat := patternFor(impl.pattern, i, impl.Size())
		s := scalar(impl.input, pat)
		res[i] = common.AbsInt(s) % 10
	}
	impl.input = res
}

func (impl *Impl) computePatternsForIndices() {
	impl.patternsForIndices = make([][]int, impl.Size())
	for i := 0; i < impl.Size(); i++ {
		impl.patternsForIndices[i] = patternFor(impl.pattern, i, impl.Size())
	}
}

func scalar(v1, v2 []int) int {
	if len(v1) != len(v2) {
		panic(fmt.Sprintf("scalar: lengths are differents v1(%v), v2(%v)", len(v1), len(v2)))
	}
	result := 0
	for idx, fromV1 := range v1 {
		result += fromV1 * v2[idx]
	}
	return result
}

func patternFor(basePattern []int, valueIdx int, size int) []int {
	// we will cut the very first one so add one to keep good size
	length := size + 1
	copies := valueIdx + 1

	q, _ := common.Eucl(length, len(basePattern)*copies)
	res := make([]int, (q+1)*len(basePattern)*copies)
	finalIdx := 0
	for loopOfBasePattern := 0; loopOfBasePattern < q+1; loopOfBasePattern++ {
		for _, patternElt := range basePattern {
			for repeat := 0; repeat < copies; repeat++ {
				res[finalIdx] = patternElt
				finalIdx = finalIdx + 1
			}
		}
	}

	// skip the first value and limit the size
	return res[1:length]
}
