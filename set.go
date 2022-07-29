package clap

import (
	"fmt"
	"strconv"
)

// NewSet is a constructor for a new empty *Set.
func NewSet() *Set {
	return &Set{
		shorthands: map[string]argName{},
		k2n:        map[argName]string{},
		keys:       argMap[IKeyValue]{},
		flags:      argMap[IFlag]{},
		positions:  argMap[IPositional]{},
		posArgs:    map[int]string{},
	}
}

// A Set is a container for a command-line Arg(s) of any Type.
type Set struct {
	shorthands map[string]argName  // shorthands for KeyValueArg/FlagArgs, map: shorthand -> argName
	k2n        map[argName]string  // ids for all args, map: id -> Arg.Name
	keys       argMap[IKeyValue]   // key-value args,	map: id -> KeyValueArg
	flags      argMap[IFlag]       // flag args,        map: id -> FlagArg
	positions  argMap[IPositional] // positional args,  map: id -> PositionalArg
	pipe       IPipe               // the PipeArg (each Set can only have 1 PipeArg)
	posArgs    map[int]string
}

// Args returns all the Arg(s) within the Set.
func (s *Set) Args() []Arg {

	keys := s.keys.List()
	flags := s.flags.List()
	indexes := s.positions.List()

	pipes := 0

	if s.pipe != nil {
		pipes = 1
	}

	out := make([]Arg, len(keys)+len(flags)+len(indexes)+pipes)

	counter := 0

	for _, a := range keys {
		out[counter] = a
		counter++
	}

	for _, a := range flags {
		out[counter] = a
		counter++
	}

	for _, a := range indexes {
		out[counter] = a
		counter++
	}

	if s.pipe != nil {
		out[counter] = s.pipe
	}

	return out
}

// Has returns true if there is an Arg with the given Identifier in the Set.
func (s *Set) Has(id Identifier) bool {
	_, exists := s.k2n[id.argName()]

	return exists
}

// Get returns an Arg with the given Identifier, if it exists in the Set.
func (s *Set) Get(id Identifier) Arg {
	if !s.Has(id) {
		return nil
	}
	an := id.argName()
	name := an.Name()

	switch an.Type() {
	case FlagType:
		return s.Flag(name)
	case KeyValueType:
		return s.Key(name)
	case PositionType:
		return s.PosS(name)
	case PipeType:
		return s.Pipe()
	}
	return nil
}

// ByShorthand returns the Arg for the given shorthand identifier, if it exists, otherwise nil.
func (s *Set) ByShorthand(sh string) Arg {

	if sh == "" {
		return nil
	}

	id, exists := s.shorthands[sh]

	if !exists {
		return nil
	}

	if id.Type() == PipeType || id.Type() == PositionType {
		return nil
	}

	return s.Get(id)
}

// Flag returns the flag for the given Id/identifier, if it exists, otherwise nil.
func (s *Set) Flag(name string) IFlag { // nolint: ireturn

	if f, exists := s.flags[name]; exists {
		return f
	}

	if len(name) == 1 {
		a := s.ByShorthand(name)
		if fl, ok := a.(IFlag); ok {
			return fl
		}
	}

	return nil
}

// Flags returns all flags(s) within the argMap.
func (s *Set) Flags() []IFlag {
	return s.flags.List()
}

// AddFlag adds a flag argument to the Parser.
//
// Returns ErrDuplicate if the key or alias already exists in the Set.
func (s *Set) AddFlag(f IFlag, opts ...Option) error {

	f.updateMetadata(opts...)

	if s.Has(f) {
		return ErrDuplicate
	}

	name := f.Name()
	k := FlagType.getIdentifier(name)

	if sh := f.Shorthand(); sh != "" {
		if an, exists := s.shorthands[sh]; exists {
			return fmt.Errorf("%w: shorthand already in use by %q", ErrDuplicate, an.Name())
		}
		s.shorthands[sh] = k
		s.k2n[FlagType.getIdentifier(sh)] = name
	}

	s.flags[name] = f
	s.k2n[k] = name

	return nil
}

// Key returns the key-value argument for the given Id/identifier, if it exists, otherwise nil.
func (s *Set) Key(name string) IKeyValue { // nolint: ireturn

	if f, exists := s.keys[name]; exists {
		return f
	}

	if len(name) == 1 {
		a := s.ByShorthand(name)
		if kv, ok := a.(IKeyValue); ok {
			return kv
		}
	}
	return nil
}

// KeyValues returns all key-value arguments within the Set.
func (s *Set) KeyValues() []IKeyValue {
	return s.keys.List()
}

// AddKeyValue adds a key-value argument to the Set.
//
// Returns ErrDuplicate if the key or alias already exists in the Set.
func (s *Set) AddKeyValue(kv IKeyValue, opts ...Option) error {

	kv.updateMetadata(opts...)

	if s.Has(kv) {
		return ErrDuplicate
	}

	name := kv.Name()
	k := KeyValueType.getIdentifier(name)

	if sh := kv.Shorthand(); sh != "" {

		if existingArgKey, shorthandExists := s.shorthands[sh]; shorthandExists {
			return fmt.Errorf("%w: shorthand already in use by %q", ErrDuplicate, existingArgKey.Name())
		}
		s.k2n[KeyValueType.getIdentifier(sh)] = name
		s.shorthands[sh] = k
	}

	s.keys[name] = kv
	s.k2n[k] = name

	return nil
}

// Pos returns the positional argument at the supplied index, if it exists, otherwise nil.
func (s *Set) Pos(index int) IPositional { // nolint: ireturn
	k, exists := s.posArgs[index]

	if !exists {
		return nil
	}

	return s.positions.Get(k)
}

// PosS is a wrapper around Pos which accepts a string representation of the integer position.
func (s *Set) PosS(sindex string) IPositional {
	i, _ := strconv.ParseInt(sindex, 10, 64)

	return s.Pos(int(i))
}

// Positions returns all positional arguments within the Set.
func (s *Set) Positions() []IPositional {
	return s.positions.List()
}

// AddPosition adds a positional argument to the Set.
//
// The Set expects positional arguments to be supplied in order and returns a wrapped ErrInvalidIndex if arguments
// are specified in an invalid order.
func (s *Set) AddPosition(a IPositional, opts ...Option) error {

	index := a.Index()
	expectedIndex := len(s.positions) + 1

	if s.Has(Position(a.Index())) {
		return ErrDuplicate
	}

	lastArg := s.Pos(len(s.positions))

	if lastArg != nil && !lastArg.IsRequired() && a.IsRequired() {
		return fmt.Errorf("%w: position %d cannot be required because position %d is optional",
			ErrInvalidIndex, index, len(s.positions))
	}

	a.updateMetadata(opts...)
	name := a.Name()
	k := PositionType.getIdentifier(name)

	if index != expectedIndex {
		return fmt.Errorf("%w: got %d but expected %d (positional arguments must be added in order)",
			ErrInvalidIndex,
			index,
			expectedIndex)
	} else if expectedIndex > 1 {
		lastPos := s.Pos(expectedIndex - 1)
		if lastPos.IsRepeatable() {
			return fmt.Errorf("%w: cannot add positional argument %d as it comes after a varying argument",
				ErrInvalidIndex,
				index)
		}
	}

	s.positions[name] = a
	s.posArgs[index] = name
	s.k2n[k] = name

	return nil
}

// Pipe returns the pipe argument for the Set, if it exists, otherwise nil.
func (s *Set) Pipe() IPipe { // nolint: ireturn
	return s.pipe
}

// AddPipe adds a pipe argument to the Set.
//
// Returns ErrPipeExists if a pipe already exists in the Set.
func (s *Set) AddPipe(p IPipe, opts ...Option) error {
	if s.pipe == nil {
		p.updateMetadata(opts...)
		s.pipe = p
		s.k2n[Pipe("").argName()] = p.Name()
	} else {
		return ErrPipeExists
	}
	return nil
}

// argMap is a collection of Arg implementations (of the same Type), indexed by their Id
type argMap[A Arg] map[string]A

// List returns all the Arg(s) within the collection.
func (s argMap[A]) List() []A {
	out := make([]A, len(s))

	i := 0

	for _, a := range s {
		out[i] = a
		i++
	}

	return out
}

// Get returns an Arg with the given Id, if it exists within the collection.
func (s argMap[A]) Get(name string) A {
	return s[name]
}
