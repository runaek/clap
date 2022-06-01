package clap

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// NewLinePipe is a constructor for a PipeArg which reads new lines from a pipe supplied via the command-line.
func NewLinePipe[T any](variable *T, parser ValueParser[T], options ...Option) *PipeArg[T] {
	return NewPipeArg[T](variable, parser, &SeparatedValuePiper{Separator: "\n"}, os.Stdin, options...)
}

// CSVPipe is a constructor for a PipeArg which reads comma-separated values from a pipe supplied via the command-line.
func CSVPipe[T any](variable *T, parser ValueParser[T], options ...Option) *PipeArg[T] {
	return NewPipeArg[T](variable, parser, &SeparatedValuePiper{Separator: ","}, os.Stdin, options...)
}

func NewPipeArg[T any](variable *T, parser ValueParser[T], piper Piper, input *os.File, options ...Option) *PipeArg[T] {
	options = append(options, pipeOptions...)

	return &PipeArg[T]{
		piper: piper,
		input: input,
		md:    NewMetadata(options...),
		v:     NewVariable[T](variable, parser),
	}
}

type IsPipe interface {
	Arg

	// ReadPipe decodes the contents of the pipe into string arguments
	ReadPipe(io.Reader) ([]string, error)

	mustEmbedPipe()
}

var pipeOptions = []Option{
	withNoShorthand(), withDefaultDisabled(), AsOptional(),
}

// Pipe is an Identifier for the PipeArg.
type Pipe string

const PipeName = "MAIN"

func (_ Pipe) identify() argName {
	return PipeType.getIdentifier(PipeName)
}

type PipeArg[T any] struct {
	piper Piper
	input *os.File

	md *Metadata
	v  Variable[T]

	parsed, supplied bool
}

func (p *PipeArg[T]) PipeInput() *os.File {
	if p.input == nil {
		p.input = os.Stdin
	}
	return p.input
}

func (p *PipeArg[T]) identify() argName {
	return PipeType.getIdentifier(p.Name())
}

func (p *PipeArg[T]) Name() string {
	return PipeName
}

func (p *PipeArg[T]) Type() Type {
	return PipeType
}

func (p *PipeArg[T]) Shorthand() rune {
	return -128
}

func (p *PipeArg[T]) Usage() string {
	return p.md.Usage()
}

func (p *PipeArg[T]) Default() string {
	return ""
}

func (p *PipeArg[T]) IsRepeatable() bool {
	return false
}

func (p *PipeArg[T]) IsRequired() bool {
	return false
}

func (p *PipeArg[T]) IsParsed() bool {
	return p.parsed
}

func (p *PipeArg[T]) IsSupplied() bool {

	if p.input == nil {
		return false
	}

	fi, err := p.input.Stat()

	if err != nil {
		return false
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		return true
	}

	return false
}

func (p *PipeArg[T]) Variable() Variable[T] {
	return p.v
}

func (p *PipeArg[T]) ReadPipe(in io.Reader) ([]string, error) {
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

	f := p.PipeInput()

	if f == nil {
		return nil
	}

	inputs, err := p.piper.Pipe(f)

	if err != nil {
		return err
	}

	return p.v.Update(inputs...)
}

func (p *PipeArg[T]) mustEmbedPipe() {}

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
	_ IsPipe        = &PipeArg[any]{}
	_ TypedArg[any] = &PipeArg[any]{}
)
