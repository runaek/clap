package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"os"
)

var (
	System = New("SYSTEM").
		AddFlag(NewFlagP[bool](nil, "help", "h", parse.Bool{}),
			AsOptional(), WithUsage("Display the help-text/usage of the program."))
)

func Ok() {
	System.Ok()
}

func Parse() error {
	System.Parse(os.Args[1:]...)
	return System.Err()
}

func SetName(name string) {
	System.Name = name
}

func SetDescription(description string) {
	System.Description = description
}

func Add(args ...Arg) error {
	System.Add(args...)
	return System.Valid()
}

func Err() error {
	return System.Err()
}

func Valid() error {
	return System.Valid()
}

func RawPositions() []string {
	return System.RawPositions()
}

func RawKeyValues() map[string][]string {
	return System.RawKeyValues()
}

func RawFlags() map[string][]string {
	return System.RawFlags()
}
