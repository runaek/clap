package clap

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
	"github.com/runaek/clap/pkg/parse"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"text/template"
)

type ErrorHandling int

const (
	ContinueOnError = ErrorHandling(flag.ContinueOnError)
	ExitOnError     = ErrorHandling(flag.ExitOnError)
	PanicOnError    = ErrorHandling(flag.PanicOnError)
)

// handleError handles some error according to the ErrorHandling, returns false if e == nil, otherwise what/whether the
// function returns is determined by the ErrHandling mode:
//
//	ContinueOnError => return true
//	PanicOnError    => panic with error message
//  ExitOnError     => write error message to w and exit 1
func handleError(w io.Writer, eh ErrorHandling, e error) bool {
	if e == nil {
		return false
	}

	me, isMultiErr := e.(*multierror.Error)

	if isMultiErr {
		e = me.ErrorOrNil()

		if e == nil {
			return false
		}
	}

	switch eh {
	case ContinueOnError:
		if e == ErrHelp {
			return false
		}
	case PanicOnError:
		if e != ErrHelp {
			panic(e)
		}
	case ExitOnError:
		if e != ErrHelp {
			_, _ = fmt.Fprint(w, e)
			os.Exit(1)
		}
		os.Exit(0)
	default:
		panic(fmt.Errorf("invalid error handling detected whilst handling error: %w", e))
	}

	return true
}

// New creates a new command-line argument Parser with ErrorHandling mode ContinueOnError.
//
// `elements` can be either concrete Arg implementations, or structs with Arg(s) defined via struct-tags, which
// will be derived at runtime.
func New(name string, elements ...any) (*Parser, error) {
	p := NewParser(name, ContinueOnError).Using(elements...)

	return p, p.Valid()
}

// Must
func Must(name string, elements ...any) *Parser {
	p, err := NewAt(name, ContinueOnError, elements...)

	if err != nil {
		panic(fmt.Errorf("unable to construct Parser: %w", err))
	}

	return p
}

// NewAt is a constructor for a Parser at a specific ErrorHandling level.
//
// If no elements are supplied, NewAt is guaranteed to return a nil error.
func NewAt(name string, errHandling ErrorHandling, elements ...any) (*Parser, error) {
	s := NewParser(name, errHandling)

	if len(elements) != 0 {
		args, err := DeriveAll(elements...)

		if err != nil {
			return s, err
		}

		s.Add(args...)
	}

	return s, s.Valid()
}

// NewParser is a constructor for a new command-line argument Parser.
func NewParser(n string, errHandling ErrorHandling) *Parser {
	validationErr := new(multierror.Error)
	validationErr.ErrorFormat = func(es []error) string {
		hdr := fmt.Sprintf("invalid parser state [%d error(s) occurred]:\n", len(es))

		msgs := make([]string, len(es))

		for i, e := range es {
			msgs[i] = fmt.Sprintf("\t* %s", e)
		}

		return hdr + strings.Join(msgs, "\n")
	}

	return &Parser{
		Id:            n,
		Set:           NewSetWithHelp(),
		ErrorHandling: errHandling,
		Stdout:        os.Stdout,
		Stdin:         os.Stdin,
		Stderr:        os.Stderr,
		flagValues:    map[string][]string{},
		keyValues:     map[string][]string{},
		argState:      map[argName]error{},
		vErr:          validationErr,
	}
}

// A Parser contains and parses a number of flag, key-value or positional inputs from the command-line.
//
// Once all desired Arg have been added to the Parser (see Add), the arguments can be parsed using the Parse method.
//
// By default, a Parser will read from os.Stdin and write to os.Stdout and os.Stderr - these fields can be changed as
// required to any FileReader or FileWriter.
type Parser struct {
	// Id for the Parser
	Id string

	// Name of the program
	Name string

	// Description of the program
	Description string

	// ErrorHandling mode to be used by the Parser
	ErrorHandling ErrorHandling

	// underlying Set for the Parser - holding all the Arg types the Parser is responsible for
	*Set

	// Shift shifts the Parser along, ignoring the first 'shifted' arguments
	Shift int

	// Strict defines whether unrecognised tokens (e.g. flag/keys) are ignored, or return ErrUnidentified
	Strict bool

	// SuppressUsage stops the Usage from being written out when an error occurs
	SuppressUsage bool

	// SuppressValidation stops validation errors (i.e. adding arguments to the Parser) from breaking
	// the program
	SuppressValidation bool

	Stdin  FileReader
	Stdout FileWriter
	Stderr FileWriter

	positionalValues []string            // raw positional arguments from the latest Parse call
	keyValues        map[string][]string // raw key-value arguments from the latest Parse call
	flagValues       map[string][]string // raw flag-value arguments from the latest Parse call

	argState map[argName]error // contains errors for each argName during Parse
	pErr     error             // pErr is the error for the *latest* call to Parse
	vErr     *multierror.Error // vErr are validation errors caught whilst adding arguments to the Parser (Set)
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
		case IFlag:
			_ = p.AddFlag(arg)
		case IKeyValue:
			_ = p.AddKeyValue(arg)
		case IPositional:
			_ = p.AddPosition(arg)
		case IPipe:
			_ = p.AddPipe(arg)
		default:
			log.Warn("Unable to add malformed Arg",
				zap.String("_t", fmt.Sprintf("%T", a)),
				zap.String("name", a.Name()),
				zap.Stringer("type", a.Type()))
		}
	}

	return p
}

// State checks the 'state' of some Arg by its Identifier *after* Parse has been called (always returning nil before).
func (p *Parser) State(id Identifier) error {
	err, exists := p.argState[id.argName()]
	if !exists || err == ok { // nolint: errorlint
		return nil
	}

	return err
}

// AddPosition adds a positional argument to the Parser.
func (p *Parser) AddPosition(a IPositional, opts ...Option) *Parser {
	if err := p.Set.AddPosition(a, opts...); err != nil {
		p.vErr = multierror.Append(p.vErr, fmt.Errorf("%w: unable to add positional argument", err))
	}
	return p
}

// AddFlag adds a flag argument to the Parser.
func (p *Parser) AddFlag(f IFlag, opts ...Option) *Parser {
	if err := p.Set.AddFlag(f, opts...); err != nil {
		p.vErr = multierror.Append(p.vErr, fmt.Errorf("%w: unable to add flag", err))
	}
	return p
}

// AddKeyValue adds a key-value argument to the Parser.
func (p *Parser) AddKeyValue(kv IKeyValue, opts ...Option) *Parser {
	if err := p.Set.AddKeyValue(kv, opts...); err != nil {
		p.vErr = multierror.Append(p.vErr, fmt.Errorf("%w: unable to add key-value", err))
	}
	return p
}

// AddPipe adds a pipe argument to the Parser.
func (p *Parser) AddPipe(pipe IPipe, options ...Option) *Parser {
	if err := p.Set.AddPipe(pipe, options...); err != nil {
		p.vErr = multierror.Append(p.vErr, fmt.Errorf("%w: unable to add pipe", err))

		return p
	}
	pipe.updateInput(p.Stdin)
	return p
}

// WithDescription sets a description for the Parser.
func (p *Parser) WithDescription(desc string) *Parser {
	p.Description = desc
	return p
}

// Parse command-line input.
//
// Parse does not return an error, instead, during the run any errors that occur are collected and stored internally
// to be retrieved via the Err method. If the ErrorHandling is not set to ContinueOnError, then errors during Parse
// will cause the program to either panic of exit.
func (p *Parser) Parse(argv ...string) {
	log.Debug("Parsing input",
		zap.String("parser", p.Id),
		zap.Int("handling", int(p.ErrorHandling)),
		zap.Bool("strict", p.Strict),
		zap.Int("shift", p.Shift),
		zap.Strings("input", argv))

	if len(argv) == 0 {
		log.Debug("Using os.Args", zap.Strings("input", os.Args))
		argv = os.Args[1+p.Shift:]
	} else if len(argv) <= p.Shift && len(argv) > 0 {
		p.pErr = fmt.Errorf("unable to parse, invalid number of arguments: got %d but want at least %d", len(argv), p.Shift)

		return
	}

	p.argState = map[argName]error{}
	p.positionalValues = nil

	parseErr := new(multierror.Error)
	parseErr.ErrorFormat = func(es []error) string {
		hdr := fmt.Sprintf("parser failure: %d error(s) occurred parsing %q:\n", len(es), strings.Join(argv, " "))

		msgs := make([]string, len(es))

		for i, e := range es {
			msgs[i] = fmt.Sprintf("\t* %s", e)
		}

		return hdr + strings.Join(msgs, "\n")
	}

	p.pErr = parseErr

	for tkns, consumed, err := argv[p.Shift:], []string{}, error(nil); err != finished; tkns, consumed, err = p.scan(tkns) {
		if err == nil {
			continue
		}

		if errors.Is(err, ErrHelp) {
			log.Debug("Help requested")
			p.pErr = multierror.Append(p.pErr, ErrHelp)
			continue
		}

		if errors.Is(err, ErrUnidentified) && p.Strict {
			scanErr := ErrScanning(err, consumed...)
			log.Warn("Error during token scan", zap.Error(scanErr))
			p.pErr = multierror.Append(p.pErr, scanErr)

			continue
		} else if errors.Is(err, ErrUnidentified) {
			log.Warn("Unidentified argument error suppressed", zap.Error(err))

			continue
		}

		if err != nil {
			scanErr := ErrScanning(err, consumed...)
			log.Warn("Error during token scan", zap.Error(scanErr))
			p.pErr = multierror.Append(p.pErr, scanErr)
		}

		log.Debug("Scanned", zap.Strings("tkns", tkns), zap.Strings("consumed", consumed), zap.Error(err))
	}

	p.parse()

	if errors.Is(p.Err(), ErrHelp) {
		return
	}
	handleError(p.Stderr, p.ErrorHandling, p.Err())
}

// Ok is a helper function that panics if there were any Parser errors. This is useful for running at the end of a chain
// of method calls for the Parser.
//
// Returns the Parser for convenience.

// Ok checks the state of the Parser (both parse and validation errors, unless SuppressValidation is set)
//
// Returns the *Parser for convenience.
func (p *Parser) Ok() *Parser {
	parserErr := p.pErr
	validationErr := p.vErr

	if validationErr.Len() > 0 && !p.SuppressValidation {
		parserErr = multierror.Append(parserErr, validationErr)
	}

	if errors.Is(parserErr, ErrHelp) {
		p.Usage()
		os.Exit(0)
	}

	if handleError(p.Stderr, p.ErrorHandling, parserErr) {
		_, _ = fmt.Fprintln(p.Stdout, parserErr)

		if !p.SuppressUsage {
			p.Usage()
		}
		os.Exit(1)
	}

	return p
}

// Valid returns validation error(s) that occurred trying to add arguments to the Parser.
func (p *Parser) Valid() error {
	if p.vErr == nil {
		return nil
	}

	return p.vErr.ErrorOrNil()
}

// Err returns the error(s) for the *latest* call to Parse.
func (p *Parser) Err() error {
	if p.pErr == nil {
		return nil
	}

	switch err := p.pErr.(type) {
	// make sure we don't return an empty non-nil multi error
	case *multierror.Error:
		if len(err.Errors) == 0 {
			return nil
		}
		return err
	default:
		return err
	}
}

const usageTemplate = `Usage: {{ .Name }} [ <args>, ... ] [ <key>=<value>, ... ] [ --<flag>=<value>, ... ]

{{- if .Description }}

{{ .Description }} 
{{- end -}}

{{- if .Pipe }}

PIPE: {{ .Pipe }}
{{- end }}

{{- if .Arguments }}

ARGUMENTS
{{- range .Arguments }}
	{{ printf "%s" . -}}
{{- end -}}
{{- end }}

{{- if .Keys }}

OPTIONS (<key>=<value>)
{{- range $name, $desc := .Keys }}
	{{ printf "%-24s : %s" $name $desc }} 
{{- end -}}
{{- end }}

{{- if .Flags }}

FLAGS
{{- range $name, $desc := .Flags }}
	{{ printf "%-24s : %s" $name $desc }} 
{{- end -}}
{{- end }}
`

type usageTemplateData struct {
	Name        string
	Description string
	Arguments   []string
	Keys        map[string]string
	Flags       map[string]string
	Pipe        string
}

// Usage writes the help-text/usage to the output.
func (p *Parser) Usage() {
	dat := usageTemplateData{
		Arguments: []string{},
		Keys:      map[string]string{},
		Flags:     map[string]string{},
	}

	if p.Name == "" {
		dat.Name = p.Id
	} else if p.Id == "SYSTEM" {
		dat.Name = "<program_name>"
	} else {
		dat.Name = p.Name
	}

	if pA := p.Pipe(); pA != nil {
		dat.Pipe = pA.Usage()
	}

	dat.Description = p.Description

	positionalArgs := make([]string, len(p.Positions()))

	for _, pa := range p.Positions() {
		k := fmt.Sprintf("%d", pa.Index())

		usage := pa.Usage()

		if pa.IsRepeatable() {
			k += " ..."
			usage = fmt.Sprintf("(repeatable) %s", usage)
		}

		positionalArgs[pa.Index()-1] = fmt.Sprintf("%-6s: %s", k, usage)
	}

	dat.Arguments = positionalArgs

	for _, a := range p.Args() {
		switch a.(type) {
		case IFlag:

			n := "--" + a.Name()

			var sh string

			if a.Shorthand() != "" {
				sh = fmt.Sprintf("[-%s]", a.Shorthand())
			}
			n = fmt.Sprintf("%-5s %s", sh, n)
			dat.Flags[n] = a.Usage()

		case IKeyValue:

			n := a.Name()
			var sh string

			if a.Shorthand() != "" {
				sh = fmt.Sprintf("[%s]", a.Shorthand())
			}

			n = fmt.Sprintf("%-4s %s", sh, n)
			dat.Keys[n] = a.Usage()

		default:
			continue
		}
	}

	tpl := template.New("help")

	if t, err := tpl.Parse(usageTemplate); err != nil {
		panic(fmt.Errorf("unable to write help to output: %w", err))
	} else if terr := t.Execute(p.Stdout, dat); err != nil {
		panic(fmt.Errorf("error executing help template: (%s) %w", err, terr))
	}
}

var (
	finished = errors.New("scan finished")
	ok       = errors.New("ok")
)

// scan some input for argument tokens (i.e. a positional argument, a key-value argument value or a
// flag argument value).
//
// Returns the tokens that were 'consumed' (successful or not), the remaining un-scanned tokens in
// the input and any errors associated with scanning the 'consumed' tokens
//
// NOTE: will *not* recursively call itself, it will scan a single token (1 or 2 elements) and
// return - it is on the caller to make repeated calls to scan to consume the entire input, which
// is when scan will return finished.
func (p *Parser) scan(input []string) (remaining, consumed []string, err error) {
	if len(input) < 1 {
		return nil, nil, finished
	}

	// split the input into the next token, and the remaining tokens
	token, left := input[0], input[1:]

	this := token
	// every time scan is called, we will *at least* consume this token
	consumed = []string{this}

	log.Debug("Scanning input", zap.String("this", this), zap.Strings("next", left))

	var argID, argValue string
	argType := Unrecognised
	flagSingleDash := false
	argValueDetected := false

	if this == "-h" || this == "--help" {
		return left, consumed, ErrHelp
	}
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
				return left, consumed, ErrInvalid
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
		return left, consumed, ErrUnknownType
	}

	defer func() {
		// At this point, we know what each of the argType is, so we can explicitly handle the case here and not have
		// to repeat it in the below block.

		log.Debug("Scanned tokens",
			zap.String("arg_id", argID),
			zap.String("arg_value", argValue),
			zap.Stringer("arg_type", argType),
			zap.Error(err))

		switch argType {
		case PositionType:
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

		if p.Has(Key(argID)) {
			return left, consumed, nil
		} else {
			return left, consumed, fmt.Errorf("%w: no such Key", ErrUnidentified)
		}
	case FlagType:
		// for FlagArg argMap we need to check if we are dealing with a BOOL flag - if we are, then we don't need to
		// consume any extra input, otherwise we may need to consume some extra input. We also need to check for
		// combined shorthands (e.g. -nWXyZ may need to be converted into -n -W -X -y -Z)

		fA := p.Flag(argID)

		if fA == nil {
			if flagSingleDash {
				if len(argID) == 1 {
					if an, exists := p.shorthands[argID]; exists && an.Type() == FlagType {
						return left, consumed, nil
					} else {
						return left, []string{token}, fmt.Errorf("%w: no such Flag", ErrUnidentified)
					}

					// if argValueDetected => not combined boolean flags
				} else if !argValueDetected {
					var newInput []string

					for _, c := range argID {
						newInput = append(newInput, fmt.Sprintf("-%c", c))

						if a, exists := p.shorthands[string(c)]; !exists || a.Type() != FlagType {
							return left, []string{token}, fmt.Errorf("%w: no such Flag", ErrUnidentified)
						}
					}

					newInput = append(newInput, left...)

					argValueDetected = true

					return newInput, []string{this}, nil
				} else {
					return left, []string{token}, fmt.Errorf("%w: no such Flag", ErrUnidentified)
				}
			} else {
				// fA == nil => a flag with Id argID does not exist
				return left, []string{token}, fmt.Errorf("%w: no such Flag", ErrUnidentified)
			}
		}

		// flag was supplied like -k=<value> or --key=<value> => no more input to be consumed
		if argValueDetected || fA == nil {
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
			case *FlagArg[parse.C]:
				// TODO: figure out why this just never gets hit when C is in this package :confused_potato:
				log.Debug("Handling COUNTER indicator")
			}
			argValueDetected = true
			return left, consumed, nil
		} else {
			// flag needs a value
			if len(left) == 0 {
				return left, consumed, ErrInvalid
			}
			argValue = left[0]
			consumed = append(consumed, argValue)
			argValueDetected = true

			// argValue might have been the last value
			if len(left) > 1 {
				return left[1:], consumed, nil
			} else {
				return left, consumed, finished
			}
		}
	}

	return left, consumed, nil
}

// parse the Arg(s) within the Parser.
//
// NOTE: assumes that input has already been scanned
func (p *Parser) parse() {
	if pipe := p.Pipe(); pipe != nil && pipe.IsSupplied() {
		if err := pipe.updateValue(); err != nil {
			p.updateState(pipe, err)

			if p.Strict {
				p.pErr = multierror.Append(p.pErr, ErrParsing(Pipe(""), err))
			}
		}
	}

	variadicIndex := -1

	var variadicValues []string

	for im1, v := range p.positionalValues {
		pA := p.Pos(im1 + 1)

		// we have reached the end of the positional arguments, or we are processing
		// a variadic argument that started before this
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
			if err := pA.updateValue(v); err != nil {
				p.updateState(pA, err)

				if pA.IsRequired() || p.Strict {
					p.pErr = multierror.Append(p.pErr, err)
				}
			}
		}
	}

	if variadicIndex > 0 {
		variadicPosArg := p.Pos(variadicIndex)
		log.Debug("Updating variadic positional arguments", zap.Strings("vals", variadicValues))

		if err := variadicPosArg.updateValue(variadicValues...); err != nil {
			p.updateState(variadicPosArg, err)

			if variadicPosArg.IsRequired() || p.Strict {
				p.pErr = multierror.Append(p.pErr, err)
			}
		}
	}

	for name, vs := range p.keyValues {
		kvA := p.Key(name)

		// no need to log an error since this would have been noticed during the scan
		if kvA == nil {
			log.Warn("No such key-value argument to parse", zap.String("name", name), zap.Strings("values", vs))

			continue
		}

		if err := kvA.updateValue(vs...); err != nil {
			// we always want to update the state to track every Arg that fails, but we only want to add an error if we
			// are running in strict mode (which prohibits any parser failures) or if the arg is required
			p.updateState(kvA, err)
			if kvA.IsRequired() || p.Strict {
				p.pErr = multierror.Append(p.pErr, err)
			}
		}
	}

	for name, vs := range p.flagValues {
		fA := p.Flag(name)

		// no need to log an error since this would have been noticed during the scan
		if fA == nil {
			log.Warn("No such flag argument to parse", zap.String("name", name), zap.Strings("values", vs))

			continue
		}

		if err := fA.updateValue(vs...); err != nil {
			p.updateState(fA, err)

			if fA.IsRequired() || p.Strict {
				p.pErr = multierror.Append(p.pErr, err)
			}
		}
	}

	for _, a := range p.Args() {
		log.Debug("Checking Argument", zap.String("name", a.Name()), zap.Stringer("type", a.Type()), zap.String("default", a.Default()))

		if a.IsParsed() {
			p.updateState(a, ok)

			continue
		}

		// error parsing has already been reported
		if err := p.State(a); err != nil {
			continue
		}

		shouldAttemptDefaultParse := true

		switch arg := a.(type) {
		case IFlag:
			if !arg.HasDefault() {
				shouldAttemptDefaultParse = false
			}
		case IKeyValue:
			if !arg.HasDefault() {
				shouldAttemptDefaultParse = false
			}
		case IPipe, IPositional:
			// default positional/pipe values do not make sense
			shouldAttemptDefaultParse = false
		}

		if shouldAttemptDefaultParse {
			if err := a.updateValue(); err != nil {
				p.updateState(a, ErrMissing)

				if p.Strict || a.IsRequired() {
					p.pErr = multierror.Append(p.pErr, fmt.Errorf("error parsing default value: %w", err))
				}
			} else {
				p.updateState(a, ok)
			}
			continue
		}

		// if strict-mode: add the error to the final pErr
		// otherwise, we just need to update the state of the Arg
		if p.Strict || a.IsRequired() {
			p.pErr = multierror.Append(p.pErr, fmt.Errorf("%w: %s (%s)", ErrMissing, a.Name(), a.Type()))
		}

		p.updateState(a, ErrMissing)
	}
}

func (p *Parser) updateState(a Arg, state error) {
	p.argState[a.argName()] = state
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
		if fl.IsIndicator() {
			cmd.Flags[fl.Name()] = predict.Nothing

			continue
		}
		cmd.Flags[fl.Name()] = predict.Something
	}
}

func (p *Parser) Using(elements ...any) *Parser {
	return p
}
