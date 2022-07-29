package main

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
)

var (
	parser = clap.Must("demo").
		Add(debugFlag, counterFlag, devFlag, idFlag, funcNamePos, argsPos, nameArg, csvPipe)
)

func main() {
	parser.Parse()
	parser.Ok()

	fmt.Printf("Name:    %s\n", name)
	fmt.Printf("Debug:   %t\n", debug)
	fmt.Printf("Dev:     %t\n", dev)
	fmt.Printf("Func:    %s\n", funcName)
	fmt.Printf("Args:    %s\n", args)
	fmt.Printf("CSVArgs: %s\n", csvArgs)
	fmt.Printf("Counter: %d\n", counter)
	fmt.Printf("Id:      %s\n", ident)
}

var (
	name    string
	nameArg = clap.NewKeyValue[string](&name, "name", parse.String{},
		clap.WithDefault("Obi-Wan Kenobi"),
		clap.WithShorthand("n"),
		clap.WithUsage("Enter your name!"))

	debug     bool
	debugFlag = clap.NewFlagP[bool](&debug, "debug", "V", parse.Bool{},
		clap.WithUsage("Toggle DEBUG mode."),
		clap.WithDefault("false"),
	)

	dev     bool
	devFlag = clap.NewFlag[bool](&dev, "dev", parse.Bool{},
		clap.WithUsage("Toggle DEV mode."),
		clap.WithDefault("false"),
	)

	ident  string
	idFlag = clap.NewFlagP[string](&ident, "id", "I", parse.String{},
		clap.WithUsage("Choose an id."),
		clap.AsRequired())

	counter     parse.C
	counterFlag = clap.NewFlagP[parse.C](&counter, "counter", "c", parse.Counter{},
		clap.WithUsage("Increment the counter."))

	funcName    string
	funcNamePos = clap.NewPosition[string](&funcName, 1, parse.String{},
		clap.WithUsage("Enter the name of a function to call."),
		clap.WithDefault("greet"))

	args    []string
	argsPos = clap.NewPositions[string](&args, 2, parse.String{},
		clap.WithUsage("Enter the arguments for the function."))

	csvArgs []string
	csvPipe = clap.CSVPipe[[]string](&csvArgs, parse.Strings{},
		clap.WithUsage("Comma-separated data from a pipe."))
)
