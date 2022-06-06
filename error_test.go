package clap

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Is(t *testing.T) {

	var testCases = map[string]struct {
		Err    error
		Target Error
		Match  bool
	}{
		"NotError": {
			Err:    fmt.Errorf("some error"),
			Target: ErrMissing,
			Match:  false,
		},
		"CausedScanErr": {
			Err:    ErrScanning(fmt.Errorf("some err: %w", ErrMissing), "mock", "token"),
			Target: ErrMissing,
			Match:  true,
		},
		"DidNotCauseScanErr": {
			Err:    ErrScanning(fmt.Errorf("some err: %w", ErrUnidentified), "mock", "token"),
			Target: ErrMissing,
		},
		"ErrMissingMatch": {
			Err:    fmt.Errorf("should be missing: %w", ErrMissing),
			Target: ErrMissing,
			Match:  true,
		},
		"ErrMissingDoesNotMatchMatch": {
			Err:    fmt.Errorf("should be missing: %w", ErrDuplicate),
			Target: ErrMissing,
		},
		"MultiErrDetectErrMissing": {
			Err:    &multierror.Error{Errors: []error{ErrMissing, ErrDuplicate}},
			Target: ErrMissing,
			Match:  true,
		},
		"MultiErrDetectErrDuplicate": {
			Err:    &multierror.Error{Errors: []error{ErrMissing, ErrDuplicate}},
			Target: ErrDuplicate,
			Match:  true,
		},
		"MultiErrShouldNotDetectErr": {
			Err:    &multierror.Error{Errors: []error{ErrMissing, ErrDuplicate}},
			Target: ErrInvalid,
		},

		"WMultiErrDetectErrMissing": {
			Err:    &multierror.Error{Errors: []error{fmt.Errorf("some err: %w", ErrMissing), ErrDuplicate}},
			Target: ErrMissing,
			Match:  true,
		},
		"WMultiErrDetectErrDuplicate": {
			Err:    &multierror.Error{Errors: []error{ErrMissing, fmt.Errorf("some err: %w", ErrDuplicate)}},
			Target: ErrDuplicate,
			Match:  true,
		},
		"WMultiErrShouldNotDetectErr": {
			Err:    &multierror.Error{Errors: []error{ErrMissing, ErrDuplicate}},
			Target: ErrInvalid,
		},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			match := errors.Is(tc.Err, tc.Target)
			a.Equal(tc.Match, match)
		})
	}

}
