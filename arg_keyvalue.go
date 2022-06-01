package clap

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

// IsKeyValue represents an Arg of Type: KeyValueType.
//
// See KeyValueArg.
type IsKeyValue interface {
	Arg

	// HasDefault returns true if the Arg has a defined default string value
	HasDefault() bool

	mustEmbedKey()
}

// NewKeyValue is a constructor for a key-value argument from the command-line.
func NewKeyValue[T any](variable *T, name string, p ValueParser[T], opts ...Option) *KeyValueArg[T] {
	return &KeyValueArg[T]{
		Key: name,
		md:  NewMetadata(opts...),
		v:   NewVariable[T](variable, p),
	}
}

// NewKeyValues is a constructor for a repeatable key-valued arguments from the command-line.
func NewKeyValues[T any](variables *[]T, name string, p ValueParser[T], opts ...Option) *KeyValueArg[[]T] {
	return &KeyValueArg[[]T]{
		Key:        name,
		md:         NewMetadata(opts...),
		v:          NewVariables[T](variables, p),
		repeatable: true,
	}
}

// A Key is an Identifier for some KeyValueArg.
type Key string

func (k Key) identify() argName {
	return KeyValueType.getIdentifier(string(k))
}

// A KeyValueArg represents a key=value argument where the key is a string and the value is a string representation for
// some type T.
//
// Should be created by the NewKeyValue and NewKeyValues functions.
type KeyValueArg[T any] struct {
	Key string
	md  *Metadata
	v   Variable[T]

	repeatable, supplied, parsed bool
}

func (k *KeyValueArg[T]) identify() argName {
	return KeyValueType.getIdentifier(k.Name())
}

func (k *KeyValueArg[T]) Default() string {
	return k.md.Default()
}

func (k *KeyValueArg[T]) Name() string {
	return k.Key
}

func (k *KeyValueArg[T]) Type() Type {
	return KeyValueType
}

func (k *KeyValueArg[T]) Usage() string {
	return k.md.Usage()
}

func (k *KeyValueArg[T]) Shorthand() rune {
	return k.md.Shorthand()
}

func (k *KeyValueArg[T]) ValueType() string {
	var zero T
	return strings.TrimPrefix(fmt.Sprintf("%T", zero), "*")
}

func (k *KeyValueArg[T]) IsRepeatable() bool {
	return k.repeatable
}

func (k *KeyValueArg[T]) IsRequired() bool {
	return k.md.IsRequired()
}

func (k *KeyValueArg[T]) IsParsed() bool {
	return k.parsed
}

func (k *KeyValueArg[T]) IsSupplied() bool {
	return k.supplied
}

func (k *KeyValueArg[T]) Variable() Variable[T] {
	return k.v
}

func (k *KeyValueArg[T]) HasDefault() bool {
	return k.md.HasDefault()
}

func (k *KeyValueArg[T]) updateMetadata(opts ...Option) {
	if k.md == nil {
		k.md = NewMetadata(opts...)
		return
	}
	k.md.updateMetadata(opts...)
}

func (k *KeyValueArg[T]) updateValue(s ...string) (err error) {

	v := k.Variable()

	log.Debug("Updating Key Value argument value",
		zap.String("kv_name", k.Name()),
		zap.String("kv_type", k.ValueType()),
		zap.String("parser_type", fmt.Sprintf("%T", v.parser())),
		zap.Strings("input", s),
		zap.Bool("parsed", k.parsed))

	if k.parsed {
		return nil
	}

	defer func() {
		if err == nil {
			k.parsed = true
			if len(s) > 0 && s[0] != "" {
				k.supplied = true
			}
		} else {
			log.Warn("Error updating Key Value argument value",
				zap.String("kv_name", k.Name()),
				zap.String("kv_type", k.ValueType()),
				zap.Error(err))
		}
	}()

	var input []string

	if len(s) > 0 {
		input = s
	} else if k.HasDefault() {
		input = []string{
			k.Default(),
		}
	}

	return v.Update(input...)
}

func (_ *KeyValueArg[T]) mustEmbedKey() {}

var (
	_ Identifier    = Key("")
	_ TypedArg[any] = &KeyValueArg[any]{}
	_ IsKeyValue    = &KeyValueArg[any]{}
)
