package clap

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

// IsPositional represents an Arg of Type: PositionType.
//
// See PositionalArg.
type IsPositional interface {
	Arg

	// Index returns the index (position) of the PositionalArg (or the starting index of the remaining args if
	// ArgRemaining is true)
	Index() int

	mustEmbedPosition()
}

func NewPosition[T any](variable *T, index int, parser ValueParser[T], opts ...Option) *PositionalArg[T] {
	opts = append(opts, positionOptions...)
	return &PositionalArg[T]{
		md:       NewMetadata(opts...),
		v:        NewVariable[T](variable, parser),
		argIndex: index,
	}
}

func NewPositions[T any](variables *[]T, fromIndex int, parser ValueParser[T], opts ...Option) *PositionalArg[[]T] {
	opts = append(opts, positionOptions...)
	return &PositionalArg[[]T]{
		md:            NewMetadata(opts...),
		v:             NewVariable[[]T](variables, SliceParser[T](parser)),
		argStartIndex: fromIndex,
	}
}

// Position is an Identifier for some PositionalArg.
type Position int

func (i Position) identify() argName {
	return PositionType.getIdentifier(fmt.Sprintf("%d", i))
}

var positionOptions = []Option{
	withNoShorthand(), withDefaultDisabled(),
}

// A PositionalArg represents a particular index (or indexes) of positional arguments representing some type T.
//
// Should be created by the NewPosition and NewPositions functions.
type PositionalArg[T any] struct {
	md *Metadata
	v  Variable[T]

	argIndex      int // > 0 indicates a single position
	argStartIndex int // > 0 indicates a variableParser number of positions (the final N args)

	parsed, supplied bool
}

func (p PositionalArg[T]) identify() argName {
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

func (p *PositionalArg[T]) Usage() string {
	return p.md.Usage()
}

func (p *PositionalArg[T]) Shorthand() rune {
	return 0
}

func (p *PositionalArg[T]) ValueType() string {
	var zero T

	return strings.TrimPrefix(fmt.Sprintf("%T", zero), "*")
}

func (p *PositionalArg[T]) IsRepeatable() bool {
	return p.argStartIndex > 0
}

func (p *PositionalArg[T]) IsRequired() bool {
	return p.md.IsRequired()
}

func (p *PositionalArg[T]) IsParsed() bool {
	return p.parsed
}

func (p *PositionalArg[T]) IsSupplied() bool {
	return p.supplied
}

func (p *PositionalArg[T]) Variable() Variable[T] {
	return p.v
}

func (p *PositionalArg[T]) Index() int {
	if p.argIndex > 0 {
		return p.argIndex
	}
	return p.argStartIndex
}

func (p *PositionalArg[T]) updateMetadata(opts ...Option) {
	opts = append(opts, positionOptions...)

	if p.md == nil {
		p.md = NewMetadata(opts...)
		return
	}

	p.md.updateMetadata(opts...)
}

func (p *PositionalArg[T]) updateValue(inputs ...string) (err error) {

	v := p.Variable()

	log.Debug("Updating PositionalArg",
		zap.String("pos_index", p.Name()),
		zap.String("pos_type", p.ValueType()),
		zap.Strings("input", inputs),
		zap.String("parser_type", fmt.Sprintf("%T", v.parser())),
		zap.Bool("parsed", p.parsed))

	if p.parsed {
		return nil
	}

	defer func() {

		if err == nil {
			p.parsed = true

			if len(inputs) > 0 {
				p.supplied = true
			}
		} else {
			log.Debug("Error updating Positional argument value",
				zap.String("pos_index", p.Name()),
				zap.String("pos_type", p.ValueType()),
				zap.Error(err))
		}
	}()

	if p.IsRepeatable() {
		return v.Update(inputs...)
	} else {
		return v.Update(inputs[0])
	}
}

func (_ *PositionalArg[T]) mustEmbedPosition() {}

var (
	_ Identifier    = Position(0)
	_ TypedArg[any] = &PositionalArg[any]{}
	_ IsPositional  = &PositionalArg[any]{}
)
