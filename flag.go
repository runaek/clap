package clap

import (
	"fmt"
	"github.com/runaek/clap/pkg/parse"
)

// IFlag is the interface satisfied by a FlagArg.
type IFlag interface {
	Arg

	// IsIndicator returns true if the FlagArg does not require a value (e.g. a bool or C)
	IsIndicator() bool

	// HasDefault returns true if the flag has a defined default string value
	HasDefault() bool

	isFlagArg()
}

// Help is a constructor for some *FlagArg[bool] (identified by '--help'/'-h') type.
func Help(helpRequested *bool, desc string) *FlagArg[bool] {
	if desc == "" {
		desc = "Display the help-text for a command or program."
	}

	fl := NewFlagP[bool](helpRequested, "help", "h", parse.Bool{}, WithUsage(desc))
	return fl
}

// NewFlag is a constructor for a new *FlagArg[T].
func NewFlag[T any](val *T, name string, parser parse.Parser[T], opts ...Option) *FlagArg[T] {
	return FlagUsingVariable[T](name, NewVariable[T](val, parser), opts...)
}

// NewFlags is a constructor for a repeatable *FlagArg[[]T].
func NewFlags[T any](val *[]T, name string, parser parse.Parser[T], options ...Option) *FlagArg[[]T] {
	return FlagsUsingVariable[T](name, NewVariables[T](val, parser), options...)
}

// NewFlagsP is a constructor for a repeatable *FlagArg[[]T] with a shorthand.
func NewFlagsP[T any](val *[]T, name string, shorthand string, parser parse.Parser[T], options ...Option) *FlagArg[[]T] {
	f := NewFlagP[[]T](val, name, shorthand, parse.Slice[T](parser), options...)
	f.repeatable = true
	return f
}

// NewFlagP is a constructor for some new *FlagArg[T] with a shorthand.
func NewFlagP[T any](val *T, name string, shorthand string, parser parse.Parser[T], opts ...Option) *FlagArg[T] {
	opts = append(opts, WithAlias(shorthand))
	f := NewFlag[T](val, name, parser, opts...)
	return f
}

// FlagUsingVariable is a constructor for a FlagArg using a Id and some Variable.
func FlagUsingVariable[T any](name string, v Variable[T], opts ...Option) *FlagArg[T] {

	if md := NewMetadata(opts...); md.Usage() == "" {
		var zero T
		opts = append(opts, WithUsage(fmt.Sprintf("%s - a %T flag variable.", name, zero)))
	}

	f := &FlagArg[T]{
		Key:     name,
		argCore: newArgCoreUsing[T](v, opts...),
	}

	if f.Usage() == "" {
		var zero T
		f.argCore.md.argUsage = fmt.Sprintf("%s - a %T flag variable.", name, zero)
	}

	var zero T

	switch interface{}(zero).(type) {
	// C is a special type which is repeatable, its value will be the number of times
	// it was supplied
	case parse.C:
		f.repeatable = true
	}

	return f
}

// FlagsUsingVariable is a constructor for a repeatable FlagArg using a Id and some Variable.
func FlagsUsingVariable[T any](name string, v Variable[[]T], opts ...Option) *FlagArg[[]T] {

	if md := NewMetadata(opts...); md.Usage() == "" {
		var zero T
		opts = append(opts, WithUsage(fmt.Sprintf("%s - a repeatable %T flag variable.", name, zero)))
	}

	f := &FlagArg[[]T]{
		Key:     name,
		argCore: newArgCoreUsing[[]T](v, opts...),
	}

	f.repeatable = true

	return f
}

// Flag is an Identifier for some FlagArg.
type Flag string

func (f Flag) argName() argName {
	return FlagType.getIdentifier(string(f))
}

// A FlagArg argument of some type T.
type FlagArg[T any] struct {
	Key string
	*argCore[T]
}

func (f *FlagArg[T]) argName() argName {
	return FlagType.getIdentifier(f.Name())
}

func (f *FlagArg[T]) Name() string {
	return f.Key
}

func (f *FlagArg[T]) Type() Type {
	return FlagType
}

// func (f *FlagArg[T]) ValueType() string {
// 	var zero T
//
// 	return strings.TrimPrefix(fmt.Sprintf("%T", zero), "*")
// }

func (f *FlagArg[T]) Variable() Variable[T] {
	return f.v
}

func (f *FlagArg[T]) IsIndicator() bool {
	var zero T

	switch interface{}(zero).(type) {
	case *bool, bool, []bool, []*bool, parse.C:
		return true
	default:
		return false
	}
}

func (f *FlagArg[T]) updateValue(s ...string) error {

	if err := f.argCore.updateValue(s...); err == nil {
		f.argCore.parsed = true

		if f.IsIndicator() {
			if len(s) == 1 && s[0] == f.Default() {
				f.argCore.supplied = false
			}
		}
	} else {
		return err
	}

	return nil
}

func (f *FlagArg[T]) isFlagArg() {}

var (
	_ Identifier    = Flag("")
	_ IFlag         = &FlagArg[any]{}
	_ TypedArg[any] = &FlagArg[any]{}
)
