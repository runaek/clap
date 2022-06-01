package clap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_Positional(t *testing.T) {
	fakeNumber := 7
	dummyFlag := &PositionalArg[any]{
		argIndex: fakeNumber,
	}

	dummyID := Position(fakeNumber)

	assert.Equal(t, dummyFlag.identify(), dummyID.identify(),
		"PositionalArg and its Identifier: Position do not return equivalent values!")
}
