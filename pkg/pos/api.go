// Package pos provides helper functions for creating clap.PositionalArg(s) with common types.
package pos

import (
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

func String(v *string, i int, opts ...clap.Option) *clap.PositionalArg[string] {
	return clap.NewPosition[string](v, i, parse.String{}, opts...)
}

func Strings(v *[]string, name int, opts ...clap.Option) *clap.PositionalArg[[]string] {
	return clap.NewPositions[string](v, name, parse.String{}, opts...)
}

func Int(v *int, name int, opts ...clap.Option) *clap.PositionalArg[int] {
	return clap.NewPosition[int](v, name, parse.Int{}, opts...)
}

func Ints(v *[]int, name int, opts ...clap.Option) *clap.PositionalArg[[]int] {
	return clap.NewPositions[int](v, name, parse.Int{}, opts...)
}
