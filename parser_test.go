package clap

import (
	"errors"
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type MockParserTestCase func(suite *TestMockParser)

func NewParserTestSuite(tc MockParserTestCase) TestMockParser {
	return TestMockParser{
		TC:             tc,
		expectedStates: map[argName]error{},
	}
}

// TestMockParser is a helper for testing the Parser with different Arg(s) and input(s).
type TestMockParser struct {
	TC MockParserTestCase

	input          []string
	setup          func(*Parser)
	args           []Arg
	expectedStates map[argName]error
	assertions     []func(*Parser, *assert.Assertions)
	expectValidErr bool
	targetValidErr error
	expectParseErr bool
	targetParseErr error
}

func (tc *TestMockParser) GivenParser(cfg func(*Parser)) *TestMockParser {
	tc.setup = cfg
	return tc
}

func (tc *TestMockParser) GivenArgAddedToParser(args ...Arg) *TestMockParser {
	tc.args = append(tc.args, args...)
	return tc
}

func (tc *TestMockParser) WhenParsedWithInput(input ...string) *TestMockParser {
	tc.input = input
	return tc
}

func (tc *TestMockParser) ThenValidShouldReturn(e error) *TestMockParser {
	tc.expectValidErr = true
	tc.targetValidErr = e
	return tc
}

func (tc *TestMockParser) ThenErrShouldReturn(e error) *TestMockParser {
	tc.expectParseErr = true
	tc.targetParseErr = e
	return tc
}

func (tc *TestMockParser) ThenTheStateForIdentifierShouldBe(id Identifier, state error) *TestMockParser {
	tc.expectedStates[id.argName()] = state
	return tc
}

func (tc *TestMockParser) Then(assrt func(p *Parser, assertions *assert.Assertions)) *TestMockParser {
	tc.assertions = append(tc.assertions, assrt)
	return tc
}

func (tc TestMockParser) Run(t *testing.T) {
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
		as(sut, a)
	}
}

func TestParser_Parse(t *testing.T) {

	// for some reason, GoLand does not recognise the TestMockParser as a test
	// unless it is wrapped in a struct like this
	var testCases = map[string]struct {
		TestMockParser
	}{
		"ValidArgs": {
			TestMockParser: NewParserTestSuite(func(suite *TestMockParser) {

				var (
					n1        = NewPosition[string](nil, 1, parse.String{})
					n2        = NewPosition[int](nil, 2, parse.Int{})
					remainder = NewPositions[string](nil, 3, parse.String{})
					k1        = NewKeyValue[string](nil, "key1", parse.String{})
					k2        = NewKeyValue[string](nil, "key2", parse.String{})
					f1        = NewFlagsP[string](nil, "f1", "f", parse.String{})
					//fb1       = NewFlagP[bool](nil, "debug", "D", parse.Bool{}, WithDefault("false"))
				)

				suite.
					GivenArgAddedToParser(n1, n2, remainder, k1, k2, f1).
					WhenParsedWithInput("cmd", "1234", "arg3", "-f=v1", "key2=hello", "-D", "key1=world", "arg4", "-f", "v2").
					ThenTheStateForIdentifierShouldBe(n1, nil).
					ThenTheStateForIdentifierShouldBe(n2, nil).
					ThenTheStateForIdentifierShouldBe(remainder, nil).
					ThenTheStateForIdentifierShouldBe(k1, nil).
					ThenTheStateForIdentifierShouldBe(k2, nil).
					Then(func(p *Parser, a *assert.Assertions) {

						a.Equal("cmd", n1.Variable().Unwrap())
						a.Equal(1234, n2.Variable().Unwrap())
						a.Equal([]string{"arg3", "arg4"}, remainder.Variable().Unwrap())
						a.Equal("hello", k2.Variable().Unwrap())
						a.Equal("world", k1.Variable().Unwrap())
						a.Equal([]string{"v1", "v2"}, f1.Variable().Unwrap())
					})
			}),
		},
		"MissingArgs": {
			TestMockParser: NewParserTestSuite(func(suite *TestMockParser) {

				var (
					n1        = NewPosition[string](nil, 1, parse.String{})
					n2        = NewPosition[int](nil, 2, parse.Int{})
					remainder = NewPositions[string](nil, 3, parse.String{})
					k1        = NewKeyValue[string](nil, "key1", parse.String{}, AsRequired())
					k2        = NewKeyValue[string](nil, "key2", parse.String{})
				)

				suite.
					GivenArgAddedToParser(n1, n2, remainder, k1, k2).
					WhenParsedWithInput("cmd", "1234", "arg3", "key2=hello", "arg4").
					ThenErrShouldReturn(ErrMissing).
					ThenTheStateForIdentifierShouldBe(n1, nil).
					ThenTheStateForIdentifierShouldBe(n2, nil).
					ThenTheStateForIdentifierShouldBe(remainder, nil).
					ThenTheStateForIdentifierShouldBe(k1, ErrMissing).
					ThenTheStateForIdentifierShouldBe(k2, nil).
					Then(func(p *Parser, a *assert.Assertions) {

						a.Equal("cmd", n1.Variable().Unwrap())
						a.Equal(1234, n2.Variable().Unwrap())
						a.Equal([]string{"arg3", "arg4"}, remainder.Variable().Unwrap())
						a.Equal("hello", k2.Variable().Unwrap())
						a.Equal("", k1.Variable().Unwrap())
					})
			}),
		},
		"CombinedBooleans": {
			TestMockParser: NewParserTestSuite(func(suite *TestMockParser) {

				var (
					f1 = NewFlagP[bool](nil, "a1", "a", parse.Bool{})
					f2 = NewFlagP[bool](nil, "b1", "b", parse.Bool{})
					f3 = NewFlagP[bool](nil, "c1", "c", parse.Bool{}, WithDefault("true"))
					f4 = NewFlagP[bool](nil, "d1", "d", parse.Bool{})
				)

				suite.
					GivenArgAddedToParser(f1, f2, f3, f4).
					WhenParsedWithInput("cmd", "1234", "--d1", "key2=hello", "-abc").
					ThenTheStateForIdentifierShouldBe(f1, nil).
					ThenTheStateForIdentifierShouldBe(f2, nil).
					ThenTheStateForIdentifierShouldBe(f3, nil).
					ThenTheStateForIdentifierShouldBe(f4, nil).
					Then(func(p *Parser, a *assert.Assertions) {

						a.Equal(true, f1.Variable().Unwrap())
						a.Equal(true, f2.Variable().Unwrap())
						a.Equal(false, f3.Variable().Unwrap())
						a.Equal(true, f4.Variable().Unwrap())
					})
			}),
		},
		"DoubleCall": {
			TestMockParser: NewParserTestSuite(func(suite *TestMockParser) {

				var (
					n1        = NewPosition[string](nil, 1, parse.String{})
					n2        = NewPosition[int](nil, 2, parse.Int{})
					remainder = NewPositions[string](nil, 3, parse.String{})
					k2        = NewKeyValue[string](nil, "key2", parse.String{})
				)

				suite.
					GivenArgAddedToParser(n1, n2, remainder, k2).
					WhenParsedWithInput("cmd", "1234", "arg3", "key2=hello", "arg4").
					ThenTheStateForIdentifierShouldBe(n1, nil).
					ThenTheStateForIdentifierShouldBe(n2, nil).
					ThenTheStateForIdentifierShouldBe(remainder, nil).
					ThenTheStateForIdentifierShouldBe(k2, nil).
					Then(func(p *Parser, a *assert.Assertions) {
						a.Equal("cmd", n1.Variable().Unwrap())
						a.Equal(1234, n2.Variable().Unwrap())
						a.Equal([]string{"arg3", "arg4"}, remainder.Variable().Unwrap())
						a.Equal("hello", k2.Variable().Unwrap())

						p.Parse([]string{"1234", "arg3", "key2=hello", "arg4"})

						a.Equal(p.positionalValues, []string{"1234", "arg3", "arg4"})
					})
			}),
		},
		"StrictModeErrorsButValuesParse": {
			TestMockParser: NewParserTestSuite(func(suite *TestMockParser) {

				var (
					f1 = NewFlagP[bool](nil, "a1", "a", parse.Bool{})
					f2 = NewFlagP[bool](nil, "b1", "b", parse.Bool{})
					f3 = NewFlagP[bool](nil, "c1", "c", parse.Bool{}, WithDefault("true"))
					f4 = NewFlagP[bool](nil, "d1", "d", parse.Bool{})
				)

				suite.
					GivenParser(func(parser *Parser) {
						parser.Strict = true
					}).
					GivenArgAddedToParser(f1, f2, f3, f4).
					WhenParsedWithInput("cmd", "1234", "--d1", "key2=hello", "-abc").
					ThenErrShouldReturn(ErrUnidentified).
					ThenTheStateForIdentifierShouldBe(f1, nil).
					ThenTheStateForIdentifierShouldBe(f2, nil).
					ThenTheStateForIdentifierShouldBe(f3, nil).
					ThenTheStateForIdentifierShouldBe(f4, nil).
					Then(func(p *Parser, a *assert.Assertions) {

						a.Equal(true, f1.Variable().Unwrap())
						a.Equal(true, f2.Variable().Unwrap())
						a.Equal(false, f3.Variable().Unwrap())
						a.Equal(true, f4.Variable().Unwrap())
					})
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, tc.Run)
	}
}

func TestParser_Usage(t *testing.T) {

	a := assert.New(t)
	p := NewParser("test-parser", ContinueOnError)

	tmpDir := t.TempDir()
	t.Logf("Creating Temp Files @%s", tmpDir)
	inF, inputErr := os.OpenFile(filepath.Join(tmpDir, "usage_input"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	a.NoError(inputErr)
	outF, outputErr := os.OpenFile(filepath.Join(tmpDir, "usage_output"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	a.NoError(outputErr)

	defer func() {
		inF.Close()
		outF.Close()
	}()

	p.Stderr = outF
	p.Stdout = outF
	p.Stdin = inF

	var (
		n1        = NewPosition[string](nil, 1, parse.String{})
		remainder = NewPositions[string](nil, 2, parse.String{})
		k1        = NewKeyValue[string](nil, "key1", parse.String{})
		f1        = NewFlagsP[string](nil, "f1", "f", parse.String{})
	)

	p.Add(n1, remainder, k1, f1)

	p.Usage()

	a.NoError(outF.Close())

	outF2, err := os.OpenFile(filepath.Join(tmpDir, "usage_output"), os.O_RDONLY, os.ModePerm)
	a.NoError(err)
	data, err := ioutil.ReadAll(outF2)

	a.NoError(err)

	a.Equal(`Usage: test-parser [ <args>, ... ] [ <key>=<value>, ... ] [ --<flag>=<value>, ... ]

ARGUMENTS
	1     : A positional argument.
	2 ... : (repeatable) Remaining positional arguments.

OPTIONS (<key>=<value>)
	     key1                : key1 - a string key-value variable.

FLAGS
	[-f]  --f1               : f1 - a []string flag variable.
`, string(data))

}
