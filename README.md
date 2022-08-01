![Tests (dev)](https://github.com/runaek/clap/workflows/test/badge.svg?branch=dev "Development")
![Tests](https://github.com/runaek/clap/workflows/test/badge.svg "Tests" )
---
## Command-Line Argument Parser (In Development)

`clap` is a command-line argument parser for Go. 

> Inspired by the standard [flag](https://pkg.go.dev/flag) package and the [clap](https://crates.io/crates/clap) Rust crate.

#### Motivation
I loved the idea behind the `flag` package, but I always found myself wanting more from this package (and any packages that extended it):
* Ability to use this same API for **any** data-type;
* Ability to use this same logic for 'different' types of input;
* Ability to access the internals of the arguments for complex use-cases;

This package aims to solve this issue by providing a unified API for dealing with 3 (4, with 1 unstable) different types
of `Arg` - each of which is bound to some Go variable and can be added to a `Parser` to allow it to be parsed from user input.

The general constructor API is pretty much identical to the `flag` package, but with Type instantiations:
```go
var (
	myValue string 
	myFlag = clap.NewFlag[string](&myValue, "my-flag-value", parse.String{})
	//myKey = clap.NewFlag[string](&myValue, "my-key-value", parse.String{})
	//myPos = clap.NewPosition[string](&myValue, 1, parse.String{})
)
```
Note: the packages `github.com/runaek/clap/pkg/{flag,keyvalue,pos}` provide a number of convenient constructors for some 
common-types which does not require the generic type instantiation or the `parse.Parser[T]` being specified.  

### Features

* Standard API for command-line arguments;
* Standard `parse.Parse[T]` interface for extending arguments to custom data-types;
* Supports dynamic generation of `Arg(s)` via struct-tags;
* Supports a number if different command-line inputs types:
  * `<key>=<vallue` style arguments: `KeyValueArg[T]`;
  * `--<flag>=<value>, -<f>=<value>` style arguments: `FlagArg[T]`;
  * `<arg1> <arg2> <argRemaining ...>` positional style arguments: `PositionalArg[T]`
  * `<some data> | my_program` pipe style arguments: `PipeArg[T]` (unstable);


### Usage
`go get -u github.com/runaek/clap` (requires go1.18+)

```go
// examples/readme/main.go
package main

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/flag"
	"github.com/runaek/clap/pkg/pos"
	"github.com/runaek/clap/pkg/keyvalue"
	"strings"
)

var (
	id     string
	idFlag = flag.String(&id, "id", clap.WithAlias("i"), clap.WithDefault("default-id"))

	words    []string
	wordsArg = pos.Strings(&words, 1, clap.WithAlias("words"))

	num    int
	numArg = keyvalue.Int(&num, "number", clap.AsRequired())

	parser = clap.Must("my-program", idFlag, numArg, wordsArg)
)

func main() {
	parser.Parse()
	parser.Ok()

	fmt.Printf("ID:     %s", id)
	fmt.Printf("Number: %d", num)
	fmt.Printf("Words:  %s", strings.Join(words, ", "))
	// do stuff with args 
}   
```

## Types 

The package provides a number of types which can be used to configure a Go program that relies on command-line input.

See `github.com/runaek/cliche` for more advanced usage. 

### Parser

The `Parser` (not to be confused with the `Parser[T]` interface from the `github.com/runaek/clap/pkg/parse` package) is 
the handler/container for a number of command-line arguments. You can add arguments to a parser, validate they are 
correct, and then you can parse command-line input with it.

The API and behaviour is very similar `flag.FlagSet` in idea, but a direct conversion is likely not possible.

### Arguments

Arguments are defined by the `Arg` interface - there are 4 distinct implementations. The main difference between all the 
implementations is where the source of their string-input comes from. 

There is also the generic `TypedArg[T]` interface, which is satisfied by all `Arg`, but less convenient to work with when 
we want to store them in a `Parser`, but if you know the type, you will always be able to extract the `TypedArg[T]`. 

Package `github.com/runaek/clap/parse` provides a generic interface which is used to parse user-supplied string (and 
default values) into concrete types: 
```go
type Parser[T any] interface {
	Parse(...string) (T, error)
}
```
At runtime, the `Parser[T]` will be provided all the input that maps to the `Arg` and the resulting type (error withstanding) will
then be assigned to the underlying variable.

A number of implementations for built-in (and some `clap`-specific) types for this `Parser` are also contained
within the package.

#### Flag: `FlagArg[T]`

Flag arguments are your traditional flag arguments as you would find in the `flag` package or any other programming language.

#### KeyValue: `KeyValueArg[T]`

Key-value arguments are the same as flag arguments, except they do not expect to be prefixed with a `--` or `-`. 

**Note: key-value args are treated essentially the same as flags, however, it would always expect a value (i.e. not ideal for
bool data types).**

#### Positional: `PositionalArg[T]`

Positional arguments are the anonymous arguments supplied to the program. 

**Note: positional arguments can only be required if the argument preceding them is also required, and 
variable numbers of positional arguments cannot be required.**

#### Pipe: `PipeArg[T]`

> **Note: this is still under development and the API/implementation may change for this type.**

Pipe arguments are mainly a convenience - the intended use was to allow them to 'replace' other argument types. The way
the implementation works is each `PipeArg` has a `Piper` which is responsible for consuming a reader (provided by the pipe)
and parsing string-input. This outputted string input is passed down into the underlying `parse.Parser[T]` just like any
other argument. 

Currently, the package provides a `SeparatedValuePiper` for piping some string-separated data, more implementations will come in time, but this can
easily be implemented by users as they wish. 

**Note: There can only ever be one pipe active within a `Parser` at a time.**

## Derived Arguments

Arguments can dynamically be generated through struct-tags. These tags should specify a `deriver` which will be mapped
to some type and if supported, the type will be used by the argument. The `deriver` referenced should be 'registered'
with the same value using `RegisterKeyValueDeriver`, `RegisterPositionalDeriver` and `RegisterFlagDeriver`. 

Struct-tags can be supplied in the following formats:

    > `cli:"#<index>|<human_name>:<deriver>"` for a positional argument
    > `cli:"#<start>...|<human_name>:<deriver>"` for positional arguments
    > `cli:"-<flag_name>|<flag_shorthand>:<deriver>"` for a flag argument
    > `cli:"@<key_name>|<key_shorthand>:<deriver>"` for a key-value argument

A `!` can be prefixed to a tag to mark the input as `Required`.

## Examples

### Arguments derived from struct-tags
```go
// examples/derive/main.go 
package main

import (
	"fmt"
	"github.com/runaek/clap"
	_ "github.com/runaek/clap/pkg/derive"
	"github.com/runaek/clap/pkg/parse"
)

type MyProgram struct {
	Name       string `cli:"!#1:string"`
	NameUsage  string
	Roles      []string `cli:"#2...:strings"`
	RolesUsage string

	Level        int `cli:"-level:int"`
	LevelUsage   string
	LevelDefault string

	Counter      parse.C `cli:"-counter:counter"`
	CounterUsage string

	Account        string `cli:"!@account:string"`
	AccountUsage   string
	AccountDefault string
}

var (
	prog = MyProgram{
		NameUsage:      "Enter your name!",
		RolesUsage:     "Enter a number of random words (roles).",
		LevelUsage:     "Enter a level (number).",
		LevelDefault:   "73",
		CounterUsage:   "A counter for the number of times the flag is supplied.",
		AccountUsage:   "A key-value for the 'account' to be used",
		AccountDefault: "[default]",
	}

	parser = clap.Must("derived_demo", &prog)
)

func main() {
	parser.Parse()
	parser.Ok()

	fmt.Printf("MyProgram Arguments:\n%+v\n", prog)
}
```

### Custom Deriver 

### Custom Parser
