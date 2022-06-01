package clap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_Flag(t *testing.T) {
	fakeName := "dummy"
	dummyFlag := &FlagArg[any]{
		Key: fakeName,
	}

	dummyID := Flag(fakeName)

	assert.Equal(t, dummyFlag.identify(), dummyID.identify(),
		"FlagArg and its Identifier: Flag do not return equivalent values!")
}
