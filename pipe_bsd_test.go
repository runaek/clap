//go:build aix || darwin || (js && wasm) || (solaris && !illumos)

package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParser_Pipe(t *testing.T) {

	a := assert.New(t)
	p := NewParser("test-pipe-parser", ContinueOnError)

	in, out, pipeErr := os.Pipe()
	a.NoError(pipeErr, "error setting up Pipe")

	_, writeErr := out.WriteString("value,uno")
	a.NoError(writeErr)

	p.Stdin = in

	var (
		f1 = NewFlagP[bool](nil, "a1", "a", parse.Bool{})
		f2 = NewFlagP[bool](nil, "b1", "b", parse.Bool{})
		p1 = CSVPipe[[]string](nil, parse.Strings{})
	)

	p1.input = in

	p.Add(f1, f2, p1)
	a.NoError(p.Valid())

	p.Parse("cmd", "1234", "-ab")

	a.Equal(true, f1.Variable().Unwrap())
	a.Equal(true, f2.Variable().Unwrap())
	a.Equal([]string{"value", "uno"}, p1.Variable().Unwrap())

}
