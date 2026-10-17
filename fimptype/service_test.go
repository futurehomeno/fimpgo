package fimptype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceMsgType(t *testing.T) {
	for in, want := range map[string]MsgTypeT{
		"cmd.lvl.set": MsgTypeCmd, "evt.lvl.report": MsgTypeEvt, "  cmd.x": MsgTypeCmd, "rsp.x": MsgTypeRsp, "foo": MsgTypeUnknown, "": MsgTypeUnknown,
	} {
		assert.Equal(t, want, InterfaceMsgType(in), in)
	}
	assert.Equal(t, "cmd", MsgTypeCmd.Str())
}

func TestService_EnsureInterfaces(t *testing.T) {
	a := Interface{Type: TypeIn, MsgType: "cmd.a", ValueType: VTypeNull, Version: "1"}
	b := Interface{Type: TypeOut, MsgType: "evt.b", ValueType: VTypeInt, Version: "1"}
	s := &Service{Interfaces: []Interface{a}}
	s.EnsureInterfaces(a, b, b)
	assert.Equal(t, []Interface{a, b}, s.Interfaces)
}

func TestService_Properties(t *testing.T) {
	s := &Service{Props: map[string]any{
		"strs": []string{"a"}, "strsAny": []any{"a", "b"}, "str": "x", "num": 1.5, "floats": []float64{1, 2},
		"floatsAny": []any{1, 2.5}, "int": 3, "intF": 4.0, "ints": []int{1}, "intsAny": []any{1.0, 2.0},
		"bool": true, "boolBad": "yes", "obj": map[string]any{"k": "v"}, "chan": make(chan int),
	}}

	assert.Equal(t, []string{"a"}, s.PropertyStrings("strs"))
	assert.Equal(t, []string{"a", "b"}, s.PropertyStrings("strsAny"))
	assert.Nil(t, s.PropertyStrings("missing"))
	assert.Nil(t, s.PropertyStrings("chan"))

	assert.Equal(t, "x", s.PropertyString("str"))
	assert.Empty(t, s.PropertyString("num"))
	assert.Empty(t, s.PropertyString("missing"))

	assert.Equal(t, []float64{1, 2}, s.PropertyFloats("floats"))
	assert.Equal(t, []float64{1, 2.5}, s.PropertyFloats("floatsAny"))
	assert.Nil(t, s.PropertyFloats("missing"))
	assert.Nil(t, s.PropertyFloats("str"))

	f, ok := s.PropertyFloat("num")
	assert.True(t, ok)
	assert.Equal(t, 1.5, f)
	f, ok = s.PropertyFloat("int")
	assert.True(t, ok)
	assert.Equal(t, 3.0, f)
	_, ok = s.PropertyFloat("str")
	assert.False(t, ok)
	_, ok = s.PropertyFloat("missing")
	assert.False(t, ok)

	assert.Equal(t, []int{1}, s.PropertyIntegers("ints"))
	assert.Equal(t, []int{1, 2}, s.PropertyIntegers("intsAny"))
	assert.Nil(t, s.PropertyIntegers("missing"))
	assert.Nil(t, s.PropertyIntegers("str"))

	for name, want := range map[string]int{"int": 3, "intF": 4} {
		i, ok := s.PropertyInteger(name)
		assert.True(t, ok, name)
		assert.Equal(t, want, i, name)
	}
	for _, name := range []string{"num", "str", "missing"} {
		_, ok = s.PropertyInteger(name)
		assert.False(t, ok, name)
	}

	assert.True(t, s.PropertyBool("bool"))
	assert.False(t, s.PropertyBool("boolBad"))
	assert.False(t, s.PropertyBool("missing"))

	var obj struct{ K string }
	assert.True(t, s.PropertyObject("obj", &obj))
	assert.Equal(t, "v", obj.K)
	assert.False(t, s.PropertyObject("missing", &obj))
	assert.False(t, s.PropertyObject("chan", &obj))
}
