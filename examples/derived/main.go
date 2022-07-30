package main

import (
	"fmt"
	"github.com/runaek/clap"
	_ "github.com/runaek/clap/pkg/derive"
	"github.com/runaek/clap/pkg/parse"
)

type MyProgram struct {
	Name       string `cli:"!#1|name:string"`
	NameUsage  string
	Roles      []string `cli:"#2...|roles:strings"`
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
