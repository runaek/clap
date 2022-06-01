package clap

import (
	"fmt"
	"github.com/runaek/clap/pkg/parsers"
	"go.uber.org/zap"
	"strings"
)

// IsFlag represents an Arg of Type: FlagType.
//
// See FlagArg.
type IsFlag interface {
	Arg

	// IsIndicator returns true if the FlagArg does not require a value (e.g. a bool or Counter)
	IsIndicator() bool

	// HasDefault returns true if the flag has a defined default string value
	HasDefault() bool

	isFlag()
}

// Help is a constructor for some *FlagArg[bool] (identified by '--help'/'-h') type.
func Help(helpRequested *bool, desc string) *FlagArg[bool] {

	if desc == "" {
		desc = "Display the help-text for a command or program."
	}

	fl := NewFlagP[bool](helpRequested, "help", 'h', parsers.Bool, WithUsage(desc), WithDefault("false"))
	return fl
}

// NewFlag is a constructor for a new *FlagArg[T].
func NewFlag[T any](val *T, name string, parser ValueParser[T], opts ...Option) *FlagArg[T] {
	return &FlagArg[T]{
		Key: name,
		md:  NewMetadata(opts...),
		v:   NewVariable(val, parser),
	}
}

// NewFlags is a constructor for a repeatable *FlagArg[[]T].
func NewFlags[T any](val *[]T, name string, parser ValueParser[T], options ...Option) *FlagArg[[]T] {
	f := NewFlag[[]T](val, name, SliceParser(parser), options...)
	f.repeatable = true
	return f
}

// NewFlagsP is a constructor for a repeatable *FlagArg[[]T] with a shorthand.
func NewFlagsP[T any](val *[]T, name string, shorthand rune, parser ValueParser[T], options ...Option) *FlagArg[[]T] {
	f := NewFlagP[[]T](val, name, shorthand, SliceParser(parser), options...)
	f.repeatable = true
	return f
}

// NewFlagP is a constructor for some new *FlagArg[T] with a shorthand.
func NewFlagP[T any](val *T, name string, shorthand rune, parser ValueParser[T], opts ...Option) *FlagArg[T] {
	opts = append(opts, WithShorthand(shorthand))
	f := NewFlag[T](val, name, parser, opts...)

	var zero T

	switch interface{}(zero).(type) {
	// Counter is a special type which is repeatable, its value will be the number of times
	// it was supplied
	case Counter:
		f.repeatable = true
	}

	return f
}

// Flag is an Identifier for some FlagArg.
type Flag string

func (f Flag) identify() argName {
	return FlagType.getIdentifier(string(f))
}

// A FlagArg argument of some type T.
type FlagArg[T any] struct {
	Key string
	v   Variable[T]
	md  *Metadata

	repeatable, supplied, parsed bool
}

func (f *FlagArg[T]) identify() argName {
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

func (f *FlagArg[T]) Shorthand() rune {
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
	case *bool, bool, []bool, []*bool, Counter:
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
		zap.String("parser_type", fmt.Sprintf("%T", v.parser())),
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

func (f *FlagArg[T]) isFlag() {}

var (
	_ Identifier    = Flag("")
	_ IsFlag        = &FlagArg[any]{}
	_ TypedArg[any] = &FlagArg[any]{}
)
