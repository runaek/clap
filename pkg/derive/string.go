package derive

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

func init() {
	// string
	clap.RegisterFlagDeriver("string", stringDeriver{})
	clap.RegisterPositionalDeriver("string", stringDeriver{})
	clap.RegisterKeyValueDeriver("string", stringDeriver{})

	// strings
	clap.RegisterFlagDeriver("strings", stringsDeriver{})
	clap.RegisterPositionalDeriver("strings", stringsDeriver{})
	clap.RegisterKeyValueDeriver("strings", stringsDeriver{})
}

const (
	ErrString clap.Error = "unable to derive 'string' Argument"
)

type stringDeriver struct{}

func (s stringDeriver) DerivePosition(a any, i int, opts ...clap.Option) (clap.IPositional, error) { // nolint: ireturn

	v, ok := a.(*string)

	if !ok {
		return nil, fmt.Errorf("%w: expected *string but got %T", ErrString, a)
	}

	return clap.NewPosition[string](v, i, parse.String{}, opts...), nil
}

func (s stringDeriver) DeriveKeyValue(a any, name string, opts ...clap.Option) (clap.IKeyValue, error) { // nolint: ireturn
	v, ok := a.(*string)

	if !ok {
		return nil, fmt.Errorf("%w: expected *string but got %T", ErrString, a)
	}

	return clap.NewKeyValue[string](v, name, parse.String{}, opts...), nil
}

func (s stringDeriver) DeriveFlag(a any, name string, opts ...clap.Option) (clap.IFlag, error) { // nolint: ireturn
	v, ok := a.(*string)

	if !ok {
		return nil, fmt.Errorf("%w: expected *string but got %T", ErrString, a)
	}

	return clap.NewFlag[string](v, name, parse.String{}, opts...), nil
}

type stringsDeriver struct{}

func (s stringsDeriver) DerivePosition(a any, i int, opts ...clap.Option) (clap.IPositional, error) { // nolint: ireturn

	v, ok := a.(*[]string)

	if !ok {
		return nil, fmt.Errorf("%w: expected *string but got %T", ErrString, a)
	}

	return clap.NewPositions[string](v, i, parse.String{}, opts...), nil
}

func (s stringsDeriver) DeriveKeyValue(a any, name string, opts ...clap.Option) (clap.IKeyValue, error) { // nolint: ireturn
	v, ok := a.(*[]string)

	if !ok {
		return nil, fmt.Errorf("%w: expected *string but got %T", ErrString, a)
	}

	return clap.NewKeyValues[string](v, name, parse.String{}, opts...), nil
}

func (s stringsDeriver) DeriveFlag(a any, name string, opts ...clap.Option) (clap.IFlag, error) { // nolint: ireturn
	v, ok := a.(*[]string)

	if !ok {
		return nil, fmt.Errorf("%w: expected *string but got %T", ErrString, a)
	}

	return clap.NewFlags[string](v, name, parse.String{}, opts...), nil
}
