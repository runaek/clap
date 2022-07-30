package clap

import (
	"fmt"
	"github.com/runaek/clap/pkg/parse"
	"go.uber.org/zap"
	"strings"
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

// NewKeyValues is a constructor for a repeatable key-valued arguments from the command-line.
//
// Automatically converts the Func[T] into a Func[[]T] via parse.Slice - use KeyValuesUsingVariable
// to be able to change this behaviour as required.
func NewKeyValues[T any](variables *[]T, name string, p parse.Parser[T], opts ...Option) *KeyValueArg[[]T] {
	v := NewVariables[T](variables, p)

	return KeyValuesUsingVariable(name, v, opts...)
}

// KeyValueUsingVariable allows a KeyValueArg to be constructed using a Variable.
func KeyValueUsingVariable[T any](name string, v Variable[T], opts ...Option) *KeyValueArg[T] {
	md := NewMetadata(opts...)

	if md.Usage() == "" {
		var zero T
		md.argUsage = fmt.Sprintf("%s - a %T key-value variable.", name, zero)
	}

	kv := &KeyValueArg[T]{
		Key: name,
		md:  md,
		v:   v,
	}

	return kv
}

// KeyValuesUsingVariable allows a repeatable KeyValueArg to be constructed using Variable.
func KeyValuesUsingVariable[T any](name string, v Variable[[]T], opts ...Option) *KeyValueArg[[]T] {
	md := NewMetadata(opts...)

	if md.Usage() == "" {
		var zero T
		md.argUsage = fmt.Sprintf("%s - a %T key-value variable.", name, zero)
	}
	kv := &KeyValueArg[[]T]{
		Key:        name,
		md:         md,
		v:          v,
		repeatable: true,
	}

	return kv
}

// A Key is an Identifier for some KeyValueArg.
type Key string

func (k Key) argName() argName {
	return KeyValueType.getIdentifier(string(k))
}

// A KeyValueArg represents a key=value argument where the key is a string and the value is a string representation for
// some type T.
//
// Should be created by the functions: NewKeyValue, NewKeyValues, KeyValueUsingVariable and KeyValuesUsingVariable.
type KeyValueArg[T any] struct {
	Key string
	md  *Metadata
	v   Variable[T]

	repeatable, supplied, parsed bool
}

func (k *KeyValueArg[T]) argName() argName {
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

func (k *KeyValueArg[T]) Shorthand() string {
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
		zap.String("parser_type", fmt.Sprintf("%T", v.Parser())),
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

func (_ *KeyValueArg[T]) isKeyValueArg() {}

var (
	_ Identifier    = Key("")
	_ TypedArg[any] = &KeyValueArg[any]{}
	_ IKeyValue     = &KeyValueArg[any]{}
)
