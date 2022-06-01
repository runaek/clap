package clap

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"go.uber.org/zap"
	"os"
	"strings"
)

type ErrorHandling int

const (
	ContinueOnError = ErrorHandling(flag.ContinueOnError)
	ExitOnError     = ErrorHandling(flag.ExitOnError)
	PanicOnError    = ErrorHandling(flag.PanicOnError)
)

// HandleError handles some error according to the ErrorHandling, returning true if the error is not nil, otherwise false
// is returned.
func HandleError(eh ErrorHandling, e error) bool {
	if e == nil {
		return false
	}

	switch eh {
	case ContinueOnError:
		return true
	case PanicOnError:
		panic(e)
	case ExitOnError:
		_, _ = fmt.Fprint(os.Stderr, e)
		os.Exit(1)
	default:
		panic("invalid error handling method")
	}

	return true
}

// New is a constructor for a new command-line argument Parser.
func New(name string, errHandling ErrorHandling, elements ...any) (*Parser, error) {

	s := newParser(name, errHandling)

	if len(elements) == 0 {
		return s, nil
	}

	return s, s.Err()
}

// Must is a constructor for a new command-line argument Parser which will panic if the constructor fails.
func Must(name string, errHandling ErrorHandling, elements ...any) *Parser {
	p, err := New(name, errHandling, elements...)

	if err != nil {
		panic(err)
	}

	return p
}

func newParser(n string, eh ErrorHandling) *Parser {

	constrErr := new(multierror.Error)
	constrErr.ErrorFormat = func(es []error) string {
		hdr := fmt.Sprintf("caught %d error(s) setting up parser %q:\n\t", len(es), n)

		msgs := make([]string, len(es))

		for i, e := range es {
			msgs[i] = fmt.Sprintf("* %s", e)
		}

		return hdr + strings.Join(msgs, "\n\t")
	}

	return &Parser{
		name:       n,
		eh:         eh,
		Set:        NewSet(),
		flagValues: map[string][]string{},
		keyValues:  map[string][]string{},
		parseErr:   new(multierror.Error),
	}
}

// A Parser Contains and parses a number of flag, key-value or positional inputs from the command-line.
type Parser struct {
	*Set        // underlying Set for the Parser - holding all the Arg types the Parser is responsible for
	Shift  int  // Shift shifts the Parser along, ignoring the first 'shifted' arguments
	Strict bool // Strict defines whether unrecognised tokens (e.g. flag/keys) are ignored, or return UnrecognisedInput

	name             string              // arbitrary name of the Parser - for debugging purposes
	eh               ErrorHandling       // type of error handling for the Parser
	positionalValues []string            // raw positional arguments from the latest Parse call
	keyValues        map[string][]string // raw key-value arguments from the latest Parse call
	flagValues       map[string][]string // raw flag-value arguments from the latest Parse call

	parseErr *multierror.Error // parserErr contains the error(s) for the *latest* call to Parse
	setErr   *multierror.Error // constrErr contains any errors caught whilst adding arguments to the Parser
}

// RawPositions returns the raw ordered positional arguments detected by the Parser.
func (p *Parser) RawPositions() []string {
	return p.positionalValues
}

// RawFlags returns the raw flag arguments and their associated values detected by the Parser.
func (p *Parser) RawFlags() map[string][]string {
	return p.flagValues
}

// RawKeyValues returns the raw key-value arguments detected by the Parser.
func (p *Parser) RawKeyValues() map[string][]string {
	return p.keyValues
}

// Add a number of Arg(s) to the Parser.
func (p *Parser) Add(args ...Arg) *Parser {
	for _, a := range args {
		switch arg := a.(type) {
		case IsFlag:
			_ = p.AddFlag(arg)
		case IsKeyValue:
			_ = p.AddKeyValue(arg)
		case IsPositional:
			_ = p.AddPosition(arg)
		case IsPipe:
			_ = p.AddPipe(arg)
		default:
			log.Warn("Skipping malformed Arg",
				zap.String("_t", fmt.Sprintf("%T", a)),
				zap.String("name", a.Name()),
				zap.Stringer("type", a.Type()))
		}
	}

	return p
}

// AddPosition adds a positional argument to the Parser.
func (p *Parser) AddPosition(a IsPositional, opts ...Option) *Parser {
	if err := p.Set.AddPosition(a, opts...); HandleError(p.eh, err) {
		p.setErr = multierror.Append(p.setErr, err)
	}
	return p
}

// AddFlag adds a flag argument to the Parser.
func (p *Parser) AddFlag(f IsFlag, opts ...Option) *Parser {

	if err := p.Set.AddFlag(f, opts...); HandleError(p.eh, err) {
		p.setErr = multierror.Append(p.setErr, fmt.Errorf("unable to add flag: %w", err))
	}
	return p
}

// AddKeyValue adds a key-value argument to the Parser.
func (p *Parser) AddKeyValue(kv IsKeyValue, opts ...Option) *Parser {
	if err := p.Set.AddKeyValue(kv, opts...); HandleError(p.eh, err) {
		p.setErr = multierror.Append(p.setErr, fmt.Errorf("unable to add key-value: %w", err))

	}
	return p
}

// AddPipe adds a pipe argument to the Parser.
func (p *Parser) AddPipe(pipe IsPipe, options ...Option) *Parser {
	if err := p.Set.AddPipe(pipe, options...); HandleError(p.eh, err) {
		p.setErr = multierror.Append(p.setErr, fmt.Errorf("unable to add position: %w", err))
	}
	return p
}

// Ok is a helper function that panics if there were any Parser errors. This is useful for running at the end of a chain
// of method calls for the Parser.
//
// Returns the Parser for convenience.
func (p *Parser) Ok() *Parser {

	if HandleError(p.eh, p.setErr.ErrorOrNil()) {
		panic(p.setErr.ErrorOrNil())
	}

	if HandleError(p.eh, p.Err()) {
		panic(p.Err())
	}

	return p
}

// Err returns the error for the *latest* call to Parse.
func (p *Parser) Err() error {
	return p.parseErr.ErrorOrNil()
}

// Parse command-line input.
func (p *Parser) Parse(argv []string) {

	if argv == nil {
		argv = os.Args[1:]
	}

	log.Debug("Parsing input", zap.String("parser", p.name), zap.Int("handling", int(p.eh)), zap.Strings("input", argv))
	p.parseErr = new(multierror.Error)
	p.parseErr.ErrorFormat = func(es []error) string {
		hdr := fmt.Sprintf("caught %d error(s) parsing input %q:\n\t", len(es), argv)

		msgs := make([]string, len(es))

		for i, e := range es {
			msgs[i] = fmt.Sprintf("* %s", e)
		}

		return hdr + strings.Join(msgs, "\n\t")
	}

	relativePos := 0

	for tkns, consumed, err := argv, []string{}, error(nil); ; tkns, consumed, err = p.scan(tkns) {

		relativePos += len(consumed)

		if err == finished {
			break
		}
		if HandleError(p.eh, err) {
			scanErr := ErrScanning(relativePos, err, consumed...)
			log.Warn("Error during token scan", zap.Error(scanErr))
			p.parseErr = multierror.Append(p.parseErr, scanErr)
		}
	}

	if pipe := p.Pipe(); pipe != nil {
		if err := pipe.updateValue(); HandleError(p.eh, err) {
			p.parseErr = multierror.Append(p.parseErr, ErrParsing(Pipe(""), err))
		}
	}

	variadicIndex := -1
	var variadicValues []string
	for im1, v := range p.positionalValues {

		pA := p.Pos(im1 + 1)

		if pA == nil && variadicIndex > 0 {
			variadicValues = append(variadicValues, v)
			continue
		} else if pA == nil {
			break
		}

		if pA.IsRepeatable() {
			variadicIndex = im1 + 1
			variadicValues = append(variadicValues, v)
		} else {
			if err := pA.updateValue(v); HandleError(p.eh, err) {
				p.parseErr = multierror.Append(p.parseErr, err)
			}
		}
	}

	if variadicIndex > 0 {
		variadicPosArg := p.Pos(variadicIndex)
		log.Debug("Updating variadic positional arguments", zap.Strings("vals", variadicValues))
		if err := variadicPosArg.updateValue(variadicValues...); HandleError(p.eh, err) {
			p.parseErr = multierror.Append(p.parseErr, err)
		}
	}

	for name, vs := range p.keyValues {

		kvA := p.Key(name)

		if kvA == nil && p.Strict {
			if HandleError(p.eh, ErrUnrecognisedToken) {
				p.parseErr = multierror.Append(p.parseErr, ErrUnrecognisedToken)
			}
		} else if kvA == nil {
			continue
		}

		if err := kvA.updateValue(vs...); HandleError(p.eh, err) {
			p.parseErr = multierror.Append(p.parseErr, err)
		}

	}

	for name, vs := range p.flagValues {

		fA := p.Flag(name)

		if fA == nil && p.Strict {
			if HandleError(p.eh, ErrUnrecognisedToken) {
				p.parseErr = multierror.Append(p.parseErr, ErrUnrecognisedToken)
			}
		} else if fA == nil {
			continue
		}

		if err := fA.updateValue(vs...); HandleError(p.eh, err) {
			p.parseErr = multierror.Append(p.parseErr, err)
		}
	}

	missing := ErrMissing()
	for _, a := range p.Args() {

		// only want to parse the default value for the Arg if it's not a Pipe
		if !a.IsParsed() && a.Type() != PipeType {
			if err := a.updateValue(); HandleError(p.eh, err) {
				p.parseErr = multierror.Append(p.parseErr, err)
			}
		}

		if p.Strict {
			if !a.IsParsed() {
				missing.Add(a)
			}
		} else {
			if !a.IsParsed() && a.IsRequired() {
				missing.Add(a)
			}
		}
	}

	if len(missing.missing) > 0 && HandleError(p.eh, missing) {
		p.parseErr = multierror.Append(p.parseErr, missing)
	}
}

var (
	finished = errors.New("scan finished")
)

// Scan some input for argument tokens (i.e. a positional argument, a key-value argument value or a flag argument value).
//
// Returns the tokens that were 'consumed' (successful or not), the remaining uns-canned tokens in the input and any errors
// associated with scanning the 'consumed' tokens
func (p *Parser) scan(input []string) (remaining, consumed []string, err error) {

	if len(input) < 1 {
		return nil, nil, finished
	}

	// split the input into the next token, and the remaining tokens
	this, left := input[0], input[1:]

	// every time scan is called, we will *at least* consume this token
	consumed = []string{this}

	log.Debug("Scanning input", zap.String("this", this), zap.Strings("next", left))

	var argID, argValue string
	argType := Unrecognised
	flagSingleDash := false
	argValueDetected := false

	// detect the Type and argID from 'this' and sanitize (additionally, calculated argValue if possible)
	switch this[0] {
	case '-':
		// any '-' prefix indicates 'this' is a FlagType
		argType = FlagType
		if this[1] != '-' {
			flagSingleDash = true
		}
		this = strings.TrimLeft(this, "-")

		fallthrough
	default:
		if strings.Contains(this, "=") {

			// argType will be == FlagType if it was detected in the above block, if it is
			// still Unrecognised, then it must be a KeyValueType
			if argType == Unrecognised {
				argType = KeyValueType
			}

			keyValue := strings.SplitN(this, "=", 2)

			// key supplied with no value
			if keyValue[1] == "" {
				return left, consumed, ErrIncompleteToken
			}

			// both argID and argValue can be detected
			argID, argValue = keyValue[0], keyValue[1]
			argValueDetected = true

		} else if argType == FlagType {
			argID = this
		} else if argType == Unrecognised {
			argType = PositionType
			argValue = this
		}
	}

	if !(Unrecognised < argType && argType <= limit) {
		return left, consumed, ErrUnrecognisedType
	}

	defer func() {
		// At this point, we know what each of the argType is, so we can explicitly handle the case here and not have
		// to repeat it in the below block.
		//
		// We don't mind if we add argId and argValue that may be broken or result in errors - that will be thrown to
		// the Parser.parseErr.
		log.Debug("Processing scanned token",
			zap.String("arg_id", argID),
			zap.String("arg_value", argValue),
			zap.Stringer("arg_type", argType))

		switch argType {
		case PositionType:
			// for Positional argMap, we only care about this specific position and so do not consume any extra values
			// from the input
			p.positionalValues = append(p.positionalValues, argValue)
		case FlagType:
			if argValueDetected {
				p.flagValues[argID] = append(p.flagValues[argID], argValue)
			}
		case KeyValueType:
			p.keyValues[argID] = append(p.keyValues[argID], argValue)
		}
	}()

	switch argType {
	case PositionType:
		// for Positional argMap, we only care about this specific position and so do not consume any extra values
		// from the input
		argValue = this
		argID = fmt.Sprintf("%d", len(p.positionalValues)+1)

		return left, consumed, nil
	case KeyValueType:
		if !p.Has(Key(argID)) {
			return left, consumed, WarnUnrecognisedArg
		} else {
			return left, consumed, nil
		}
	case FlagType:
		// for FlagArg argMap we need to check if we are dealing with a BOOL flag - if we are, then we don't need to
		// consume any extra input, otherwise we may need to consume some extra input. We also need to check for
		// combined shorthands (e.g. -nWXyZ may need to be converted into -n -W -X -y -Z)

		fA := p.Flag(argID)

		if fA == nil {

			if flagSingleDash {
				// with a '-' prefix, argID can have 1 of 3 meanings:
				//
				//	1. full name of a flag
				//	2. the shorthand name of a flag
				//	3. the combination of the shorthands of multiple boolean flags
				//
				// We just need to check 3. as 1/2 will sort themselves out

				var newInput []string

				for _, c := range argID {

					newInput = append(newInput, fmt.Sprintf("-%c", c))
					if a, exists := p.shorthands[c]; !exists || a.Type() != FlagType {
						return left, []string{this}, fmt.Errorf("%w by shorthand: %c", ErrUnrecognisedToken, c)
					}
				}

				for _, i := range left {
					newInput = append(newInput, i)
				}

				return newInput, []string{this}, nil

			} else {
				// fA == nil => a flag with name argID does not exist
				return left, []string{this}, WarnUnrecognisedArg
			}
		}

		// flag was supplied like -k=<value> or --key=<value> => no more input to be consumed
		if argValueDetected {
			break
		}

		// we don't know the argValue yet, but it may be that the Flag is an indicator and doesn't require a value - if
		// this is the case, we can go and set some sensible defaults if they do not exist (e.g. false for bool flags).
		// Otherwise, we know we need to consume (or try to) the next token to get the value for the flag
		if fA.IsIndicator() {
			log.Debug("Handling INDICATOR flag", zap.String("_t", fmt.Sprintf("%T", fA)))
			switch indF := fA.(type) {
			case *FlagArg[bool]:

				log.Debug("Handling BOOL indicator")

				dflt, err := ValidateDefaultValue[bool](indF)

				if err != nil {
					// assume that when no default can be properly found, then default is false
					log.Warn("unable to validate *FlagArg[bool] default value", zap.Error(err))
					dflt = false
				}

				log.Debug("Setting BOOL flag value string", zap.Bool("value_string", !dflt))
				argValue = fmt.Sprintf("%t", !dflt)
				argValueDetected = true
			case *FlagArg[Counter]:
				// TODO: figure out why this just never gets hit :confused_potato:
				log.Debug("Handling COUNTER indicator")
			}
			argValueDetected = true
			return left, consumed, nil
		} else {
			if len(left) < 1 {
				return left, consumed, ErrIncompleteToken
			}
			argValue = left[0]
			consumed = append(consumed, argValue)
			left = left[1:]
			argValueDetected = true

			return left[1:], consumed, nil
		}
	}

	return left, consumed, ErrIncompleteToken
}

// Complete attaches arguments and flags to the completion Command for autocompletion support.
func (p *Parser) Complete(cmd *complete.Command) {

	var argPredictions []complete.Predictor

	if cmd.Args != nil {
		argPredictions = []complete.Predictor{
			cmd.Args,
		}
	}

	if len(p.Set.positions) > 0 {
		argPredictions = append(argPredictions, predict.Something)
	}

	argSet := make([]string, len(p.KeyValues()))

	for i, k := range p.KeyValues() {
		argSet[i] = fmt.Sprintf("%s=", k)
	}

	if len(p.KeyValues()) > 0 {
		argPredictions = append(argPredictions, predict.Set(argSet))
	}
	cmd.Args = predict.Or(argPredictions...)

	if cmd.Flags == nil {
		cmd.Flags = map[string]complete.Predictor{}
	}

	for _, fl := range p.Flags() {
		cmd.Flags[fl.Name()] = predict.Something
	}
}
