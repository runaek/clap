// Code generated by cmd/generate_drivers. DO NOT EDIT

package derive

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

const (
	ErrFloat clap.Error = "unable to derive 'float' Argument"
)

type floatDeriver struct{}

func (_ floatDeriver) DeriveKeyValue(a any, s string, opts ...clap.Option) (clap.IKeyValue, error) {

	v, ok := a.(*float64)

	if !ok {
		return nil, fmt.Errorf("%w: want *int but got %T", ErrFloat, v)
	}

	return clap.NewKeyValue[float64](v, s, parse.Float64{}, opts...), nil
}

func (_ floatDeriver) DerivePosition(a any, s int, opts ...clap.Option) (clap.IPositional, error) {

	v, ok := a.(*float64)

	if !ok {
		return nil, fmt.Errorf("%w: want *int but got %T", ErrFloat, v)
	}

	return clap.NewPosition[float64](v, s, parse.Float64{}, opts...), nil
}

func (_ floatDeriver) DeriveFlag(a any, s string, opts ...clap.Option) (clap.IFlag, error) {

	v, ok := a.(*float64)

	if !ok {
		return nil, fmt.Errorf("%w: want *int but got %T", ErrFloat, v)
	}

	return clap.NewFlag[float64](v, s, parse.Float64{}, opts...), nil
}

func init() {
	clap.RegisterFlagDeriver("float", floatDeriver{})
	clap.RegisterPositionalDeriver("float", floatDeriver{})
	clap.RegisterKeyValueDeriver("float", floatDeriver{})
}
