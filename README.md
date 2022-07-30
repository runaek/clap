![Tests (dev)](https://github.com/runaek/clap/workflows/test/badge.svg?branch=dev "Development")
![Tests](https://github.com/runaek/clap/workflows/test/badge.svg "Tests" )
---
## Command-Line Argument Parser (In Development)

`clap` is a command-line argument parser for Go. It loosely mimics/extends the API in the standard library package `flag`.
The idea is that you can easily define/use different types of arguments in your program:

```go
var (
    myString string
    myArg = clap.NewKeyValue[string](&myString, "my-arg", parse.String{})
		
    parser = clap.Must("my-program", myArg)
)

    
func main() {
    parser.Parse(nil)
    parser.Ok()
    // do stuff with myArg
}   
```

### Usage
`go get -u github.com/runaek/clap` (requires go1.18+)

### Features

* Supports key-value `KeyValueArg`, flag `FlagArg`, positional `PositionalArg` and pipe `PipeArg` arguments as input;
* Supports extensions to user-defined types;
* Auto-generate `Arg(s)` based on struct-definitions;

### Examples

#### Arguments derived from struct-tags

Struct-tags can be supplied in the following formats: 

    > `cli:"#<index>:<deriver>"` for a positional argument
    > `cli:"#<start>...:<deriver>"` for positional arguments
    > `cli:"-<flag_name>|<flag_shorthand>:<deriver>"` for a flag argument
    > `cli:"@<key_name>|<key_shorthand>:<deriver>"` for a key-value argument

A `!` can be prefixed to tags to mark the input as `Required` (this will not apply to variable positional arguments).

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