package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMetadata(t *testing.T) {

	var testCases = map[string]struct {
		Input   Arg
		Updates []Option
		Output  Metadata
	}{
		"ChangeToRequired": {
			Input: NewFlag[string](nil, "f1", parse.String{}),
			Updates: []Option{
				AsRequired(),
			},
			Output: Metadata{
				argUsage:     "f1 - a string flag variable.",
				argShorthand: "",
				argDefault:   "",
				hasDefault:   false,
				argRequired:  true,
			},
		},
		"ChangeDefault": {
			Input: NewFlag[string](nil, "f1", parse.String{}, WithDefault("first_default")),
			Updates: []Option{
				WithDefault("second_default"),
			},
			Output: Metadata{
				argUsage:     "f1 - a string flag variable.",
				argShorthand: "",
				argDefault:   "second_default",
				hasDefault:   true,
				argRequired:  false,
			},
		},
		"AddShorthand": {
			Input: NewFlag[string](nil, "f1", parse.String{}, WithDefault("first_default")),
			Updates: []Option{
				WithDefault("second_default"),
				WithAlias("f"),
			},
			Output: Metadata{
				argUsage:     "f1 - a string flag variable.",
				argShorthand: "f",
				argDefault:   "second_default",
				hasDefault:   true,
				argRequired:  false,
			},
		},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {
			a := assert.New(t)

			UpdateMetadata(tc.Input, tc.Updates...)

			a.Equal(tc.Output.Usage(), tc.Input.Usage())
			a.Equal(tc.Output.IsRequired(), tc.Input.IsRequired())

		})
	}

}
