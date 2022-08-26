package parse

// A Parser is responsible for parsing some type T from string input.
type Parser[T any] interface {
	// Parse should attempt to parse a new variable from the given input, returning an informative error on failure
	Parse(input ...string) (T, error)
}

// A Func is a convenient Parser[T] implementation.
type Func[T any] func(input ...string) (output T, err error)

func (fn Func[T]) Parse(input ...string) (T, error) {
	return fn(input...)
}

// Slice is a helper function to convert a Parser[T] into a Parser[[]T].
//
// Return type is a Func[[]T].
func Slice[T any](p Parser[T]) Parser[[]T] {
	return Func[[]T](func(input ...string) ([]T, error) {
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

// Anonymize is a helper function for converting a typed Parser[T] into an untyped Parser[any].
//
// Returned type is a Func[any].
func Anonymize[T any](parser Parser[T]) Parser[any] {
	return Func[any](func(input ...string) (output any, err error) {
		return parser.Parse(input...)
	})
}
