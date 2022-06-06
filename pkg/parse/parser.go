package parse

// A Parser is responsible for parsing some type T from string input.
type Parser[T any] interface {
	// Parse should attempt to parse a new variable from the given input, returning an informative error on failure
	Parse(input ...string) (T, error)
}

// A FuncTypeParser is a convenient Parser[T] implementation.
type FuncTypeParser[T any] func(input ...string) (output T, err error)

func (fn FuncTypeParser[T]) Parse(input ...string) (T, error) {
	return fn(input...)
}

// SliceParser is a helper function to convert a FuncTypeParser[T] into a FuncTypeParser[[]T].
func SliceParser[T any](p Parser[T]) Parser[[]T] {
	return FuncTypeParser[[]T](func(input ...string) ([]T, error) {

		out := make([]T, len(input))

		for i, in := range input {
			if v, err := p.Parse(in); err != nil {
				return out, err
			} else {
				out[i] = v
			}
		}
		return out, nil
	})
}

// AnonymizeParser is a helper function for converting a typed FuncTypeParser into an untyped FuncTypeParser.
func AnonymizeParser[T any](parser Parser[T]) Parser[any] {
	return FuncTypeParser[any](func(input ...string) (output any, err error) {
		return parser.Parse(input...)
	})
}
