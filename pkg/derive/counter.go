// Code generated by cmd/generate_drivers. DO NOT EDIT

package derive 

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

const (
	ErrCounter clap.Error = "unable to derive 'counter' Argument"
)

type counterDeriver struct {}

func (_ counterDeriver) DeriveKeyValue(a any, s string, opts ...clap.Option) (clap.IKeyValue, error) {
	
	v, ok := a.(*parse.C)
	
	if !ok {
		return nil, fmt.Errorf("%w: want *int but got %T", ErrCounter, v)
	}
	
	return clap.NewKeyValue[parse.C](v, s, parse.Counter{}, opts...), nil 
}

func (_ counterDeriver) DerivePosition(a any, s int, opts ...clap.Option) (clap.IPositional, error) {
	
	v, ok := a.(*parse.C)
	
	if !ok {
		return nil, fmt.Errorf("%w: want *int but got %T", ErrCounter, v)
	}
	
	return clap.NewPosition[parse.C](v, s, parse.Counter{}, opts...), nil 
}

func (_ counterDeriver) DeriveFlag(a any, s string, opts ...clap.Option) (clap.IFlag, error) {
	
	v, ok := a.(*parse.C)
	
	if !ok {
		return nil, fmt.Errorf("%w: want *int but got %T", ErrCounter, v)
	}
	
	return clap.NewFlag[parse.C](v, s, parse.Counter{}, opts...), nil 
}

func init() {
	clap.RegisterFlagDeriver("counter", counterDeriver{})
	clap.RegisterPositionalDeriver("counter", counterDeriver{})
	clap.RegisterKeyValueDeriver("counter", counterDeriver{})
}
