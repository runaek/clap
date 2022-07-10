package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_Flag(t *testing.T) {
	fakeName := "dummy"
	dummyFlag := &FlagArg[any]{
		Key: fakeName,
	}

	dummyID := Flag(fakeName)

	assert.Equal(t, dummyFlag.argName(), dummyID.argName(),
		"FlagArg and its Identifier: Flag do not return equivalent values!")
}

func TestFlagArg_Constructors(t *testing.T) {

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

			f := NewFlag[string](&tc.v, "testflag", parse.String{}, tc.Options...)
			a.False(f.IsRepeatable(), "Singular Flag should not be repeatable")
			a.Equal(testMd.Shorthand(), f.Shorthand(), "Flag: Shorthand does not match expected value")
			a.Equal(testMd.Default(), f.Default(), "Flag: Default does not match expected value")
			a.Equal(testMd.IsRequired(), f.IsRequired(), "Flag: Required option does not match expected value")

			fs := NewFlags[string](&tc.vs, "testflags", parse.String{}, tc.Options...)
			a.True(fs.IsRepeatable(), "Multiple Flags should always be repeatable")
			a.Equal(testMd.Shorthand(), fs.Shorthand(), "Flags: Shorthand does not match expected value")
			a.Equal(testMd.Default(), fs.Default(), "Flags: Default does not match expected value")
			a.Equal(testMd.IsRequired(), f.IsRequired(), "Flags: Required option does not match expected value")

		})
	}
}
