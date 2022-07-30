// Code generated by github.com/runaek/clap/cmd/generate_drivers. DO NOT EDIT

package derive

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

const (
	ErrIndicator clap.Error = "unable to derive 'indicator' Argument"
)

type indicatorDeriver struct{}

func (_ indicatorDeriver) DeriveFlag(a any, s string, opts ...clap.Option) (clap.IFlag, error) {
	v, ok := a.(*parse.I)

	if !ok {
		return nil, fmt.Errorf("%w: want *parse.I but got %T", ErrIndicator, v)
	}

	return clap.NewFlag[parse.I](v, s, parse.Indicator{}, opts...), nil
}

func init() {
	clap.RegisterFlagDeriver("indicator", indicatorDeriver{})
}