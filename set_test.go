package clap

import (
	"github.com/runaek/clap/pkg/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_AddFlag(t *testing.T) {

	var testCases = map[string]struct {
		Run          func(*Set, *assert.Assertions)
		Exists       []Identifier
		DoesNotExist []Identifier
	}{
		"SameNameReturnsErrDuplicate": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					f1 = NewFlagsP[string](nil, "f1", "f", parse.String{})
				)

				a.NoError(set.AddFlag(f1))
				a.ErrorIs(set.AddFlag(f1), ErrDuplicate)
			},
			Exists: []Identifier{
				Flag("f1"),
				Flag("f"),
			},
		},
		"DuplicateAliasReturnsErrDuplicate": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					f1 = NewFlagsP[string](nil, "f1", "f", parse.String{})
					f2 = NewFlagP[string](nil, "f2", "f", parse.String{})
				)

				a.NoError(set.AddFlag(f1))
				a.ErrorIs(set.AddFlag(f2), ErrDuplicate)
				//a.True(errors.Is(set.AddFlag(f2), ErrDuplicate))
			},
			Exists: []Identifier{
				Flag("f1"),
				Flag("f"),
			},
			DoesNotExist: []Identifier{
				Flag("f2"),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			s := NewSet()

			if tc.Run == nil {
				a.Fail("Test has nil Run")
			} else {
				tc.Run(s, a)
			}

			for _, shouldExist := range tc.Exists {
				a.NotNil(s.Get(shouldExist))
			}

			for _, shouldNotExist := range tc.DoesNotExist {
				a.Nil(s.Get(shouldNotExist))
			}

		})
	}
}

func TestSet_AddKeyValue(t *testing.T) {
	var testCases = map[string]struct {
		Run          func(*Set, *assert.Assertions)
		Exists       []Identifier
		DoesNotExist []Identifier
	}{
		"SameNameReturnsErrDuplicate": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					k1 = NewKeyValue[string](nil, "key1", parse.String{}, WithShorthand("k"))
				)

				a.NoError(set.AddKeyValue(k1))
				a.ErrorIs(set.AddKeyValue(k1), ErrDuplicate)
				//a.True(errors.Is(set.AddKeyValue(k1), ErrDuplicate))
			},
			Exists: []Identifier{
				Key("key1"),
				Key("k"),
			},
		},
		"DuplicateAliasReturnsErrDuplicate": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					k1 = NewKeyValue[string](nil, "key1", parse.String{}, WithShorthand("k"))
					k2 = NewKeyValue[string](nil, "key2", parse.String{}, WithShorthand("k"))
				)

				a.NoError(set.AddKeyValue(k1))
				a.ErrorIs(set.AddKeyValue(k2), ErrDuplicate)
				//a.True(errors.Is(set.AddKeyValue(k2), ErrDuplicate))
			},
			Exists: []Identifier{
				Key("key1"),
				Key("k"),
			},
			DoesNotExist: []Identifier{
				Key("key2"),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			s := NewSet()

			if tc.Run == nil {
				a.Fail("Test has nil Run")
			} else {
				tc.Run(s, a)
			}

			for _, shouldExist := range tc.Exists {
				a.NotNil(s.Get(shouldExist))
			}

			for _, shouldNotExist := range tc.DoesNotExist {
				a.Nil(s.Get(shouldNotExist))
			}

		})
	}
}

func TestSet_AddPosition(t *testing.T) {
	var testCases = map[string]struct {
		Run          func(*Set, *assert.Assertions)
		Exists       []Identifier
		DoesNotExist []Identifier
	}{
		"DuplicateIndexReturnsErr": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					p1  = NewPosition[string](nil, 1, parse.String{})
					p2  = NewPosition[string](nil, 2, parse.String{})
					p22 = NewPosition[string](nil, 2, parse.String{})
				)

				a.NoError(set.AddPosition(p1))
				a.NoError(set.AddPosition(p2))
				a.ErrorIs(set.AddPosition(p22), ErrDuplicate)
				//a.True(errors.Is(set.AddPosition(p22), ErrDuplicate))
			},
			Exists: []Identifier{
				Position(1),
				Position(2),
			},
		},
		"PositionAfterVariadicReturnsErr": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					p1 = NewPosition[string](nil, 1, parse.String{})
					p2 = NewPosition[string](nil, 2, parse.String{})
					p3 = NewPositions[string](nil, 3, parse.String{})
					p4 = NewPosition[string](nil, 4, parse.String{})
				)

				a.NoError(set.AddPosition(p1))
				a.NoError(set.AddPosition(p2))
				a.NoError(set.AddPosition(p3))
				a.ErrorIs(set.AddPosition(p4), ErrInvalidIndex)
				//e := set.AddPosition(p4)
				//a.Truef(errors.Is(e, ErrInvalidIndex), "unexpected error: want %s but got %s", ErrInvalidIndex, e)

			},
			Exists: []Identifier{
				Position(1),
				Position(2),
				Position(3),
			},
		},
		"RequiredCannotComeAfterOptional": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					p1 = NewPosition[string](nil, 1, parse.String{})
					p2 = NewPosition[string](nil, 2, parse.String{}, AsRequired())
				)

				a.NoError(set.AddPosition(p1))
				e := set.AddPosition(p2)
				a.ErrorIs(e, ErrInvalidIndex)
			},
			Exists: []Identifier{
				Position(1),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			s := NewSet()

			if tc.Run == nil {
				a.Fail("Test has nil Run")
			} else {
				tc.Run(s, a)
			}

			for _, shouldExist := range tc.Exists {
				a.NotNil(s.Get(shouldExist))
			}

			for _, shouldNotExist := range tc.DoesNotExist {
				a.Nil(s.Get(shouldNotExist))
			}

		})
	}
}

func TestSet_AddPipe(t *testing.T) {
	var testCases = map[string]struct {
		Run          func(*Set, *assert.Assertions)
		Exists       []Identifier
		DoesNotExist []Identifier
	}{
		"OnlyOnePipe": {
			Run: func(set *Set, a *assert.Assertions) {

				var (
					p1 = NewPipeArg[string](nil, parse.String{}, nil, nil)
					p2 = NewPipeArg[string](nil, parse.String{}, nil, nil)
				)

				a.NoError(set.AddPipe(p1))
				a.ErrorIs(set.AddPipe(p2), ErrPipeExists)
				//e := set.AddPipe(p2)
				//a.Truef(errors.Is(e, ErrPipeExists), "unexpected error: want %s but got %s", ErrPipeExists, e)
			},
			Exists: []Identifier{
				Pipe(""),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			s := NewSet()

			if tc.Run == nil {
				a.Fail("Test has nil Run")
			} else {
				tc.Run(s, a)
			}

			for _, shouldExist := range tc.Exists {
				a.NotNil(s.Get(shouldExist))
			}

			for _, shouldNotExist := range tc.DoesNotExist {
				a.Nil(s.Get(shouldNotExist))
			}

		})
	}
}
