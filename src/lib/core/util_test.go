package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDumpFail(t *testing.T) {
	for _, th := range []*Thunk{
		OutOfRangeError(),
	} {
		v := PApp(Dump, th).Eval()
		t.Log(v)

		if _, ok := v.(ErrorType); !ok {
			t.Fail()
		}
	}
}

func TestInternalStrictDumpPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	strictDump(nil)
}

func TestInternalStrictDumpFail(t *testing.T) {
	for _, th := range []*Thunk{
		NewList(OutOfRangeError()),
		NewDictionary([]Value{Nil.Eval()}, []*Thunk{OutOfRangeError()})} {
		if _, err := strictDump(th.Eval()); err == nil {
			t.Fail()
		}
	}
}

func TestStrictDump(t *testing.T) {
	for _, th := range []*Thunk{
		Nil,
		True,
		False,
		EmptyList,
		EmptyDictionary,
		NewNumber(42),
		NewString("foo"),
		NewOutput(Nil),
	} {
		s, err := StrictDump(th.Eval())
		assert.NotEqual(t, "", s)
		assert.Equal(t, (*Thunk)(nil), err)
	}
}

func TestIdentity(t *testing.T) {
	for _, th := range []*Thunk{Nil, NewNumber(42), True, False, NewString("foo")} {
		assert.True(t, testEqual(PApp(identity, th), th))
	}
}

func TestTypeOf(t *testing.T) {
	for _, test := range []struct {
		typ   string
		thunk *Thunk
	}{
		{"nil", Nil},
		{"list", NewList()},
		{"list", EmptyList},
		{"bool", True},
		{"number", NewNumber(123)},
		{"string", NewString("foo")},
		{"dict", NewDictionary(nil, nil)},
		{"error", NewError("MyError", "This is error.")},
		{"function", identity},
		{"function", PApp(Partial, identity)},
	} {
		v := PApp(TypeOf, test.thunk).Eval()
		t.Log(v)
		assert.Equal(t, test.typ, string(v.(StringType)))
	}
}

func TestTypeOfFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		} else if _, ok := r.(error); !ok {
			t.Fail()
		}
	}()

	PApp(TypeOf, Normal("foo")).Eval()
}
