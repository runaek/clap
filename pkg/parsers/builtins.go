package parsers

import (
	"strconv"
)

// TODO: add more here

// Bool is a parser for a bool-value
func Bool(in ...string) (bool, error) {

	if len(in) == 0 {
		return false, nil
	}

	return strconv.ParseBool(in[0])
}

// String is a parser for a string-value
func String(in ...string) (string, error) {
	return in[0], nil
}

// Strings is a parser for a []string-value supplied as <sep> separated arguments
func Strings(in ...string) ([]string, error) {
	return in, nil
}

// Int is a parser for an int-value.
func Int(in ...string) (int, error) {
	i64, err := strconv.ParseInt(in[0], 10, 64)

	if err != nil {
		return 0, err
	}

	return int(i64), nil
}

// Int64 is a parser for an int64-value.
func Int64(in ...string) (int64, error) {
	return strconv.ParseInt(in[0], 10, 64)
}

// Float64 is a parser for a float64-value.
func Float64(in ...string) (float64, error) {
	return strconv.ParseFloat(in[0], 64)
}

// Float32 is a parser for a float32-value.
func Float32(in ...string) (float32, error) {
	f32, err := strconv.ParseFloat(in[0], 32)

	if err != nil {
		return 0, err
	}

	return float32(f32), nil
}
