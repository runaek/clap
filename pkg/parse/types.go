package parse

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"os"
	"strconv"
	"strings"
)

var (
	ErrEmptyInput = errors.New("no input to parse")
)

// Bool is a Parser[bool].
type Bool struct{}

func (_ Bool) Parse(input ...string) (bool, error) {
	if len(input) == 0 {
		return false, nil
	}

	b, err := strconv.ParseBool(input[0])

	if err != nil {
		err = fmt.Errorf("unable to parse bool from: %s", input[0])
	}

	return b, err
}

// String is a Parser[string].
type String struct{}

func (_ String) Parse(input ...string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	return input[0], nil
}

// Strings is a Parser[[]string].
type Strings struct{}

func (_ Strings) Parse(input ...string) ([]string, error) {
	return input, nil
}

// Int is a Parser[int].
type Int struct{}

func (_ Int) Parse(in ...string) (int, error) {
	if len(in) == 0 {
		return 0, ErrEmptyInput
	}

	i64, err := strconv.ParseInt(in[0], 10, 64)

	if err != nil {
		return 0, fmt.Errorf("unable to parse integer from: %s", in[0])
	}

	return int(i64), nil
}

// Ints is a Parser[[]int].
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

// Float64 is a Parser[float64]
type Float64 struct{}

func (_ Float64) Parse(input ...string) (float64, error) {
	if len(input) == 0 {
		return 0, ErrEmptyInput
	}

	if f, err := strconv.ParseFloat(input[0], 64); err != nil {
		return f, fmt.Errorf("unable to pase float64 from: %s", input[0])
	} else {
		return f, nil
	}

}

// C represents a 'counter' type.
//
// It is similar to a bool and an I - it does not expect a value to be supplied
// with the input. Its value is the number of times the input is detected.
//
// It can be parsed using a Counter.
type C int

// Counter is a parser for some C - it tracks how many times the argument was supplied rather than the value supplied.

// Counter is a Parser[C].
type Counter struct{}

func (_ Counter) Parse(input ...string) (C, error) {
	return C(len(input)), nil
}

// I represents an 'indicator' type.
//
// It is similar to a bool and a C - it does not expect a value to be supplied
// with the input. Each time the input is detected, its value is flipped.
//
// It can be parsed using an Indicator.
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

// File is a Parser[*os.File].
//
// It expects to be passed a filename, and it will attempt to create and return
// a file handle when parsed.
type File struct {
	// (Optional) the permissions to open the file with
	Mode *os.FileMode
	// The mode of operation
	Flag *int
}

func (f File) Parse(input ...string) (*os.File, error) {
	if len(input) == 0 {
		return nil, errors.New("unable to parse file: missing filepath")
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

// Map is a Parser[map[K]V].
type Map[K comparable, V any] struct {
	KeyParser   Parser[K]
	ValueParser Parser[V]
}

func (m Map[K, V]) Parse(input ...string) (map[K]V, error) {

	out := map[K]V{}

	for _, i := range input {
		parts := strings.SplitN(i, "=", 2)

		if len(parts) != 2 {
			return out, fmt.Errorf("invalid syntax: expected <key>=<value>: %s", i)
		}

		key, keyErr := m.KeyParser.Parse(parts[0])

		if keyErr != nil {
			return out, fmt.Errorf("unable to parse key: %w", keyErr)
		}

		val, valErr := m.ValueParser.Parse(parts[1])

		if valErr != nil {
			return out, fmt.Errorf("unable to parse value for key %v: %w", key, valErr)
		}

		out[key] = val
	}
	return out, nil
}

// StringMap is a Parser[string, V].
type StringMap[V any] Map[string, V]

func (s StringMap[V]) Parse(input ...string) (map[string]V, error) {
	return Map[string, V](s).Parse(input...)
}

var (
	_ Parser[map[int]any]    = Map[int, any]{}
	_ Parser[map[string]any] = StringMap[any]{}

	_ Parser[C]        = Counter{}
	_ Parser[I]        = Indicator{}
	_ Parser[int]      = Int{}
	_ Parser[bool]     = Bool{}
	_ Parser[float64]  = Float64{}
	_ Parser[string]   = String{}
	_ Parser[[]string] = Strings{}
)
