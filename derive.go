package clap

import (
	"errors"
	"fmt"
	"github.com/runaek/clap/internal/derive"
)

// A PositionalDeriver is responsible for constructing a PositionalArg dynamically.
type PositionalDeriver interface {
	DerivePosition(any, int, ...Option) (IPositional, error)
}

// A KeyValueDeriver is responsible for constructing a KeyValueArg dynamically.
type KeyValueDeriver interface {
	DeriveKeyValue(any, string, ...Option) (IKeyValue, error)
}

// A FlagDeriver is responsible for constructing a FlagArg dynamically.
type FlagDeriver interface {
	DeriveFlag(any, string, ...Option) (IFlag, error)
}

func RegisterPositionalDeriver(name string, d PositionalDeriver) {
	positionalDerivers[name] = d
}

func RegisterKeyValueDeriver(name string, d KeyValueDeriver) {
	keyValueDerivers[name] = d
}

func RegisterFlagDeriver(name string, d FlagDeriver) {
	flagDerivers[name] = d
}

var (
	positionalDerivers = map[string]PositionalDeriver{}
	keyValueDerivers   = map[string]KeyValueDeriver{}
	flagDerivers       = map[string]FlagDeriver{}
)

func deriveArgs(src any) ([]Arg, error) {

	program, err := derive.Parse(src)

	if err != nil {
		return nil, err
	}

	var out []Arg

	for _, v := range program.Args() {

		opts := []Option{
			WithUsage(v.Usage),
		}

		if v.Required {
			opts = append(opts, AsRequired())
		}

		if v.Alias != "" {
			opts = append(opts, WithShorthand(v.Alias[:1]))
		}

		if v.Default != "" {
			opts = append(opts, WithDefault(v.Default))
		}

		var arg Arg

		switch v.Type {
		case derive.KeyType:

			deriver := lookupKVDeriver(v.Deriver)

			if deriver == nil {
				return nil, fmt.Errorf("unable to derive: %s", v.Deriver)
			}

			if a, err := deriver.DeriveKeyValue(v.Field(), v.Identifier, opts...); err != nil {
				return nil, fmt.Errorf("unable to derive key-value argument: %w", err)
			} else {
				arg = a
			}

		case derive.PosType:

			deriver := lookupPDeriver(v.Deriver)

			if deriver == nil {
				return nil, fmt.Errorf("unable to derive: %s", v.Deriver)
			}

			if a, err := deriver.DerivePosition(v.Field(), v.Pos(), opts...); err != nil {
				return nil, fmt.Errorf("unable to derive positional argument: %w", err)
			} else {
				arg = a
			}
		case derive.FlagType:

			deriver := lookupFDeriver(v.Deriver)

			if deriver == nil {
				return nil, fmt.Errorf("unable to derive: %s", v.Deriver)
			}
			if a, err := deriver.DeriveFlag(v.Field(), v.Identifier, opts...); err != nil {
				return nil, fmt.Errorf("unable to derive flag argument: %w", err)
			} else {
				arg = a
			}
		default:
			fmt.Printf("uh oh")
			return nil, fmt.Errorf("invalid Type for derived Argument: %s", v.Type)
		}

		if arg == nil {
			return nil, errors.New("an error occurred deriving an Argument")
		}

		out = append(out, arg)
	}

	return out, nil
}

func lookupKVDeriver(name string) KeyValueDeriver {
	return keyValueDerivers[name]
}

func lookupPDeriver(name string) PositionalDeriver {
	return positionalDerivers[name]
}

func lookupFDeriver(name string) FlagDeriver {
	return flagDerivers[name]
}

// Derive attempts to construct a number of Arg dynamically from some struct-tags.
func Derive(source any) ([]Arg, error) {
	return deriveArgs(source)
}

// DeriveAll is a convenient wrapper around Derive.
func DeriveAll(sources ...any) ([]Arg, error) {
	out := make([]Arg, 0, len(sources))

	for _, s := range sources {

		if a, ok := s.(Arg); ok {
			out = append(out, a)
			continue
		}

		more, err := deriveArgs(s)

		if err != nil {
			return out, err
		}

		out = append(out, more...)

	}

	return out, nil
}
