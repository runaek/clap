package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_KeyValue(t *testing.T) {
	fakeName := "dummy"
	dummyArg := &KeyValueArg[any]{
		Key: fakeName,
	}

	dummyID := Key(fakeName)

	assert.Equal(t, dummyArg.argName(), dummyID.argName(),
		"KeyValueArg and its Identifier: Key do not return equivalent values!")
}

func TestKeyValueArg_Constructors(t *testing.T) {

	var testCases = map[string]struct {
		Options []Option
		v       string
		vs      []string
	}{
		"NoOptions": {},
		"ShorthandApplies": {
			Options: []Option{WithShorthand("t")},
		},
		"DefaultApplies": {
			Options: []Option{
				WithDefault("my-default"),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)

			testMd := NewMetadata(tc.Options...)

			kv := NewKeyValue[string](&tc.v, "testarg", parse.String{}, tc.Options...)
			a.False(kv.IsRepeatable(), "Singular KeyValue should not be repeatable")
			a.Equal(testMd.Shorthand(), kv.Shorthand(), "KeyValue: Shorthand does not match expected value")
			a.Equal(testMd.Default(), kv.Default(), "KeyValue: Default does not match expected value")
			a.Equal(testMd.IsRequired(), kv.IsRequired(), "KeyValue: Required option does not match expected value")

			kvs := NewKeyValues[string](&tc.vs, "testmultiargs", parse.String{}, tc.Options...)
			a.True(kvs.IsRepeatable(), "Multiple KeyValue(s) should always be repeatable")
			a.Equal(testMd.Shorthand(), kvs.Shorthand(), "KeyValues: Shorthand does not match expected value")
			a.Equal(testMd.Default(), kvs.Default(), "KeyValues: Default does not match expected value")
			a.Equal(testMd.IsRequired(), kv.IsRequired(), "KeyValues: Required option does not match expected value")

		})
	}
}
