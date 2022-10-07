package clap

import (
	"fmt"
	"github.com/runaek/clap/pkg/parse"
)

// IKeyValue is the interface satisfied by a KeyValueArg.
type IKeyValue interface {
	Arg

	// HasDefault returns true if the Arg has a defined default string value
	HasDefault() bool

	isKeyValueArg()
}

// NewKeyValue is a constructor for a key-value argument from the command-line.
func NewKeyValue[T any](variable *T, name string, p parse.Parser[T], opts ...Option) *KeyValueArg[T] {
	v := NewVariable[T](variable, p)

	return KeyValueUsingVariable[T](name, v, opts...)
}

// NewKeyValues is a constructor for a repeatable key-valued arguments from the
// command-line.
//
// Automatically converts the Func[T] into a Func[[]T] via parse.Slice - use
// KeyValuesUsingVariable to be able to change this behaviour as required.
func NewKeyValues[T any](variables *[]T, name string, p parse.Parser[T], opts ...Option) *KeyValueArg[[]T] {
	v := NewVariables[T](variables, p)

	return KeyValuesUsingVariable(name, v, opts...)
}

// KeyValueUsingVariable allows a KeyValueArg to be constructed using a Variable.
func KeyValueUsingVariable[T any](name string, v Variable[T], opts ...Option) *KeyValueArg[T] {

	if md := NewMetadata(opts...); md.Usage() == "" {
		var zero T
		opts = append(opts, WithUsage(fmt.Sprintf("%s - a %T key-value variable.", name, zero)))
	}

	kv := &KeyValueArg[T]{
		Key:     name,
		argCore: newArgCoreUsing[T](v, opts...),
	}

	return kv
}

// KeyValuesUsingVariable allows a repeatable KeyValueArg to be constructed
// using a Variable.
func KeyValuesUsingVariable[T any](name string, v Variable[[]T], opts ...Option) *KeyValueArg[[]T] {

	if md := NewMetadata(opts...); md.Usage() == "" {
		var zero T
		opts = append(opts, WithUsage(fmt.Sprintf("%s - a repeatable %T key-value variable.", name, zero)))
	}

	core := newArgCoreUsing[[]T](v, opts...)
	core.repeatable = true

	kv := &KeyValueArg[[]T]{
		Key:     name,
		argCore: core,
	}

	return kv
}

// A Key is an Identifier for some KeyValueArg.
type Key string

func (k Key) argName() argName {
	return KeyValueType.getIdentifier(string(k))
}

// A KeyValueArg represents a key=value argument where the key is a string and
// the value is a string representation for some type T.
//
// Should be created by the functions: NewKeyValue, NewKeyValues,
// KeyValueUsingVariable and KeyValuesUsingVariable.
type KeyValueArg[T any] struct {
	Key string
	*argCore[T]
}

func (k *KeyValueArg[T]) argName() argName {
	return KeyValueType.getIdentifier(k.Name())
}

func (k *KeyValueArg[T]) Name() string {
	return k.Key
}

func (k *KeyValueArg[T]) Type() Type {
	return KeyValueType
}

func (_ *KeyValueArg[T]) isKeyValueArg() {}

var (
	_ Identifier    = Key("")
	_ TypedArg[any] = &KeyValueArg[any]{}
	_ IKeyValue     = &KeyValueArg[any]{}
)
