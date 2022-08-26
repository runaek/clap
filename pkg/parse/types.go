package parse

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"os"
	"strconv"
)

var (
	ErrEmptyInput = errors.New("no input to parse")
)

type Bool struct{}

func (_ Bool) Parse(input ...string) (bool, error) {
	if len(input) == 0 {
		return false, nil
	}

	return strconv.ParseBool(input[0])
}

type String struct{}

func (_ String) Parse(input ...string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	return input[0], nil
}

type Strings struct{}

func (_ Strings) Parse(input ...string) ([]string, error) {
	return input, nil
}

type Int struct{}

func (_ Int) Parse(in ...string) (int, error) {
	if len(in) == 0 {
		return 0, ErrEmptyInput
	}

	i64, err := strconv.ParseInt(in[0], 10, 64)

	if err != nil {
		return 0, err
	}

	return int(i64), nil
}

type Ints struct{}

func (_ Ints) Parse(input ...string) ([]int, error) {
	out := make([]int, len(input))

	res := new(multierror.Error)

	ip := Int{}

	for i, in := range input {
		if iv, err := ip.Parse(in); err != nil {
			res = multierror.Append(res, err)
		} else {
			out[i] = iv
		}
	}

	return out, res.ErrorOrNil()
}

type Float64 struct{}

func (_ Float64) Parse(input ...string) (float64, error) {
	if len(input) == 0 {
		return 0, ErrEmptyInput
	}

	return strconv.ParseFloat(input[0], 64)
}

// C is a special 'counter' type that tracks the number of times a particular argument is supplied, expecting no value to
// be supplied - like a boolean flag.
type C int

// Counter is a parser for some C - it tracks how many times the argument was supplied rather than the value supplied.
type Counter struct{}

func (_ Counter) Parse(input ...string) (C, error) {
	return C(len(input)), nil
}

// I is a special 'indicator' that acts as a boolean switch that can be supplied multiple times.
type I bool

// Indicator is a parser for some I - it flips the value for each time the input is supplied/detected.
type Indicator struct {
	Initial bool
}

func (i Indicator) Parse(input ...string) (I, error) {
	active := i.Initial

	if len(input)%2 == 0 {
		return I(active), nil
	}
	return I(!active), nil
}

type File struct {
	Mode *os.FileMode
	Flag *int
}

func (f File) Parse(input ...string) (*os.File, error) {
	if len(input) == 0 {
		return nil, errors.New("missing filepath")
	}

	var (
		perm os.FileMode
		mode int
	)

	if f.Flag == nil {
		mode = os.O_CREATE | os.O_RDWR
	} else {
		mode = *f.Flag
	}

	if f.Mode == nil {
		perm = os.ModePerm
	} else {
		perm = *f.Mode
	}

	return os.OpenFile(input[0], mode, perm)
}
