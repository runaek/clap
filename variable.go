package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"io"
	"os"
)

// FileWriter represents some type of file-like object that would be used as an output we can write to.
//
// Usually, in this package, values of this type will be set to os.Stdout by default.
type FileWriter interface {
	io.Writer
	Stat() (os.FileInfo, error)
	Fd() uintptr
}

type FileInfo = os.FileInfo

// FileReader represents some time of file-like object that would be used as an input we can read from.
//
// Usually, in this package, values of this type will be set to os.Stdin by default
type FileReader interface {
	io.Reader
	Stat() (os.FileInfo, error)
	Fd() uintptr
}

// A Variable refers to some program variable that can be parsed from string input.
//
// Internally, a Variable uses a FuncTypeParser to parse some string input into the underlying variable.
type Variable[T any] interface {
	// Update parses the given input and attempts to update the underlying variable
	Update(...string) error
	// Ref returns a reference to the underlying variable
	Ref() *T
	// Unwrap returns the value of the underlying variable
	Unwrap() T
	// Parser returns the FuncTypeParser for the variable
	Parser() parse.Parser[T]
}

// NewVariables is a constructor for a Variable that has an underlying argumentVariable of slice-type.
//
// This is a helper function which converts the supplied FuncTypeParser into one that supports slices via SliceParser.
//
// A specific FuncTypeParser can be used by creating a Variable using NewVariablesWithParser.
func NewVariables[T any](variables *[]T, p parse.Parser[T]) Variable[[]T] {
	if variables == nil {
		var anon []T
		variables = &anon
	}
	return &argumentVariable[[]T]{
		v: variables,
		p: parse.SliceParser(p),
	}
}

// NewVariablesWithParser is a constructor for a Variable that has an underlying FuncTypeParser of slice-type.
func NewVariablesWithParser[T any](variables *[]T, p parse.Parser[[]T]) Variable[[]T] {
	return &argumentVariable[[]T]{
		v: variables,
		p: p,
	}
}

// NewVariable is a constructor for a Variable.
func NewVariable[T any](variable *T, p parse.Parser[T]) Variable[T] {
	if variable == nil {
		variable = new(T)
	}
	return &argumentVariable[T]{
		v: variable,
		p: p,
	}
}

// argumentVariable wraps some variable of type T and is responsible for updating the underlying value of this variable
// when required.
//
// The Parser is responsible for the generating the new value (from string inpu) to be held by the variable.
type argumentVariable[T any] struct {
	v *T              // v is the underlying go-variable being maintained by the argumentVariable
	p parse.Parser[T] // p is the Parser responsible for creating a new value for the go-variable
}

func (b *argumentVariable[T]) Update(ss ...string) error {

	v, err := b.p.Parse(ss...)

	if err != nil {
		return err
	}

	*b.v = v

	return nil
}

func (b *argumentVariable[T]) Ref() *T {
	return b.v
}

func (b *argumentVariable[T]) Unwrap() T {
	if b.v == nil {
		return *new(T)
	}
	return *b.v
}

func (b *argumentVariable[T]) Parser() parse.Parser[T] {
	return b.p
}

var (
	_ Variable[any] = &argumentVariable[any]{}
)
