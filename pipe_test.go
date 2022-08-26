package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIdentifier_Pipe(t *testing.T) {
	fakeName := "dummy"
	dummyArg := &PipeArg[any]{}

	dummyID := Pipe(fakeName)

	assert.Equal(t, dummyArg.argName(), dummyID.argName(),
		"PipeArg and its Identifier: Pipe do not return equivalent values!")
}

func TestPipeArg_Constructors(t *testing.T) {

	var testCases = map[string]struct {
		Options []Option
		v       int
	}{
		"NoOptions": {},
		"WithDefaultIsNoOp": {
			Options: []Option{
				WithDefault("999"),
			},
		},
		"WithShorthandIsNoOp": {
			Options: []Option{
				WithAlias("t"),
			},
		},
		"AsRequiredIsNoOp": {
			Options: []Option{
				AsRequired(),
			},
		},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			p := NewPipeArg[int](&tc.v, parse.Int{}, nil, os.Stdin, tc.Options...)

			a.False(p.md.HasDefault(), "PipeArg cannot have a default value")
			a.Equal("", p.Default(), "PipeArg default value should be an empty string")
			a.Equalf("", p.Shorthand(), "PipeArg shorthand should be ''")
			a.False(p.IsRequired(), "PipeArg should not be required")

		})
	}

}
