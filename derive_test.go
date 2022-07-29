//go:build clap_mocks

package clap

import (
	"github.com/golang/mock/gomock"
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AllOptionalArgs struct {
	Value1 string `cli:"@k1:string"`
	Value2 string `cli:"-f1:string"`
}

func nf(name string, opts ...Option) IFlag {
	return NewFlag[string](nil, name, parse.String{}, opts...)
}

func nk(name string, opts ...Option) IKeyValue {
	return NewKeyValue[string](nil, name, parse.String{}, opts...)
}

func TestDerive(t *testing.T) {

	ctrl := gomock.NewController(t)
	fDeriver := NewMockFlagDeriver(ctrl)
	kDeriver := NewMockKeyValueDeriver(ctrl)

	f1 := nf("f1")
	k1 := nk("k1")

	fDeriver.EXPECT().
		DeriveFlag(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(f1, nil)

	kDeriver.EXPECT().
		DeriveKeyValue(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(k1, nil)

	RegisterFlagDeriver("string", fDeriver)
	RegisterKeyValueDeriver("string", kDeriver)

	a := assert.New(t)

	exp := []argName{
		FlagType.getIdentifier("f1"),
		KeyValueType.getIdentifier("k1"),
	}

	optArgs := AllOptionalArgs{}

	t.Logf("Testing Derive: %+v", optArgs)

	act, err := Derive(&optArgs)

	a.NoError(err)

	actNames := make([]argName, len(act))

	for i, ac := range act {
		actNames[i] = ac.argName()
	}

	a.ElementsMatch(exp, actNames)

	p := New("test_parser")

	p.Add(act...)

	a.NoError(p.Valid())

	p.Parse("test", "args", "--f1=my_value", "k1=my_other_value")

	tk1 := k1.(*KeyValueArg[string])
	tf1 := f1.(*FlagArg[string])

	a.Equal(tk1.Variable().Unwrap(), "my_other_value")
	a.Equal(tf1.Variable().Unwrap(), "my_value")

	a.NoError(p.Err())

}
