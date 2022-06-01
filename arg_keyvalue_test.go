package clap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_KeyValue(t *testing.T) {
	fakeName := "dummy"
	dummyArg := &KeyValueArg[any]{
		Key: fakeName,
	}

	dummyID := Key(fakeName)

	assert.Equal(t, dummyArg.identify(), dummyID.identify(),
		"KeyValueArg and its Identifier: Key do not return equivalent values!")
}
