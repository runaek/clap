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

// Derive attempts to construct a number of Arg dynamically from some struct-tags.
//
// The tags are applied on the fields to be bound to the Arg (i.e. at runtime, the Field will hold
// the value of the Arg). There are 4 formats for each of the types and all can be prefixed with a
// '!' which will mark the Arg as required:
//
// 1. KeyValueArg: `cli:"@<name>|<alias>:<deriver>"
//
// 2. FlagArg: `cli:"-<name>|<alias>:<deriver>"
//
// 3. PositionalArg: `cli:"#<index>:<deriver>"
//
// 4. PositionalArg(s): `cli:"#<index>...:<deriver>"
//
// Where the 'deriver' is the name assigned to the deriver when it was registered (using one of
// RegisterKeyValueDeriver, RegisterFlagDeriver or RegisterPositionalDeriver).
//
// Usage strings can be supplied by creating a string field in the struct called <Name>Usage, the
// same can be done for a default value (i.e. <Name>Default). At runtime, the values held by the
// variable will be used for the respective usage/default.
//
// Package `github.com/runaek/clap/pkg/derivers` provides support to a number of the
// built-in Go types. The program `github.com/runaek/clap/cmd/generate_derivers` is
// a codegen tool for helping generate FlagDeriver, KeyValueDeriver and PositionalDeriver
// implementations using a specified parse.Parse[T].
//
//
// 	type Args struct {
//		// a "name" or "N" KeyValueArg
//		Name string `cli:"!@name|N:string"`
//
//		// usage/description for the KeyValueArg
//		NameUsage string
//
//		// default value for the KeyValueArg
//		NameDefault string
//
//		// a variable number of string args, starting from the first index
//		Values []string `cli:"#1...:string"`
//
//		// usage/description for the positional args
//		ValuesUsage string
//	}
//
//	...
//
// 	var (
//		args = Args{
//			NameUsage:   "Supply your name.",
//			NameDefault: "John",
//			ValuesUsage: "A number of values to be supplied."
//		}
//
//		parser = clap.Must("demo", &args)
//	)
//
//	func main() {
//		parser.Parse()
//		parser.Ok()
//
//		// use args
//	}
// NOTE: Fields should be defined in the same order you would expect to add them to a Parser manually (i.e. Field
// referring to the 1st positional argument must come before the Field referring to the 2nd positional argument).
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

var (
	positionalDerivers = map[string]PositionalDeriver{}
	keyValueDerivers   = map[string]KeyValueDeriver{}
	flagDerivers       = map[string]FlagDeriver{}
)

// deriveArgs is a helper function for reading and creating Arg implementations using struct-tags and
// FlagDeriver, KeyValueDeriver and PositionalDeriver implementations that have been registered.
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
			opts = append(opts, WithShorthand(v.Alias))
		}

		if v.Default != "" {
			opts = append(opts, WithDefault(v.Default))
		}

		var arg Arg

		switch v.Type {
		case derive.KeyType:

			deriver := getKVD(v.Deriver)

			if deriver == nil {
				return nil, fmt.Errorf("unable to derive: %s", v.Deriver)
			}

			if a, err := deriver.DeriveKeyValue(v.Field(), v.Identifier, opts...); err != nil {
				return nil, fmt.Errorf("unable to derive key-value argument: %w", err)
			} else {
				arg = a
			}

		case derive.PosType:

			deriver := getPD(v.Deriver)

			if deriver == nil {
				return nil, fmt.Errorf("unable to derive: %s", v.Deriver)
			}

			if a, err := deriver.DerivePosition(v.Field(), v.Pos(), opts...); err != nil {
				return nil, fmt.Errorf("unable to derive positional argument: %w", err)
			} else {
				arg = a
			}
		case derive.FlagType:

			deriver := getFD(v.Deriver)

			if deriver == nil {
				return nil, fmt.Errorf("unable to derive: %s", v.Deriver)
			}
			if a, err := deriver.DeriveFlag(v.Field(), v.Identifier, opts...); err != nil {
				return nil, fmt.Errorf("unable to derive flag argument: %w", err)
			} else {
				arg = a
			}
		default:
			return nil, fmt.Errorf("invalid Type for derived Argument: %s", v.Type)
		}

		if arg == nil {
			return nil, errors.New("an error occurred deriving an Argument")
		}

		out = append(out, arg)
	}

	return out, nil
}

func getKVD(name string) KeyValueDeriver {
	return keyValueDerivers[name]
}

func getPD(name string) PositionalDeriver {
	return positionalDerivers[name]
}

func getFD(name string) FlagDeriver {
	return flagDerivers[name]
}
