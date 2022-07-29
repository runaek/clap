package main

import (
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/derive"
	"github.com/runaek/clap/pkg/parse"
	"io"
	"os"
)

var (
	tpl     = derive.NewTemplate()
	output  string
	outputF = clap.NewFlagP[string](&output, "output", "o", parse.String{},
		clap.WithUsage("(Optional) A output path to write the generated code to."))
	parser = clap.Must("codegen", outputF, &tpl).
		WithDescription(description)
)

func main() {
	parser.Parse()
	parser.Ok()

	var out io.Writer

	if output == "" {
		out = os.Stdout
	} else {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, os.ModePerm)

		if err != nil {
			panic(err)
		}

		out = f
		defer f.Close()
	}

	if err := tpl.Process(out); err != nil {
		panic(err)
	}

}

const (
	description = `codegen is a code-generating tool for derived clap.Arg implementations from struct-tags.

The program will produce a Go-file containing:

	* clap.FlagDeriver, clap.KeyValueDeriver and clap.PositionalDeriver implementations for some
	  user-defined type T which has a parse.Parse[T] implementation

	* an init() function to register the different Deriver implementations

The package containing the implementations should be imported for effects into any programs 
that want to derive the arguments from struct-tags.`
)
