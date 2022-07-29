//go:build clap_mocks

package derive

import (
	"fmt"
	"github.com/runaek/clap"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stringTestData struct {
	Value1        string `cli:"!@k1:string"`
	Value1Usage   string
	Value1Default string
	Value2        string `cli:"-f1:string"`
	Value2Usage   string
	Value3        string `cli:"-f2:string"`
	Value3Usage   string
	Value3Default string
	Args          []string `cli:"#1...:strings"`
	ArgsUsage     string
}

func TestDerive_String(t *testing.T) {

	a := assert.New(t)

	targetStrukt := stringTestData{
		Value1Usage:   "value-1 usage",
		Value2Usage:   "value-2 usage",
		Value3Usage:   "value-3 usage",
		Value3Default: "hello",
		ArgsUsage:     "args usage",
	}

	args, err := clap.Derive(&targetStrukt)

	a.NoError(err)
	a.Len(args, 4)

	p := clap.New("test_parser").Add(args...)

	p.Parse("test", "args", "--f1=my_value", "k1=my_other_value")
	a.NoError(p.Err())

	a.Equal("my_value", targetStrukt.Value2)
	a.Equal("my_other_value", targetStrukt.Value1)
	a.Equal([]string{"test", "args"}, targetStrukt.Args)

	for _, arg := range args {
		t.Logf("Identified Argument: %s %s %s", arg.Name(), arg.Type(), arg.Usage())
		msg := fmt.Sprintf("unexpected usage string for arg: %s", arg.Name())
		switch arg.Name() {
		case "k1":
			a.Equal("value-1 usage", arg.Usage(), msg)
			a.True(arg.IsRequired(), "expected k1 to be Required")
		case "f2":
			a.Equal("value-3 usage", arg.Usage(), msg)
			a.Equal("hello", arg.Default(), "expected f2 to have Default Value")
		case "f1":
			a.Equal("value-2 usage", arg.Usage(), msg)
			a.False(arg.IsRequired(), "expected f1 to *not* be Required")
		case "1":
			a.Equal("args usage", arg.Usage(), msg)
			a.True(arg.IsRepeatable(), "expected positional arguments to be Repeatable")
		default:
			a.Failf("unexpected Argument", "%v", arg)
		}
	}

}
