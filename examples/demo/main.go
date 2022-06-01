package main

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parsers"
	"os"
)

var (
	name    string
	nameArg = clap.NewKeyValue[string](&name, "name", parsers.String,
		clap.WithDefault("Obi-Wan Kenobi"),
		clap.WithShorthand('n'),
		clap.WithUsage("Enter your name!"))

	debug     bool
	debugFlag = clap.NewFlagP[bool](&debug, "debug", 'V', parsers.Bool,
		clap.WithUsage("Toggle DEBUG mode."),
		clap.WithDefault("false"),
	)

	dev     bool
	devFlag = clap.NewFlagP[bool](&dev, "dev", 'd', parsers.Bool,
		clap.WithUsage("Toggle DEV mode."),
		clap.WithDefault("true"),
	)

	counter     clap.Counter
	counterFlag = clap.NewFlagP[clap.Counter](&counter, "counter", 'c', clap.CounterParser,
		clap.WithUsage("Increment the counter."))

	funcName    string
	funcNamePos = clap.NewPosition[string](&funcName, 1, parsers.String,
		clap.WithUsage("Enter the name of a function to call."),
		clap.WithDefault("greet"))

	args    []string
	argsPos = clap.NewPositions(&args, 2, parsers.String,
		clap.WithUsage("Enter the arguments for the function."))

	csvArgs []string
	csvPipe = clap.CSVPipe(&csvArgs, parsers.Strings)
)

func main() {

	parser := clap.Must("demo", clap.ContinueOnError).
		AddFlag(debugFlag).
		AddFlag(counterFlag).
		AddFlag(devFlag).
		AddPosition(funcNamePos).
		AddPosition(argsPos).
		AddKeyValue(nameArg).
		AddPipe(csvPipe).
		Ok()

	parser.Parse(os.Args)
	fmt.Printf("name:    %s\n", name)
	fmt.Printf("Debug:   %t\n", debug)
	fmt.Printf("Dev:     %t\n", dev)
	fmt.Printf("Func:    %s\n", funcName)
	fmt.Printf("Args:    %s\n", args)
	fmt.Printf("Counter: %d\n", counter)

	fmt.Println(parser.RawPositions())
	fmt.Println(parser.RawFlags())
	fmt.Println(parser.RawKeyValues())
}
