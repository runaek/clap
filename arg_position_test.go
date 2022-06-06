package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifier_Positional(t *testing.T) {
	fakeNumber := 7
	dummyFlag := &PositionalArg[any]{
		argIndex: fakeNumber,
	}

	dummyID := Position(fakeNumber)

	assert.Equal(t, dummyFlag.argName(), dummyID.argName(),
		"PositionalArg and its Identifier: Position do not return equivalent values!")
}

func TestPositionalArg_Constructors(t *testing.T) {

	var testCases = map[string]struct {
		Options []Option
		v       int
		vs      []int
	}{
		"NoOptions": {},
		"WithDefaultIsNoOp": {
			Options: []Option{
				WithDefault("999"),
			},
		},
		"WithShorthandIsNoOp": {
			Options: []Option{
				WithShorthand('t'),
			},
		},
		"AsRequiredIsNoOpForVariadicPositions": {
			Options: []Option{
				AsRequired(),
			},
		},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			p := NewPosition[int](&tc.v, 1, parse.Int{}, tc.Options...)

			a.False(p.md.HasDefault(), "PositionalArg should never have a default value")
			a.Equal("", p.Default(), "PositionalArg default value should be an empty string")
			a.Equalf(noShorthand, p.Shorthand(), "PositionalArg shorthand should be %d", noShorthand)

			ps := NewPositions[int](&tc.vs, 2, parse.Int{}, tc.Options...)

			a.False(p.md.HasDefault(), "PositionalArg should never have a default value")
			a.Equal("", p.Default(), "PositionalArg default value should be an empty string")
			a.Equalf(noShorthand, p.Shorthand(), "PositionalArg shorthand should be %d", noShorthand)
			a.False(ps.IsRequired(), "PositionalArgs should not be required")
		})
	}

}
