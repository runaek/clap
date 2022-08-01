// Package keyvalue provides helper functions for creating clap.KeyValueArg(s) with common types.
package keyvalue

import (
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

func String(v *string, name string, opts ...clap.Option) *clap.KeyValueArg[string] {
	return clap.NewKeyValue[string](v, name, parse.String{}, opts...)
}

func Strings(v *[]string, name string, opts ...clap.Option) *clap.KeyValueArg[[]string] {
	return clap.NewKeyValues[string](v, name, parse.String{}, opts...)
}

func Bool(v *bool, name string, opts ...clap.Option) *clap.KeyValueArg[bool] {
	return clap.NewKeyValue[bool](v, name, parse.Bool{}, opts...)
}

func Int(v *int, name string, opts ...clap.Option) *clap.KeyValueArg[int] {
	return clap.NewKeyValue[int](v, name, parse.Int{}, opts...)
}

func Ints(v *[]int, name string, opts ...clap.Option) *clap.KeyValueArg[[]int] {
	return clap.NewKeyValues[int](v, name, parse.Int{}, opts...)
}
