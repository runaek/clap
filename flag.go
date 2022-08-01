package clap

import (
	"fmt"
	"github.com/runaek/clap/pkg/parse"
	"go.uber.org/zap"
	"strings"
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
	md := NewMetadata(opts...)

	if md.Usage() == "" {
		var zero T
		md.argUsage = fmt.Sprintf("%s - a %T flag variable.", name, zero)
	}
	f := &FlagArg[T]{
		Key: name,
		md:  md,
		v:   v,
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
	md := NewMetadata(opts...)

	if md.Usage() == "" {
		var zero T
		md.argUsage = fmt.Sprintf("%s - a repeatable %T flag variable.", name, zero)
	}

	f := &FlagArg[[]T]{
		Key: name,
		md:  md,
		v:   v,
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
	v   Variable[T]
	md  *Metadata

	repeatable, supplied, parsed bool
}

func (f *FlagArg[T]) argName() argName {
	return FlagType.getIdentifier(f.Name())
}

func (f *FlagArg[T]) Default() string {
	return f.md.Default()
}

func (f *FlagArg[T]) Name() string {
	return f.Key
}

func (f *FlagArg[T]) Type() Type {
	return FlagType
}

func (f *FlagArg[T]) Usage() string {
	return f.md.Usage()
}

func (f *FlagArg[T]) Shorthand() string {
	return f.md.Shorthand()
}

func (f *FlagArg[T]) ValueType() string {
	var zero T

	return strings.TrimPrefix(fmt.Sprintf("%T", zero), "*")
}

func (f *FlagArg[T]) IsRepeatable() bool {
	return f.repeatable
}

func (f *FlagArg[T]) IsRequired() bool {
	return f.md.IsRequired()
}

func (f *FlagArg[T]) IsParsed() bool {
	return f.parsed
}

func (f *FlagArg[T]) IsSupplied() bool {
	return f.supplied
}

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

func (f *FlagArg[T]) HasDefault() bool {
	return f.md.HasDefault()
}

func (f *FlagArg[T]) updateValue(s ...string) (err error) {
	v := f.Variable()

	log.Debug("Updating FlagArg argument value",
		zap.String("flag_name", f.Name()),
		zap.String("flag_type", f.ValueType()),
		zap.String("parser_type", fmt.Sprintf("%T", v.Parser())),
		zap.Strings("input", s),
		zap.Bool("parsed", f.parsed))

	if f.parsed {
		return nil
	}

	defer func() {
		if err == nil {
			f.parsed = true

			if f.IsIndicator() {
				if len(s) == 1 && s[0] != f.Default() {
					f.supplied = true
				}
			} else {
				if len(s) > 0 {
					f.supplied = true
				}
			}
		} else {
			log.Warn("Error updating FlagArg argument value",
				zap.String("flag_name", f.Name()),
				zap.String("flag_type", f.ValueType()),
				zap.Error(err))
		}
	}()

	var input []string

	if len(s) > 0 {
		input = s
	} else if f.HasDefault() {
		input = []string{
			f.Default(),
		}
	}

	return v.Update(input...)
}

func (f *FlagArg[T]) updateMetadata(opts ...Option) {
	if f.md == nil {
		f.md = NewMetadata(opts...)

		return
	}

	f.md.updateMetadata(opts...)
}

func (f *FlagArg[T]) isFlagArg() {}

var (
	_ Identifier    = Flag("")
	_ IFlag         = &FlagArg[any]{}
	_ TypedArg[any] = &FlagArg[any]{}
)
