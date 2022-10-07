package clap

import (
	"fmt"
	"strings"
)

const (
	ErrUnidentified Error = "unidentified argument"
	ErrMissing      Error = "argument is missing"
	ErrDuplicate    Error = "identifier already exists"
	ErrInvalidIndex Error = "invalid positional index"
	ErrPipeExists   Error = "pipe already exists"
	ErrUnknownType  Error = "unrecognised argument type"
	ErrInvalid      Error = "invalid argument syntax"
	ErrHelp         Error = "help requested"
)

// Error is a simple indicator for some error that occurs.
//
// The value of the error should indicate the 'cause' of the problem and context
// should be provided by the returning process.
type Error string

func (err Error) Error() string {
	return string(err)
}

func (err Error) Is(target error) bool {
	parseErr, ok := target.(Error) // nolint: errorlint

	if !ok {
		return false
	}

	return parseErr == err
}

// ErrScanning is a constructor for a ScanError.
func ErrScanning(cause error, tokens ...string) *ScanError {
	return &ScanError{
		Tokens: tokens,
		Cause:  cause,
	}
}

// ScanError indicates that there was an error scanning the command-line inpu.
type ScanError struct {
	Tokens []string // Tokens are the offending tokens
	Cause  error    // Cause is the error that the tokens lead to
}

func (err *ScanError) Error() string {
	return fmt.Sprintf("error scanning %q: %s", strings.Join(err.Tokens, " "), err.Cause)
}

func (err *ScanError) Unwrap() error {
	return err.Cause
}

// ErrParsing is a constructor for a ParseError.
func ErrParsing(id Identifier, cause error) *ParseError {
	return &ParseError{
		Id:    id,
		Cause: cause,
	}
}

// ParseError indicates that an error occurred parsing the variable for some
// argument.
type ParseError struct {
	Id    Identifier
	Cause error
}

func (err *ParseError) Unwrap() error {
	return err.Cause
}

func (err *ParseError) Error() string {
	a := err.Id.argName()

	return fmt.Sprintf("error parsing (%s) %q): %s", a.Type(), a.Name(), err.Cause)
}
