package fft

import (
	"adventofcode2019/common"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

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

// Result returns a string containing the first n digits
// n is clamped to fit the max length of result
func (impl *Impl) Result(n int) string {
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
func (impl *Impl) ProcessNSteps(n int) string {
	for i := 0; i < n; i++ {
		color.Cyan("iteration: %v\n", i)
		impl.processStep()
	}
	return impl.Result(8)
}

func (impl *Impl) processStep() {
	// assuming offset is large enough to get only zeros and ones from the pattern
	for i := impl.Size() - 2; i >= impl.offset; i-- {
		digit := common.AbsInt(impl.input[i+1]+impl.input[i]) % 10
		impl.input[i] = digit
	}
}
