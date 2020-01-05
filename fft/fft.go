package fft

import (
	"adventofcode2019/common"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

var pattern = [4]int{0, 1, 0, -1}

// Impl contains all details for Flawed Frequency Transmission
type Impl struct {
	input  []int
	offset int
}

// New creates a new instance of FFT
func New(line string) (*Impl, error) {
	impl := &Impl{}
	err := impl.Parse(line)
	if err != nil {
		return nil, err
	}
	offset, err := strconv.Atoi(line[0:7])
	if err != nil {
		return nil, err
	}
	impl.offset = offset

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

// String just to print some information about the data
func (impl *Impl) String() string {
	return fmt.Sprintf("{ len: %v, offset: %v}", impl.Size(), impl.offset)
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

// StringOfFirstNWithOffset returns a string containing the first n digits
// n is clamped to fit the max length of result
func (impl *Impl) StringOfFirstNWithOffset(n int) string {
	size := common.MinInt(len(impl.input)-impl.offset, n)
	var builder strings.Builder
	for _, digit := range impl.input[impl.offset : impl.offset+size] {
		builder.WriteString(strconv.Itoa(digit))
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
		color.Cyan("iteration: %v\n", i)
		impl.processStep()
	}
}

func (impl *Impl) processStep() {
	res := make([]int, impl.Size())
	for i := impl.offset; i < impl.Size(); i++ {
		computationIdx := i
		s := 0
		for idx := impl.offset; idx < impl.Size(); idx++ {
			coef := coefFromPattern(computationIdx, idx)
			s += impl.input[idx] * coef
		}
		digit := common.AbsInt(s) % 10
		res[i] = digit
	}
	copy(impl.input, res)
}

func coefFromPattern(computationIdx, idx int) int {
	return pattern[((idx+1)/(computationIdx+1))%4]
}
