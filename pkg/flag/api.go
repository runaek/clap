// Package flag provides helper functions for creating clap.FlagArg(s) with common types.
package flag

import (
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

func String(v *string, name string, opts ...clap.Option) *clap.FlagArg[string] {
	return clap.NewFlag[string](v, name, parse.String{}, opts...)
}

func Strings(v *[]string, name string, opts ...clap.Option) *clap.FlagArg[[]string] {
	return clap.NewFlags[string](v, name, parse.String{}, opts...)
}

func Bool(v *bool, name string, opts ...clap.Option) *clap.FlagArg[bool] {
	return clap.NewFlag[bool](v, name, parse.Bool{}, opts...)
}

func Int(v *int, name string, opts ...clap.Option) *clap.FlagArg[int] {
	return clap.NewFlag[int](v, name, parse.Int{}, opts...)
}

func Ints(v *[]int, name string, opts ...clap.Option) *clap.FlagArg[[]int] {
	return clap.NewFlags[int](v, name, parse.Int{}, opts...)
}

func Counter(v *parse.C, name string, opts ...clap.Option) *clap.FlagArg[parse.C] {
	return clap.NewFlag[parse.C](v, name, parse.Counter{}, opts...)
}

func Indicator(v *parse.I, name string, opts ...clap.Option) *clap.FlagArg[parse.I] {
	return clap.NewFlag[parse.I](v, name, parse.Indicator{}, opts...)
}
