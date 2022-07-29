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
	outputF = clap.NewFlag[string](&output, "output", parse.String{},
		clap.WithUsage("(Optional) A output path to write the generated code to."))

	parser = clap.Must("codegen", outputF, &tpl)
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
