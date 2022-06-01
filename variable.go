package clap

// A Counter is a special type that tracks the number of times a particular FlagArg is supplied, expecting no value to
// be supplied - like a boolean flag.
type Counter int

func CounterParser(input ...string) (Counter, error) {
	return Counter(len(input)), nil
}

// A Variable refers to some program variable that can be parsed from string input.
//
// TODO: does this even need to be an interface?
type Variable[T any] interface {
	// Update attempts to parse a new value for the Variable and update it
	Update(...string) error

	// Ref returns a reference to the underlying Variable
	Ref() *T
	// Unwrap returns the value of the underlying Variable
	Unwrap() T

	// Parser returns the ValueParser for the type T
	parser() ValueParser[T]
}

// A ValueParser is responsible for parsing some type T from some string-input.
type ValueParser[T any] func(input ...string) (output T, err error)

// SliceParser is a helper function to convert a ValueParser[T] into a ValueParser[[]T].
func SliceParser[T any](p ValueParser[T]) ValueParser[[]T] {
	return func(input ...string) ([]T, error) {

		out := make([]T, len(input))

		for i, in := range input {

			if v, err := p(in); err != nil {
				return out, err
			} else {
				out[i] = v
			}

		}

		return out, nil
	}

}

// AnonymizeParser is a helper function for converting a typed ValueParser into an untyped ValueParser.
func AnonymizeParser[T any](parser ValueParser[T]) ValueParser[any] {
	return func(input ...string) (output any, err error) {
		return parser(input...)
	}
}

// NewVariables is a constructor for a Variable that has an underlying variableParser of slice-type.
//
// This is a helper function which converts the supplied parser into one that supports slices via SliceParser, specific
// parsers can be used by creating a Variable using NewVariablesWithParser.
func NewVariables[T any](variables *[]T, p ValueParser[T]) Variable[[]T] {
	return &variableParser[[]T]{
		v: variables,
		p: SliceParser(p),
	}
}

// NewVariablesWithParser is a constructor for a Variable that has an underlying variableParser of slice-type.
func NewVariablesWithParser[T any](variables *[]T, p ValueParser[[]T]) Variable[[]T] {
	return &variableParser[[]T]{
		v: variables,
		p: p,
	}
}

// NewVariable is a constructor for a Variable.
func NewVariable[T any](variable *T, p ValueParser[T]) Variable[T] {
	return &variableParser[T]{
		v: variable,
		p: p,
	}
}

// variableParser wraps some variable of type T and is responsible for updating the underlying value of this variable
// when required.
//
// The parser is responsible for the generating the new value (from string input) to be held by the variable.
type variableParser[T any] struct {
	v *T             // v is the underlying go-variable being maintained by the variableParser
	p ValueParser[T] // p is the parser responsible for creating a new value for the go-variable
}

func (b *variableParser[T]) Update(ss ...string) error {

	v, err := b.p(ss...)

	if err != nil {
		return err
	}

	*b.v = v

	return nil
}

func (b *variableParser[T]) Ref() *T {
	return b.v
}

func (b *variableParser[T]) Unwrap() T {
	if b.v == nil {
		return *new(T)
	}
	return *b.v
}

func (b *variableParser[T]) parser() ValueParser[T] {
	return b.p
}

var (
	_ Variable[any] = &variableParser[any]{}
)
