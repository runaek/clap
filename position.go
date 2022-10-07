package clap

import (
	"fmt"
	"github.com/runaek/clap/pkg/parse"
)

// IPositional is the interface satisfied by a PositionalArg.
type IPositional interface {
	Arg

	// Index returns the index (position) of the PositionalArg (or the starting index of the remaining args if
	// ArgRemaining is true)
	Index() int

	isPositionalArg()
}

// NewPosition is a constructor for a positional argument at some index.
func NewPosition[T any](variable *T, index int, parser parse.Parser[T], opts ...Option) *PositionalArg[T] {
	return PositionUsingVariable[T](index, NewVariable[T](variable, parser), opts...)
}

// NewPositions is a constructor for a number of positional arguments starting from some index.
func NewPositions[T any](variables *[]T, fromIndex int, parser parse.Parser[T], opts ...Option) *PositionalArg[[]T] {
	return PositionsUsingVariable[T](fromIndex, NewVariable[[]T](variables, parse.Slice[T](parser)), opts...)
}

func PositionUsingVariable[T any](index int, v Variable[T], opts ...Option) *PositionalArg[T] {
	opts = append(opts, positionOptions...)

	if md := NewMetadata(opts...); md.Usage() == "" {
		opts = append(opts, WithUsage("A positional argument."))
	}

	core := newArgCoreUsing[T](v, opts...)

	return &PositionalArg[T]{
		argCore:  core,
		argIndex: index,
	}
}

func PositionsUsingVariable[T any](fromIndex int, v Variable[[]T], opts ...Option) *PositionalArg[[]T] {
	opts = append(opts, positionOptions...)

	if md := NewMetadata(opts...); md.Usage() == "" {
		opts = append(opts, WithUsage("Remaining positional arguments."))
	}

	core := newArgCoreUsing[[]T](v, opts...)
	core.repeatable = true

	return &PositionalArg[[]T]{
		argCore:       core,
		argIndex:      0,
		argStartIndex: fromIndex,
	}
}

// Position is an Identifier for some PositionalArg.
type Position int

func (i Position) argName() argName {
	return PositionType.getIdentifier(fmt.Sprintf("%d", i))
}

var positionOptions = []Option{
	withDefaultDisabled(),
}

// A PositionalArg represents a particular index (or indexes) of positional arguments representing some type T.
//
// Should be created by the NewPosition and NewPositions functions.
type PositionalArg[T any] struct {
	*argCore[T]
	argIndex      int // > 0 indicates a single position
	argStartIndex int // > 0 indicates a variable number of positions (the final N args)
}

func (p PositionalArg[T]) argName() argName {
	return PositionType.getIdentifier(p.Name())
}

func (p *PositionalArg[T]) Default() string {
	return ""
}

func (p *PositionalArg[T]) Name() string {
	if p.IsRepeatable() {
		return fmt.Sprintf("%d", p.argStartIndex)
	}
	return fmt.Sprintf("%d", p.argIndex)
}

func (p *PositionalArg[T]) Type() Type {
	return PositionType
}

func (p *PositionalArg[T]) IsRepeatable() bool {
	return p.argStartIndex > 0
}

func (p *PositionalArg[T]) IsRequired() bool {
	if p.argStartIndex > 0 {
		return false
	}
	return p.md.IsRequired()
}

func (p *PositionalArg[T]) Index() int {
	if p.argIndex > 0 {
		return p.argIndex
	}
	return p.argStartIndex
}

func (p *PositionalArg[T]) updateValue(inputs ...string) (err error) {
	v := p.Variable()

	if p.parsed {
		return nil
	}
	defer func() {

		if len(inputs) > 0 {
			p.supplied = true
		}

		if err == nil {
			p.parsed = true
		}
	}()
	if p.IsRepeatable() {
		return v.Update(inputs...)
	} else {
		return v.Update(inputs[0])
	}
}

func (_ *PositionalArg[T]) isPositionalArg() {}

var (
	_ Identifier         = Position(0)
	_ TypedArg[any]      = &PositionalArg[any]{}
	_ TypedArg[[]string] = &PositionalArg[[]string]{}
	_ IPositional        = &PositionalArg[any]{}
)
