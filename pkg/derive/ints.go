// Code generated by github.com/runaek/clap/cmd/generate_drivers. DO NOT EDIT

package derive

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

const (
	ErrInts clap.Error = "unable to derive 'ints' Argument"
)

type intsDeriver struct{}

func (_ intsDeriver) DeriveKeyValue(a any, s string, opts ...clap.Option) (clap.IKeyValue, error) {
	v, ok := a.(*[]int)

	if !ok {
		return nil, fmt.Errorf("%w: want *[]int but got %T", ErrInts, v)
	}

	return clap.NewKeyValue[[]int](v, s, parse.Ints{}, opts...), nil
}

func (_ intsDeriver) DerivePosition(a any, s int, opts ...clap.Option) (clap.IPositional, error) {
	v, ok := a.(*[]int)

	if !ok {
		return nil, fmt.Errorf("%w: want *[]int but got %T", ErrInts, v)
	}

	return clap.NewPosition[[]int](v, s, parse.Ints{}, opts...), nil
}

func (_ intsDeriver) DeriveFlag(a any, s string, opts ...clap.Option) (clap.IFlag, error) {
	v, ok := a.(*[]int)

	if !ok {
		return nil, fmt.Errorf("%w: want *[]int but got %T", ErrInts, v)
	}

	return clap.NewFlag[[]int](v, s, parse.Ints{}, opts...), nil
}

func init() {
	clap.RegisterFlagDeriver("ints", intsDeriver{})
	clap.RegisterPositionalDeriver("ints", intsDeriver{})
	clap.RegisterKeyValueDeriver("ints", intsDeriver{})
}
