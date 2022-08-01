package main

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/runaek/clap/pkg/flag"
	"github.com/runaek/clap/pkg/keyvalue"
	"github.com/runaek/clap/pkg/pos"
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

	fmt.Printf("ID:     %s\n", id)
	fmt.Printf("Number: %d\n", num)
	fmt.Printf("Words:  %s\n", strings.Join(words, ", "))
	// do stuff with myArg
}
