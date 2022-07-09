package clap

import (
	"errors"
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ParserTestCase func(suite *ParserTestSuite)

func NewParserTestSuite(tc ParserTestCase) ParserTestSuite {
	return ParserTestSuite{
		TC:             tc,
		expectedStates: map[argName]error{},
	}
}

type ParserTestSuite struct {
	TC ParserTestCase

	input          []string
	setup          func(*Parser)
	args           []Arg
	expectedStates map[argName]error
	assertions     []func(*assert.Assertions)
	expectValidErr bool
	targetValidErr error
	expectParseErr bool
	targetParseErr error
}

func (tc *ParserTestSuite) GivenParser(cfg func(*Parser)) *ParserTestSuite {
	tc.setup = cfg
	return tc
}

func (tc *ParserTestSuite) GivenArgAddedToParser(args ...Arg) *ParserTestSuite {
	tc.args = append(tc.args, args...)
	return tc
}

func (tc *ParserTestSuite) WhenParsedWithInput(input ...string) *ParserTestSuite {
	tc.input = input
	return tc
}

func (tc *ParserTestSuite) ThenValidShouldReturn(e error) *ParserTestSuite {
	tc.expectValidErr = true
	tc.targetValidErr = e
	return tc
}

func (tc *ParserTestSuite) ThenErrShouldReturn(e error) *ParserTestSuite {
	tc.expectParseErr = true
	tc.targetParseErr = e
	return tc
}

func (tc *ParserTestSuite) ThenTheStateForIdentifierShouldBe(id Identifier, state error) *ParserTestSuite {
	tc.expectedStates[id.argName()] = state
	return tc
}

func (tc *ParserTestSuite) ThenTheParserShould(p *Parser) *ParserTestSuite {

}

func (tc *ParserTestSuite) Then(assrt func(assertions *assert.Assertions)) *ParserTestSuite {
	tc.assertions = append(tc.assertions, assrt)
	return tc
}

func (tc ParserTestSuite) Execute(t *testing.T) {
	tc.TC(&tc)

	a := assert.New(t)
	sut := New("test-parser")

	if tc.setup != nil {
		tc.setup(sut)
	}

	sut.Add(tc.args...)

	if !tc.expectValidErr {
		a.NoError(sut.Valid(), "unexpected validation error")
	} else if tc.targetValidErr != nil {
		a.Truef(errors.Is(sut.Valid(), tc.targetValidErr), "unexpected validation error: expected to detect %s from error but got %s", tc.targetValidErr, sut.Valid())
	} else {
		a.NotNilf(sut.Valid(), "expected a non-nil error")
	}

	sut.Parse(tc.input)

	if !tc.expectParseErr {
		a.NoError(sut.Err(), "unexpected parser error")
	} else if tc.targetParseErr != nil {
		a.Truef(errors.Is(sut.Err(), tc.targetParseErr), "unexpected parser error: expected to detect %s from error but got %s", tc.targetValidErr, sut.Valid())
	} else {
		a.NotNilf(sut.Err(), "expected non-nil error")
	}

	for k, expState := range tc.expectedStates {

		actState := sut.State(k)

		a.Truef(errors.Is(actState, expState), "unexpected state (%s): expected to detect %s but got %s",
			k, expState, actState)
	}

	for _, as := range tc.assertions {
		as(a)
	}
}

func TestParser_Parse(t *testing.T) {

	var testCases = map[string]ParserTestSuite{
		"Basic": NewParserTestSuite(func(suite *ParserTestSuite) {

			var (
				n1        = NewPosition[string](nil, 1, parse.String{})
				n2        = NewPosition[int](nil, 2, parse.Int{})
				remainder = NewPositions[string](nil, 3, parse.String{})
				k1        = NewKeyValue[string](nil, "key1", parse.String{})
				k2        = NewKeyValue[string](nil, "key2", parse.String{})
			)

			suite.
				GivenArgAddedToParser(n1, n2, remainder, k1, k2).
				WhenParsedWithInput("cmd", "1234", "arg3", "key2=hello", "key1=world", "arg4").
				ThenTheStateForIdentifierShouldBe(n1, nil).
				ThenTheStateForIdentifierShouldBe(n2, nil).
				ThenTheStateForIdentifierShouldBe(remainder, nil).
				ThenTheStateForIdentifierShouldBe(k1, nil).
				ThenTheStateForIdentifierShouldBe(k2, nil).
				Then(func(a *assert.Assertions) {

					a.Equal("cmd", n1.Variable().Unwrap())
					a.Equal(1234, n2.Variable().Unwrap())
					a.Equal([]string{"arg3", "arg4"}, remainder.Variable().Unwrap())
					a.Equal("hello", k2.Variable().Unwrap())
					a.Equal("world", k1.Variable().Unwrap())
				})
		}),
		"SubsequentCallsAreSame": NewParserTestSuite(func(suite *ParserTestSuite) {
			var (
				n1        = NewPosition[string](nil, 1, parse.String{})
				n2        = NewPosition[int](nil, 2, parse.Int{})
				remainder = NewPositions[string](nil, 3, parse.String{})
				k1        = NewKeyValue[string](nil, "key1", parse.String{})
				k2        = NewKeyValue[string](nil, "key2", parse.String{})
			)

			suite.
				GivenArgAddedToParser(n1, n2, remainder, k1, k2).
				WhenParsedWithInput("cmd", "1234", "arg3", "key2=hello", "key1=world", "arg4").
				WhenParsedWithInput("cmd", "1234", "arg3", "key2=hello", "key1=world", "arg4").
				ThenTheStateForIdentifierShouldBe(n1, nil).
				ThenTheStateForIdentifierShouldBe(n2, nil).
				ThenTheStateForIdentifierShouldBe(remainder, nil).
				ThenTheStateForIdentifierShouldBe(k1, nil).
				ThenTheStateForIdentifierShouldBe(k2, nil).
				Then(func(a *assert.Assertions) {
					a.Equal("cmd", n1.Variable().Unwrap())
					a.Equal(1234, n2.Variable().Unwrap())
					a.Equal([]string{"arg3", "arg4"}, remainder.Variable().Unwrap())
					a.Equal("hello", k2.Variable().Unwrap())
					a.Equal("world", k1.Variable().Unwrap())
				})
		}),
	}

	for name, tc := range testCases {
		t.Run(name, tc.Execute)
	}
}
