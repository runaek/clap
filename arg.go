package clap

import (
	"fmt"
	"strings"
)

// UpdateArgValue is a helper function to update the value of some Arg.
//
// If no input is supplied, then the default is attempted to be used
func UpdateArgValue(a Arg, input ...string) error {
	return a.updateValue(input...)
}

// UpdateArgMetadata is a helper function to update some mutable metadata of an Arg.
func UpdateArgMetadata(a Arg, options ...Option) {
	a.updateMetadata(options...)
}

type Identifier interface {
	identify() argName
}

type Updater interface {
	updateValue(...string) error
	updateMetadata(...Option)
}

// argName is the private key used to index/identify an Arg
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

func (id argName) Name() string {
	return strings.SplitN(string(id), ":", 2)[1]
}

func (id argName) identify() argName {
	return id
}

// Arg is the shared behaviour of all command-line input types (FlagType, KeyValueType and PositionType). It essentially
// exposes an API similar to the behaviour seen in the standard 'flags' package, extending support to key-value and
// positional arguments.
//
// Metadata is attached to each Arg giving each arg some mutable properties which can be managed via Option(s).
//
// Each Arg has its own go-variableParser which it is responsible for updating:
//
//	var (
//		strVal string
//		myArg = NewKeyValue(&strVal, "arg1", ...)
//		parser = clap.Must("program-name", clap.ContinueOnError).
//				AddKeyValue(myArg).
//				OK()
//	)
//
//	...
//
//	func main() {
//		parser.Parse(os.Args)
//
//		... // do stuff with 'strVal'
//	}
type Arg interface {
	Identifier

	Updater

	// Name returns the key/identifier of the Arg, this should be unique for each Type.
	//
	// This usually be used to identify the Arg from the command-line
	Name() string

	// Type indicates the type of command-line argument the Arg is, returns one of:
	//
	// 	> FlagType;
	//
	//	> KeyValueType;
	//
	//	> PositionType;
	Type() Type

	// Shorthand is the single character alias/identifier for the Arg, if applicable
	//
	// Can be updated via the WithShorthand Option
	Shorthand() rune

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

// TypedArg is the generic interface that is satisfied by all Arg implementations.
//
// This is a convenience interface that provides access to the underlying generic Variable.
type TypedArg[T any] interface {
	Arg
	// Variable returns the underlying Variable for the Arg
	Variable() Variable[T]
}

// ValidateDefaultValue is a helper function for retrieving and attempting to parse the actual typed default value
// of an Arg.
func ValidateDefaultValue[T any](arg TypedArg[T]) (T, error) {

	var zero T
	switch a := arg.(type) {
	case IsFlag:
		switch a.(type) {
		case *FlagArg[bool], *FlagArg[Counter]:
			if !a.HasDefault() {
				return zero, nil
			}
		}
	}

	v := arg.Variable()

	parser := v.parser()

	defaultValue := arg.Default()

	val, err := parser(defaultValue)

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
