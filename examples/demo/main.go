package main

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/parse"
	"os"
)

var (
	name    string
	nameArg = clap.NewKeyValue[string](&name, "name", parse.String{},
		clap.WithDefault("Obi-Wan Kenobi"),
		clap.WithShorthand('n'),
		clap.WithUsage("Enter your name!"))

	debug     bool
	debugFlag = clap.NewFlagP[bool](&debug, "debug", 'V', parse.Bool{},
		clap.WithUsage("Toggle DEBUG mode."),
		clap.WithDefault("false"),
	)

	dev     bool
	devFlag = clap.NewFlagP[bool](&dev, "dev", 'd', parse.Bool{},
		clap.WithUsage("Toggle DEV mode."),
		clap.WithDefault("false"),
	)

	ident  string
	idFlag = clap.NewFlag[string](&ident, "id", parse.String{},
		clap.WithUsage("Choose an id."),
		clap.AsRequired())

	counter     parse.C
	counterFlag = clap.NewFlagP[parse.C](&counter, "counter", 'c', parse.Counter{},
		clap.WithUsage("Increment the counter."))

	funcName    string
	funcNamePos = clap.NewPosition[string](&funcName, 1, parse.String{},
		clap.WithUsage("Enter the name of a function to call."),
		clap.WithDefault("greet"))

	args    []string
	argsPos = clap.NewPositions[string](&args, 2, parse.String{},
		clap.WithUsage("Enter the arguments for the function."))

	csvArgs []string
	csvPipe = clap.CSVPipe[[]string](&csvArgs, parse.Strings{})
)

func main() {

	parser := clap.Must("demo", clap.ContinueOnError).
		AddFlag(debugFlag).
		AddFlag(counterFlag).
		AddFlag(devFlag).
		AddFlag(idFlag, clap.AsOptional()).
		AddPosition(funcNamePos).
		AddPosition(argsPos).
		AddKeyValue(nameArg).
		AddPipe(csvPipe).
		Ok()

	parser.Parse(os.Args)
	//parser.Ok()
	fmt.Printf("name:    %s\n", name)
	fmt.Printf("Debug:   %t\n", debug)
	fmt.Printf("Dev:     %t\n", dev)
	fmt.Printf("Func:    %s\n", funcName)
	fmt.Printf("Args:    %s\n", args)
	fmt.Printf("CSVArgs: %s\n", csvArgs)
	fmt.Printf("C: %d\n", counter)
	fmt.Printf("Id:      %s\n", ident)

	fmt.Println(parser.Err())

	fmt.Println(parser.RawPositions())
	fmt.Println(parser.RawFlags())
	fmt.Println(parser.RawKeyValues())

}
