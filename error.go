package clap

import (
	"errors"
	"fmt"
	"strings"
)

func ErrScanning(at int, cause error, tokens ...string) *ScanError {
	return &ScanError{
		Index:  at,
		Tokens: tokens,
		Cause:  cause,
	}
}

// ScanError indicates that there was an error decoding the command-line input tokens.
type ScanError struct {
	Index  int
	Tokens []string
	Cause  error
}

func (err *ScanError) Error() string {
	return fmt.Sprintf("error scanning %q (index=%d): %s", err.Tokens, err.Index, err.Cause)
}

func (err *ScanError) Unwrap() error {
	return err.Cause
}

func ErrParsing(id Identifier, cause error) *ParseError {
	return &ParseError{
		Id:    id,
		Cause: cause,
	}
}

// ParseError indicates that an error occurred parsing the variable for some argument.
type ParseError struct {
	Id    Identifier
	Cause error
}

func (err *ParseError) Unwrap() error {
	return err.Cause
}

func (err *ParseError) Error() string {
	a := err.Id.identify()
	return fmt.Sprintf("error parsing %s argument %q: %s", a.Type(), a.Name(), err.Cause)
}

var (
	ErrUnrecognisedType  = errors.New("unable to determine argument type")
	ErrUnrecognisedToken = errors.New("unrecognised token")
	ErrIncompleteToken   = errors.New("argument is missing its value")
	WarnUnrecognisedArg  = errors.New("unrecognised argument name")
)

// ErrMissing is a constructor for an *empty* MissingInputError.
//
// ErrMissing *always* returns a non-nil *MissingInputError - it is on the caller to add errors to this  and make it
// meaningful (i.e. it will still be considered an error if there are no ids 'missing' - use IsEmpty to check)
func ErrMissing(ids ...Identifier) *MissingInputError {
	err := &MissingInputError{
		missing: map[argName]struct{}{},
	}

	err.Add(ids...)

	return err
}

// MissingInputError indicates that there are missing argument(s).
type MissingInputError struct {
	missing map[argName]struct{}
}

func (err *MissingInputError) IsEmpty() bool {
	return len(err.missing) == 0
}

func (err *MissingInputError) Add(ids ...Identifier) {
	if err.missing == nil {
		err.missing = map[argName]struct{}{}
	}

	for _, id := range ids {
		err.missing[id.identify()] = struct{}{}
	}
}

func (err *MissingInputError) Error() string {
	hdr := fmt.Sprintf("%d missing argument(s):", len(err.missing))

	msgs := make([]string, len(err.missing))
	i := 0
	for a, _ := range err.missing {
		msgs[i] = fmt.Sprintf("\t> (%s) %s", a.Type(), a.Name())
		i++
	}

	return fmt.Sprintf("%s\n%s", hdr, strings.Join(msgs, "\n"))
}

func (err *MissingInputError) IsMissing(id Identifier) bool {
	_, missing := err.missing[id.identify()]
	return missing
}
