package clap

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrDuplicateID  = errors.New("identifier already exists")
	ErrPipeExists   = errors.New("pipe already exists")
	ErrInvalidIndex = errors.New("invalid positional index")
)

// NewSet is a constructor for a new empty *Set.
func NewSet() *Set {
	return &Set{
		shorthands: map[rune]argName{},
		k2n:        map[argName]string{},
		keys:       argMap[IsKeyValue]{},
		flags:      argMap[IsFlag]{},
		positions:  argMap[IsPositional]{},
		posArgs:    map[int]string{},
	}
}

// A Set is a container for a command-line Arg(s) of any Type.
type Set struct {
	shorthands map[rune]argName     // shorthands for KeyValueArg/FlagArgs, map: shorthand -> argName
	k2n        map[argName]string   // ids for all args, map: id -> Arg.Name
	keys       argMap[IsKeyValue]   // key-value args,	map: id -> KeyValueArg
	flags      argMap[IsFlag]       // flag args,        map: id -> FlagArg
	positions  argMap[IsPositional] // positional args,  map: id -> PositionalArg
	pipe       IsPipe               // the PipeArg (each Set can only have 1 PipeArg)
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
	_, exists := s.k2n[id.identify()]
	return exists
}

// Get returns an Arg with the given Identifier, if it exists in the Set.
func (s *Set) Get(id Identifier) Arg {
	if !s.Has(id) {
		return nil
	}
	an := id.identify()
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
func (s *Set) ByShorthand(shorthand rune) Arg {

	if shorthand == noShorthand {
		return nil
	}

	id, exists := s.shorthands[shorthand]

	if !exists {
		return nil
	}

	if id.Type() == PipeType || id.Type() == PositionType {
		return nil
	}

	return s.Get(id)
}

// Flag returns the flag for the given name/identifier, if it exists, otherwise nil.
func (s *Set) Flag(name string) IsFlag {

	if f, exists := s.flags[name]; exists {
		return f
	}

	if len(name) == 1 {
		a := s.ByShorthand(rune(name[0]))
		if fl, ok := a.(IsFlag); ok {
			return fl
		}
	}

	return nil
}

// Flags returns all flags(s) within the argMap.
func (s *Set) Flags() []IsFlag {
	return s.flags.List()
}

// AddFlag adds a flag argument to the Parser.
func (s *Set) AddFlag(f IsFlag, opts ...Option) error {

	f.updateMetadata(opts...)

	if s.Has(f) {
		return ErrDuplicateID
	}

	name := f.Name()
	k := FlagType.getIdentifier(name)

	if sh := f.Shorthand(); sh != noShorthand {
		if an, exists := s.shorthands[sh]; exists {
			return fmt.Errorf("shorthand %w for %q", ErrDuplicateID, an.Name())
		}
		s.shorthands[sh] = k
		s.k2n[FlagType.getIdentifier(fmt.Sprintf("%c", sh))] = name
	}

	s.flags[name] = f
	s.k2n[k] = name

	return nil
}

// Key returns the key-value argument for the given name/identifier, if it exists, otherwise nil.
func (s *Set) Key(name string) IsKeyValue {

	if f, exists := s.keys[name]; exists {
		return f
	}

	if len(name) == 1 {
		a := s.ByShorthand(rune(name[0]))
		if kv, ok := a.(IsKeyValue); ok {
			return kv
		}
	}
	return nil
}

// KeyValues returns all key-value arguments within the Set.
func (s *Set) KeyValues() []IsKeyValue {
	return s.keys.List()
}

// AddKeyValue adds a key-value argument to the Set.
func (s *Set) AddKeyValue(kv IsKeyValue, opts ...Option) error {

	kv.updateMetadata(opts...)

	if s.Has(kv) {
		return ErrDuplicateID
	}

	name := kv.Name()
	k := KeyValueType.getIdentifier(name)

	if sh := kv.Shorthand(); sh != noShorthand {

		if existingArgKey, shorthandExists := s.shorthands[sh]; shorthandExists {
			return fmt.Errorf("shorthand %w for %q", ErrDuplicateID, existingArgKey.Name())
		}
		s.k2n[KeyValueType.getIdentifier(fmt.Sprintf("%c", sh))] = name
		s.shorthands[sh] = k
	}

	s.keys[name] = kv
	s.k2n[k] = name

	return nil
}

// Pos returns the positional argument at the supplied index, if it exists, otherwise nil.
func (s *Set) Pos(index int) IsPositional {
	k, exists := s.posArgs[index]

	if !exists {
		return nil
	}

	return s.positions.Get(k)
}

// PosS is a wrapper around Pos which accepts a string representation of the integer position.
func (s *Set) PosS(sindex string) IsPositional {
	i, _ := strconv.ParseInt(sindex, 10, 64)

	return s.Pos(int(i))
}

// Positions returns all positional arguments within the Set.
func (s *Set) Positions() []IsPositional {
	return s.positions.List()
}

// AddPosition adds a positional argument to the Set.
//
// The Set expects positional arguments to be supplied in order.
func (s *Set) AddPosition(a IsPositional, opts ...Option) error {

	index := a.Index()
	expectedIndex := len(s.positions) + 1

	if s.Has(Position(a.Index())) {
		return ErrDuplicateID
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
func (s *Set) Pipe() IsPipe {
	return s.pipe
}

// AddPipe adds a pipe argument to the Set.
func (s *Set) AddPipe(p IsPipe, opts ...Option) error {
	if s.pipe == nil {
		p.updateMetadata(opts...)
		s.pipe = p
	} else {
		return ErrPipeExists
	}
	return nil
}

// argMap is a collection of Arg implementations (of the same Type), indexed by their name
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

// Get returns an Arg with the given name, if it exists within the collection.
func (s argMap[A]) Get(name string) A {
	return s[name]
}

//// A Collection of Arg(s) indexed by their respective Identifier.
////
//// TODO: finish this and wrap Set in this interface to make for easier testing
//type Collection interface {
//	// Args return all the Arg(s) within the Collection
//	Args() []Arg
//
//	// Has returns true if the Identifier can be found in the Collection, otherwise false
//	Has(Identifier) bool
//
//	// Get returns an Arg by its Identifier if it can be found in the Collection, otherwise nil
//	Get(Identifier) Arg
//
//	// Flag returns a flag argument by its name/shorthand if it exists, otherwise nil
//	Flag(string) IsFlag
//
//	// Flags returns all flags(s) within the argMap.
//	Flags() []IsFlag
//
//	// Key returns a key-value argument by its key/shorthand if it exists, otherwise nil
//	Key(string) IsKeyValue
//}
