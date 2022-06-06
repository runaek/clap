package clap

import (
	"bytes"
	"errors"
	"github.com/runaek/clap/pkg/parse"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// NewLinePipe is a constructor for a PipeArg which reads new lines from a pipe supplied via the command-line.
func NewLinePipe[T any](variable *T, parser parse.Parser[T], options ...Option) *PipeArg[T] {
	return NewPipeArg[T](variable, parser, &SeparatedValuePiper{Separator: "\n"}, os.Stdin, options...)
}

// CSVPipe is a constructor for a PipeArg which reads comma-separated values from a pipe supplied via the command-line.
func CSVPipe[T any](variable *T, parser parse.Parser[T], options ...Option) *PipeArg[T] {
	return NewPipeArg[T](variable, parser, &SeparatedValuePiper{Separator: ","}, os.Stdin, options...)
}

func NewPipeArg[T any](variable *T, parser parse.Parser[T], piper Piper, input FileReader, options ...Option) *PipeArg[T] {
	options = append(options, pipeOptions...)
	return &PipeArg[T]{
		piper:    piper,
		input:    input,
		md:       NewMetadata(options...),
		v:        NewVariable[T](variable, parser),
		supplied: nil,
	}
}

func PipeUsingVariable[T any](piper Piper, input FileReader, v Variable[T], options ...Option) *PipeArg[T] {
	options = append(options, pipeOptions...)
	return &PipeArg[T]{
		piper:    piper,
		input:    input,
		md:       NewMetadata(options...),
		v:        v,
		supplied: nil,
	}
}

// IPipe is the interface satisfied by a PipeArg.
type IPipe interface {
	Arg

	// PipeInput returns the FileReader wrapping the underlying input for the pipe - this is usually os.input
	PipeInput() FileReader

	// PipeDecode decodes the contents of the pipe into string arguments to be parsed later
	PipeDecode(io.Reader) ([]string, error)

	updateInput(r FileReader)
}

var pipeOptions = []Option{
	withNoShorthand(), withDefaultDisabled(), AsOptional(),
}

// Pipe is an Identifier for the PipeArg.
//
// Any value of Pipe will Identify the singular PipeArg in a Parser.
type Pipe string

const PipeName = "MAIN"

func (_ Pipe) argName() argName {
	return PipeType.getIdentifier(PipeName)
}

// PipeArg represents *the* (there can only be a single PipeArg defined per Parser) command-line inpu provided from the
// Stdout of another program.
type PipeArg[T any] struct {
	piper Piper
	input FileReader
	md    *Metadata
	v     Variable[T]

	// data is the full read data from the pipe
	data     []byte
	parsed   bool
	supplied *bool
}

func (p *PipeArg[T]) updateInput(r FileReader) {
	if r == nil {
		return
	}
	p.input = r
}

func (p *PipeArg[T]) PipeInput() FileReader {
	if p.input == nil {
		p.input = os.Stdin
	}
	return p.input
}

func (p *PipeArg[T]) argName() argName {
	return PipeType.getIdentifier(p.Name())
}

func (p *PipeArg[T]) Name() string {
	return PipeName
}

func (p *PipeArg[T]) Type() Type {
	return PipeType
}

func (p *PipeArg[T]) Shorthand() rune {
	return noShorthand
}

func (p *PipeArg[T]) Usage() string {
	return p.md.Usage()
}

// Default returns an empty string - a pipe cannot have a default value.
func (p *PipeArg[T]) Default() string {
	return ""
}

// IsRepeatable returns false - a pipe can only be supplied once.
func (p *PipeArg[T]) IsRepeatable() bool {
	return false
}

// IsRequired returns false - a pipe is not required.
func (p *PipeArg[T]) IsRequired() bool {
	return false
}

func (p *PipeArg[T]) IsParsed() bool {
	return p.parsed
}

// IsSupplied checks if a pipe has been supplied & data has been written to the pipe.
func (p *PipeArg[T]) IsSupplied() (cond bool) {

	if p.supplied != nil {
		return *p.supplied
	}

	defer func() {
		if cond {
			t := true
			p.supplied = &t
		} else {
			f := false
			p.supplied = &f
		}
	}()

	if p.input == nil {
		return false
	}

	fi, err := p.input.Stat()

	if err != nil {
		return false
	}

	log.Debug("Program Input", zap.String("fn", fi.Name()), zap.Stringer("fm", fi.Mode()), zap.Int64("size", fi.Size()))

	if fi.Mode()&os.ModeNamedPipe != 0 && fi.Size() > 0 {
		return true
	}

	return false
}

func (p *PipeArg[T]) Variable() Variable[T] {
	return p.v
}

func (p *PipeArg[T]) PipeDecode(in io.Reader) ([]string, error) {
	return p.piper.Pipe(in)
}

func (p *PipeArg[T]) updateMetadata(options ...Option) {
	options = append(options, pipeOptions...)

	if p.md == nil {
		p.md = NewMetadata(options...)
		return
	}

	p.md.updateMetadata(options...)
}

func (p *PipeArg[T]) updateValue(_ ...string) error {

	if p.parsed {
		return nil
	}

	var fullData []byte

	if len(p.data) > 0 {
		fullData = p.data
	} else {
		f := p.PipeInput()

		if f == nil {
			return errors.New("pipe has no inpu")
		}
		b, err := ioutil.ReadAll(f)

		if err != nil {
			return errors.New("error reading from pipe")
		}
		fullData = b
		p.data = fullData
	}

	dat := bytes.NewReader(fullData)

	inputs, err := p.piper.Pipe(dat)

	if err != nil {
		return err
	}

	if err := p.v.Update(inputs...); err != nil {
		return err
	} else {
		p.parsed = true
	}

	return nil
}

// A Piper is responsible for decoding the data from a pipe into raw command-line arguments
type Piper interface {
	Pipe(in io.Reader) ([]string, error)
}

type SeparatedValuePiper struct {
	Separator string
}

func (s *SeparatedValuePiper) Pipe(in io.Reader) ([]string, error) {
	dat, err := ioutil.ReadAll(in)

	if err != nil {
		return nil, err
	}
	return strings.Split(string(dat), s.Separator), nil
}

var (
	_ Identifier    = Pipe("")
	_ IPipe         = &PipeArg[any]{}
	_ TypedArg[any] = &PipeArg[any]{}
)
