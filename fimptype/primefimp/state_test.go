package primefimp

import (
	"testing"

	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func attrVal(vt fimptype.ValueTypeT, val any) StateAttributeValue {
	return StateAttributeValue{Timestamp: "ts", ValType: vt, Val: val, Props: map[string]string{"p": "v"}}
}

func testState() State {
	return State{Devices: StateDevices{
		{Id: 1, Services: []*StateService{
			{Name: "meter_elec", Attributes: []StateAttribute{{Name: "meter"}}},
			{Name: "battery"},
		}},
		{Id: 2},
		{Id: 3, Services: []*StateService{{Name: "battery"}}},
	}}
}

func TestStateFilters(t *testing.T) {
	devs := testState().Devices
	assert.Len(t, devs.FilterDevicesByService("battery"), 2)
	assert.Len(t, devs.FilterDevicesByService("none"), 0)
	assert.Nil(t, StateDevices{}.FilterDevicesByService("battery"))
	assert.Len(t, devs.FilterDevicesByAttribute("meter"), 1)
	assert.Nil(t, StateDevices{}.FilterDevicesByAttribute("meter"))
	assert.Len(t, devs.FilterDevicesByFunc(func(d *StateDevice) bool { return d.Id > 1 }), 2)

	s := testState()
	s.WithFilteredDevices(func(d *StateDevice) bool { return d.Id == 1 })
	require.Len(t, s.Devices, 1)

	sd := s.Devices[0]
	assert.True(t, sd.ContainsService("battery"))
	assert.False(t, sd.ContainsService("none"))
	assert.False(t, devs[1].ContainsService("battery"))
	assert.Len(t, sd.FilterServices([]string{"battery", "meter_elec", "none"}), 2)
	sd.WithFilteredServices(func(ss *StateService) bool { return ss.Name == "battery" })
	assert.Len(t, sd.Services, 1)

	svc := StateService{Attributes: []StateAttribute{{Name: "a"}}}
	_, ok := svc.FindAttribute("a")
	assert.True(t, ok)
	_, ok = svc.FindAttribute("b")
	assert.False(t, ok)
}

func TestStateAttributeGettersEmpty(t *testing.T) {
	sa := StateAttribute{Name: "x"}
	_, err := sa.GetFirstValue()
	require.Error(t, err)
	_, err = sa.GetFirstStringValue()
	require.Error(t, err)
	_, err = sa.GetFirstIntValue()
	require.Error(t, err)
	_, err = sa.GetFirstFloatValue()
	require.Error(t, err)
	_, _, err = sa.GetFirstBinaryValue()
	require.Error(t, err)
	_, err = sa.GetFirstStrArrayValue()
	require.Error(t, err)
	_, err = sa.GetFirstStrMapValue()
	require.Error(t, err)
	_, err = sa.GetFirstIntArrayValue()
	require.Error(t, err)
	_, err = sa.GetFirstIntMapValue()
	require.Error(t, err)
	_, err = sa.GetFirstFloatArrayValue()
	require.Error(t, err)
	_, err = sa.GetFirstFloatMapValue()
	require.Error(t, err)
	_, err = sa.GetFirstBoolMapValue()
	require.Error(t, err)
	_, err = sa.GetFirstPropAsString("p")
	require.Error(t, err)
}

func TestStateAttributeGetters(t *testing.T) {
	first := func(vt fimptype.ValueTypeT, val any) StateAttribute {
		return StateAttribute{Name: "x", Values: []StateAttributeValue{attrVal(vt, val)}}
	}

	v, err := first(fimptype.VTypeString, "s").GetFirstValue()
	require.NoError(t, err)
	assert.Equal(t, "ts", v.Timestamp)

	s, err := first(fimptype.VTypeString, "s").GetFirstStringValue()
	require.NoError(t, err)
	assert.Equal(t, "s", s)

	i, err := first(fimptype.VTypeInt, 5).GetFirstIntValue()
	require.NoError(t, err)
	assert.Equal(t, 5, i)

	f, err := first(fimptype.VTypeFloat, 1.5).GetFirstFloatValue()
	require.NoError(t, err)
	assert.Equal(t, 1.5, f)

	b, ts, err := first(fimptype.VTypeBool, true).GetFirstBinaryValue()
	require.NoError(t, err)
	assert.True(t, b)
	assert.Equal(t, "ts", ts)

	sa, err := first(fimptype.VTypeStrArray, []string{"a", "b"}).GetFirstStrArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, sa)

	sm, err := first(fimptype.VTypeStrMap, map[string]string{"k": "v"}).GetFirstStrMapValue()
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"k": "v"}, sm)

	ia, err := first(fimptype.VTypeIntArray, []int{1, 2}).GetFirstIntArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2}, ia)

	im, err := first(fimptype.VTypeIntMap, map[string]int{"k": 1}).GetFirstIntMapValue()
	require.NoError(t, err)
	assert.Equal(t, map[string]int{"k": 1}, im)

	fa, err := first(fimptype.VTypeFloatArray, []float64{1.5}).GetFirstFloatArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []float64{1.5}, fa)

	fm, err := first(fimptype.VTypeFloatMap, map[string]float64{"k": 1.5}).GetFirstFloatMapValue()
	require.NoError(t, err)
	assert.Equal(t, map[string]float64{"k": 1.5}, fm)

	bm, err := first(fimptype.VTypeBoolMap, map[string]bool{"k": true}).GetFirstBoolMapValue()
	require.NoError(t, err)
	assert.Equal(t, map[string]bool{"k": true}, bm)

	p, err := first(fimptype.VTypeString, "s").GetFirstPropAsString("p")
	require.NoError(t, err)
	assert.Equal(t, "v", p)
	_, err = first(fimptype.VTypeString, "s").GetFirstPropAsString("none")
	require.Error(t, err)
}

func TestStateAttributeValueMismatch(t *testing.T) {
	getters := map[string]func(StateAttributeValue) error{
		"string":     func(v StateAttributeValue) error { _, err := v.GetStringValue(); return err },
		"int":        func(v StateAttributeValue) error { _, err := v.GetIntValue(); return err },
		"float":      func(v StateAttributeValue) error { _, err := v.GetFloatValue(); return err },
		"bool":       func(v StateAttributeValue) error { _, err := v.GetBoolValue(); return err },
		"strArray":   func(v StateAttributeValue) error { _, err := v.GetStrArrayValue(); return err },
		"strMap":     func(v StateAttributeValue) error { _, err := v.GetStrMapValue(); return err },
		"intArray":   func(v StateAttributeValue) error { _, err := v.GetIntArrayValue(); return err },
		"floatArray": func(v StateAttributeValue) error { _, err := v.GetFloatArrayValue(); return err },
		"floatMap":   func(v StateAttributeValue) error { _, err := v.GetFloatMapValue(); return err },
		"intMap":     func(v StateAttributeValue) error { _, err := v.GetIntMapValue(); return err },
		"boolMap":    func(v StateAttributeValue) error { _, err := v.GetBoolMapValue(); return err },
	}
	// the value type is "null", so parse leaves Val untouched and the typed assertion fails
	for name, get := range getters {
		require.Error(t, get(attrVal("null", nil)), name)
	}
	require.Error(t, getters["int"](attrVal(fimptype.VTypeInt, "x")))
	require.Error(t, getters["float"](attrVal(fimptype.VTypeFloat, "x")))
	require.Error(t, getters["bool"](attrVal(fimptype.VTypeBool, "x")))
	require.Error(t, getters["strArray"](attrVal(fimptype.VTypeStrArray, "x")))
	require.Error(t, getters["intMap"](attrVal(fimptype.VTypeIntMap, map[string]string{"k": "x"})))
	require.Error(t, getters["floatMap"](attrVal(fimptype.VTypeFloatMap, map[string]string{"k": "x"})))
	require.Error(t, getters["boolMap"](attrVal(fimptype.VTypeBoolMap, map[string]string{"k": "x"})))
}
