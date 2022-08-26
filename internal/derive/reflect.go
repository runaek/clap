package derive

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"reflect"
	"strconv"
	"strings"
)

const (
	sigRequired     = "!"
	nameDescription = "Description"
	nameUsage       = "Usage"
	nameDefault     = "Default"

	FlagType = "-"
	KeyType  = "@"
	PosType  = "#"
)

// TODO: improve error handling/reporting

// Parse a struct into a *Program.
func Parse(input any) (*Program, error) {
	p := &Program{}

	return p, p.Decode(input)
}

// A Program represents a structure with a number of Argument that can be automatically derived and parsed
// into clap.Arg in order to be parsed.
type Program struct {
	// arguments for the Program
	arguments map[string]*Argument

	order        []string
	descriptions map[string]string
	defaults     map[string]string
}

// Args returns the *Argument in the order they are defined in the struct.
func (p *Program) Args() []*Argument {
	out := make([]*Argument, len(p.order))

	for i, name := range p.order {
		a := p.arguments[name]

		if d, ok := p.descriptions[name]; ok {
			a.Usage = d
		}

		if d, ok := p.defaults[name]; ok {
			a.Default = d
		}

		out[i] = a
	}

	return out
}

// Decode some input type and attempt to derive a number of *Argument from struct-tags and fields.
func (p *Program) Decode(input any) error {
	if p.arguments == nil {
		p.arguments = map[string]*Argument{}
	}
	if p.descriptions == nil {
		p.descriptions = map[string]string{}
	}
	if p.defaults == nil {
		p.defaults = map[string]string{}
	}

	var strukt reflect.Value

	rawVal := reflect.ValueOf(input)

	if rawVal.Kind() == reflect.Pointer {
		strukt = rawVal.Elem()
	} else {
		strukt = rawVal
	}

	if strukt.Kind() != reflect.Struct {
		return fmt.Errorf("unable to Parse struct-tags from non-struct type %T", input)
	}

	result := &multierror.Error{
		ErrorFormat: func(es []error) string {
			msgs := make([]string, len(es)+1)
			msgs[0] = fmt.Sprintf("%d error(s) occurred deriving arguments:", len(es))

			for i := 1; i <= len(es); i++ {
				msgs[i] = fmt.Sprintf("\t> %s", es[i-1].Error())
			}

			return strings.Join(msgs, "\n")
		},
	}

	for i := 0; i < strukt.NumField(); i++ {
		fieldType := strukt.Type().Field(i)

		if !fieldType.IsExported() {
			continue
		}

		fieldValue := strukt.Field(i)

		arg, err := p.parseArgument(fieldType, fieldValue)

		if err != nil {
			result = multierror.Append(result, err)

			continue
		} else if arg == nil {
			continue
		}

		p.arguments[arg.Name] = arg
		p.order = append(p.order, arg.Name)
	}

	for k, v := range p.descriptions {
		if a, ok := p.arguments[k]; ok {
			a.Usage = v
		}
	}

	for k, v := range p.defaults {
		if a, ok := p.arguments[k]; ok {
			a.Default = v
		}
	}

	return result.ErrorOrNil()
}

// isReference is a helper function for checking whether a struct-field is referring to another field (via the same
// name - its prefix).
//
// It returns the cleaned name, and true if it is a reference, otherwise false and name is returned
func isReference(name string) (string, bool) {
	isRef := true
	if strings.HasSuffix(name, nameUsage) {
		name = strings.TrimSuffix(name, nameUsage)
	} else if strings.HasSuffix(name, nameDescription) {
		name = strings.TrimSuffix(name, nameDescription)
	} else if strings.HasSuffix(name, nameDefault) {
		name = strings.TrimSuffix(name, nameDefault)
	} else {
		isRef = false
	}
	return name, isRef
}

// parseArgument
func (p *Program) parseArgument(sf reflect.StructField, val reflect.Value) (*Argument, error) {
	name := sf.Name

	if cleaned, isRef := isReference(name); isRef {
		name = cleaned
	}

	arg, argExists := p.arguments[name]

	t := sf.Tag.Get("cli")
	switch strings.TrimPrefix(sf.Name, name) {
	case nameUsage, nameDescription:

		s, ok := val.Interface().(string)

		if !ok {
			return nil, fmt.Errorf("invalid 'Usage' (%s) type: want string but got %T", sf.Name, val.Interface())
		}

		p.descriptions[name] = s
		return nil, nil
	case nameDefault:
		s, ok := val.Interface().(string)

		if !ok {
			return nil, fmt.Errorf("invalid 'Default' (%s) type: want string but got %T", sf.Name, val.Interface())
		}

		p.defaults[name] = s
		return nil, nil
	default:

		if t == "" {
			return nil, nil
		}

		if !argExists {
			if newArg, err := newArgument(sf); err != nil {
				return nil, err
			} else {
				arg = newArg
			}
		}

		switch val.Kind() {
		case reflect.Pointer:

			elem := val.Elem()

			arg.ptr = elem.Interface()
		default:
			if val.CanAddr() {
				arg.ptr = val.Addr().Interface()
			} else {
				return nil, fmt.Errorf("unable to 'Address' Field Value (%s)", sf.Name)
			}
		}
	}
	return arg, nil
}

// newArgument attempts to construct an *Argument for some struct-field.
//
// It returns nil if no *Argument is defined by the struct-field.
func newArgument(sf reflect.StructField) (*Argument, error) {
	if !sf.IsExported() {
		return nil, nil
	}

	name := sf.Name

	tag := sf.Tag.Get("cli")

	cleanedName, isRef := isReference(sf.Name)

	if isRef {
		name = cleanedName
	} else if tag == "" {
		return nil, nil
	}

	if strings.HasSuffix(sf.Name, nameUsage) {
		name = strings.TrimSuffix(sf.Name, nameUsage)
	} else if strings.HasSuffix(sf.Name, nameDescription) {
		name = strings.TrimSuffix(sf.Name, nameDescription)
	} else if strings.HasSuffix(sf.Name, nameDefault) {
		name = strings.TrimSuffix(sf.Name, nameDefault)
	}

	out := &Argument{
		Name: name,
	}

	if sf.Name != name {
		return out, nil
	}

	switch tag[:1] {
	case sigRequired:
		out.Required = true
		tag = tag[1:]
	}

	switch tag[:1] {
	case FlagType, PosType, KeyType:
		out.Type = tag[:1]
		tag = tag[1:]
	default:
		return nil, fmt.Errorf("invalid Tag syntax (type) at struct-field: %q - %s", name, tag)
	}

	elements := strings.SplitN(tag, ":", 2)

	if len(elements) != 2 {
		return nil, fmt.Errorf("invalid Tag syntax at struct-field: %q - %s", name, tag)
	}

	ident, deriver := elements[0], elements[1]

	identItems := strings.Split(ident, "|")

	if len(identItems) > 2 {
		return nil, fmt.Errorf("invalid Tag syntax (identifier) at struct-field: %q - %s", name, tag)
	} else if len(identItems) == 2 {
		ident = identItems[0]
		out.Alias = identItems[1]
	}

	switch out.Type {
	case FlagType, KeyType:
		out.Identifier = ident
	case PosType:

		if strings.HasSuffix(ident, "...") {
			p, err := strconv.ParseInt(strings.TrimSuffix(ident, "..."), 10, 64)

			if err != nil || p < 1 {
				return nil, fmt.Errorf("invalid Index for Positional Args: %s - %w", ident, err)
			}

			out.IndexFrom = int(p)
			out.Identifier = fmt.Sprintf("%d", out.IndexFrom)
			out.Repeatable = true
		} else {
			p, err := strconv.ParseInt(ident, 10, 64)
			if err != nil || p < 1 {
				return nil, fmt.Errorf("invalid Index for Positional Arg: %s - %w", ident, err)
			}

			out.Index = int(p)
			out.Identifier = fmt.Sprintf("%d", out.Index)
		}

	default:
		panic("unreachable")
	}

	out.Deriver = deriver

	return out, nil
}

// An Argument represents a clap.Arg defined via struct-tags.
type Argument struct {
	Name       string // name of the struct-field
	Deriver    string // name of the Deriver implementation for the Argument
	Identifier string // identifier for the Argument
	Alias      string // (not implemented) single character alias for the Argument (only valid for Flag/KeyValue)
	Default    string // the default value for the Argument
	Type       string // indicates the type of Argument
	Required   bool   // indicates if the Argument is required
	Repeatable bool   // indicates the Argument is repeatable

	// Index of the Arg
	//
	// Only specified for non-repeatable PositionalArg
	Index int

	// Index of the start of the Arg
	//
	// Only specified for repeatable PositionalArg
	IndexFrom int

	// Usage/description of the Arg
	Usage string

	//dflt any // pointer to a Default Value specified by a field of the struct
	ptr any // pointer to the underlying Field of the struct
}

// Field returns a pointer to the underlying variable for the Argument.
func (a *Argument) Field() any {
	return a.ptr
}

func (a *Argument) Pos() int {
	if a.Index > 0 {
		return a.Index
	}

	return a.IndexFrom
}
