package clap

// NewMetadata is a constructor for a new Metadata.
func NewMetadata(opts ...Option) *Metadata {
	md := &Metadata{
		argShorthand: noShorthand,
	}

	for _, opt := range opts {
		opt.apply(md)
	}

	return md
}

// WithShorthand adds a shorthand/alias to some Arg.
//
// This is a no-op for IsPositional types.
func WithShorthand(sh rune) Option {
	return shorthandOpt{Shorthand: sh}
}

// WithUsage adds a usage string to some Arg.
func WithUsage(usage string) Option {
	return usageOpt{Usage: usage}
}

// WithDefault adds a default string value for some Arg.
//
// This is a no-op for IsPositional types.
func WithDefault(defaultValue string) Option {
	return defaultOpt{Default: defaultValue}
}

// AsRequired makes an Arg required.
func AsRequired() Option {
	return &requiredOpt{
		reqVal: true,
	}
}

func withDefaultDisabled() Option {
	return noDefaultOpt{}
}

func withNoShorthand() Option {
	return noShorthandOpt{}
}

// AsOptional makes an Arg optional.
func AsOptional() Option {
	return &requiredOpt{
		reqVal: false,
	}
}

// Metadata holds the metadata for some Arg.
//
// Metadata can only be updated through Option implementations.
type Metadata struct {
	argUsage     string // describes how to use the argument
	argShorthand rune   // a single character that can identify the argument
	argDefault   string // the default str to be used when the argument is not supplied
	hasDefault   bool   // indicates argDefault has been set (to differentiate between "" and not supplied)
	argRequired  bool   // indicates if the argument is mandatory
}

func (m *Metadata) Usage() string {
	return m.argUsage
}

func (m *Metadata) updateMetadata(opts ...Option) {
	for _, opt := range opts {
		opt.apply(m)
	}
}

func (m *Metadata) IsRequired() bool {
	return m.argRequired
}

func (m *Metadata) Shorthand() rune {
	return m.argShorthand
}

func (m *Metadata) Default() string {
	return m.argDefault
}

func (m *Metadata) HasDefault() bool {
	return m.hasDefault
}

// An Option describes a change to some *Metadata, such as changing the Key/setting the argument as required.
type Option interface {
	apply(*Metadata)
}

type noShorthandOpt struct{}

var (
	noShorthand rune = -255
)

func (_ noShorthandOpt) apply(metadata *Metadata) {
	metadata.argShorthand = noShorthand
}

type noDefaultOpt struct{}

func (_ noDefaultOpt) apply(metadata *Metadata) {
	metadata.argDefault = ""
	metadata.hasDefault = false
}

type usageOpt struct {
	Usage string
}

func (n usageOpt) apply(metadata *Metadata) {
	metadata.argUsage = n.Usage
}

type shorthandOpt struct {
	Shorthand rune
}

func (n shorthandOpt) apply(metadata *Metadata) {
	metadata.argShorthand = n.Shorthand
}

type defaultOpt struct {
	Default string
}

func (n defaultOpt) apply(metadata *Metadata) {
	metadata.argDefault = n.Default
	metadata.hasDefault = true
}

type requiredOpt struct {
	reqVal bool
}

func (r *requiredOpt) apply(metadata *Metadata) {
	metadata.argRequired = r.reqVal
}
