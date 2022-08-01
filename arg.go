package clap

import (
	"errors"
	"fmt"
	"github.com/runaek/clap/pkg/parse"
	"strings"
)

//go:generate mockgen -destination=mocks.go -copyright_file=internal/mock_header.txt -package=clap github.com/runaek/clap FileReader,FileWriter,FileInfo,PositionalDeriver,KeyValueDeriver,FlagDeriver

// Arg is the shared behaviour of all command-line input types (FlagType, KeyValueType and PositionType). It essentially
// exposes an API similar to the behaviour seen in the standard 'flags' package, extending support to key-value and
// positional arguments.
//
// Metadata is attached to each Arg giving each arg some mutable properties which can be managed via Option(s).
//
// Each Arg has its own go-argumentVariable which it is responsible for updating:
//
//  import (
// 		. "github.com/runaek/clap"
//		"github.com/runaek/clap/pkg/parse"
//	)
//
//	var (
//		strVal string
//		numVal int
//		myArg = NewKeyValue[string](&strVal, "arg1", parse.String, <options>...) // arg1=Test123 => strVal=Test123
//		numFlag = NewFlag[int](&numVal, "amount", parse.Int, <options>...)       // --amount=123 => numVal=123
//
//		parser = clap.New("program-Id").
//				Add(myArg, numFlag).
//				OK()
//	)
//
//	...
//
//	func main() {
//		parser.Parse()
//		parser.OK()
//
//		... // do stuff with 'strVal' and 'numVal'
//	}
type Arg interface {
	Identifier

	Updater

	// Name returns the key/identifier of the Arg, this should be unique for each Type.
	Name() string

	// Type indicates the type of command-line argument the Arg is, returns one of:
	//
	// 	> FlagType;
	//
	//	> KeyValueType;
	//
	//	> PositionType;
	//
	//	> PipeType;
	Type() Type

	// Shorthand is the single character alias/identifier for the Arg, if applicable
	//
	// Can be updated via the WithAlias Option
	Shorthand() string

	// Usage returns a usage of the Arg for the user
	//
	// Can be updated via the WithUsage Option
	Usage() string

	// Default returns the default string value for the Arg
	//
	// Can be updated via the WithDefault Option, if applicable
	Default() string

	// IsRequired returns true if the Arg is required
	//
	// Can be updated via the AsRequired, AsOptional Option(s)
	IsRequired() bool

	// IsParsed returns true if the Arg has been parsed
	IsParsed() bool

	// IsSupplied returns true if the Arg was supplied by the user
	IsSupplied() bool

	// IsRepeatable returns true if the Arg can be supplied multiple times
	IsRepeatable() bool
}

// Generalize is a helper function for converting some TypedArg[T] -> Arg.
func Generalize[T any](typed ...TypedArg[T]) []Arg {
	out := make([]Arg, len(typed))

	for i, t := range typed {
		out[i] = t.(Arg) //nolint:forcetypeassert
	}
	return out
}

// TypedArg is the generic interface that is satisfied by all Arg implementations.
//
// This is a convenience interface that provides access to the underlying generic Variable.
type TypedArg[T any] interface {
	Arg
	// Variable returns the underlying Variable for the Arg
	Variable() Variable[T]
}

// An Identifier represents something we can use to identify an Arg.
type Identifier interface {
	// argName returns an identifier for an Arg
	argName() argName
}

// NameOf is a helper function for getting the name of an Arg from an Identifier.
func NameOf(id Identifier) string {
	return id.argName().Name()
}

// TypeOf is a helper function for getting the type of Arg from an Identifier.
func TypeOf(id Identifier) Type {
	return id.argName().Type()
}

const (
	ErrInvalidType Error = "invalid type"
)

// ValueOf is a helper function for retrieving the value of an Arg.
func ValueOf[T any](arg Arg) (T, error) {
	var zero T

	if typed, ok := arg.(TypedArg[T]); ok {
		return typed.Variable().Unwrap(), nil
	}

	return zero, fmt.Errorf("%w: unable to read %T value from %T", ErrInvalidType, zero, arg)
}

// ReferenceTo is a helper function for retrieving a reference to the value for an Arg.
func ReferenceTo[T any](arg Arg) (*T, error) {
	if arg == nil {
		return nil, errors.New("nil reference")
	}

	if typ, ok := arg.(TypedArg[T]); ok {
		return typ.Variable().Ref(), nil
	}

	var zero T

	return nil, fmt.Errorf("%w: unable to reference variable of type %T from %T", ErrInvalidType, zero, arg)
}

// Updater is the shared private behaviour shared by all Arg which allows mutable *Metadata fields to be updated
// and the underlying value/variable of an Arg to be updated.
type Updater interface {
	updateValue(...string) error
	updateMetadata(...Option)
}

// UpdateValue is a public helper function to update the value of some Arg.
func UpdateValue(u Updater, input ...string) error {
	if len(input) == 0 {
		return u.updateValue()
	}
	return u.updateValue(input...)
}

// UpdateMetadata is the public helper function to update the Metadata of some Arg.
func UpdateMetadata(u Updater, options ...Option) {
	u.updateMetadata(options...)
}

// argName is the private key used to index/Id an Arg.
//
// This is a string formatted: <Type>:<Id>
type argName string

func (id argName) Type() Type {
	parts := strings.SplitN(string(id), ":", 2)

	switch strings.ToLower(parts[0]) {
	case KeyValueType.String():
		return KeyValueType
	case FlagType.String():
		return FlagType
	case PositionType.String():
		return PositionType
	case PipeType.String():
		return PipeType
	}

	return Unrecognised
}

func (id argName) String() string {
	return string(id)
}

func (id argName) Name() string {
	return strings.SplitN(string(id), ":", 2)[1]
}

func (id argName) argName() argName {
	return id
}

// ValidateDefaultValue is a helper function for retrieving and attempting to parse the actual typed default value
// of an Arg.
func ValidateDefaultValue[T any](arg TypedArg[T]) (T, error) {
	var zero T
	switch a := arg.(type) {
	case IFlag:
		switch a.(type) {
		case *FlagArg[bool], *FlagArg[parse.C]:
			if !a.HasDefault() {
				return zero, nil
			}
		}
	}

	v := arg.Variable()

	parser := v.Parser()

	defaultValue := arg.Default()

	val, err := parser.Parse(defaultValue)

	if err != nil {
		return zero, fmt.Errorf("invalid default value %q: %w", defaultValue, err)
	}

	return val, nil
}

// Type indicates the 'type' of input being processed (either a flag, key-value argument or a positional argument).
type Type int

func (t Type) String() string {
	switch t {
	case FlagType:
		return "flag"
	case KeyValueType:
		return "key-value"
	case PositionType:
		return "position"
	case PipeType:
		return "pipe"
	}
	return "unrecognised"
}

func (t Type) getIdentifier(name string) argName {
	return argName(t.String() + ":" + name)
}

const (
	Unrecognised Type = iota
	FlagType
	KeyValueType
	PositionType
	PipeType
	limit
)
