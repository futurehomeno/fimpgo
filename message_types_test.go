package fimpgo

import (
	"math"
	"testing"
	"time"

	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkGetter[T any](t *testing.T, msg *FimpMessage, get func(*FimpMessage) (T, error), want T) {
	t.Helper()

	parsed := roundTrip(t, msg)

	for _, m := range []*FimpMessage{msg, parsed} {
		got, err := get(m)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}

	_, err := get(NewNullMessage("evt.test.report", "test", nil, nil, nil))
	assert.Error(t, err)
}

func roundTrip(t *testing.T, msg *FimpMessage) *FimpMessage {
	t.Helper()

	raw, err := msg.SerializeToJson()
	require.NoError(t, err)

	parsed, err := NewMessageFromBytes(raw)
	require.NoError(t, err)

	return parsed
}

func TestFimpMessage_TypedGettersRoundTrip(t *testing.T) {
	const i, s = "evt.test.report", fimptype.ServiceNameT("test")

	checkGetter(t, NewIntMessage(i, s, 7, nil, nil, nil), (*FimpMessage).GetIntValue, 7)
	checkGetter(t, NewStringMessage(i, s, "x", nil, nil, nil), (*FimpMessage).GetStringValue, "x")
	checkGetter(t, NewBoolMessage(i, s, true, nil, nil, nil), (*FimpMessage).GetBoolValue, true)
	checkGetter(t, NewFloatMessage(i, s, 1.5, nil, nil, nil), (*FimpMessage).GetFloatValue, 1.5)
	checkGetter(t, NewStrArrayMessage(i, s, []string{"a", "b"}, nil, nil, nil), (*FimpMessage).GetStrArrayValue, []string{"a", "b"})
	checkGetter(t, NewIntArrayMessage(i, s, []int{1, 2}, nil, nil, nil), (*FimpMessage).GetIntArrayValue, []int{1, 2})
	checkGetter(t, NewFloatArrayMessage(i, s, []float64{1.5, 2}, nil, nil, nil), (*FimpMessage).GetFloatArrayValue, []float64{1.5, 2})
	checkGetter(t, NewBoolArrayMessage(i, s, []bool{true, false}, nil, nil, nil), (*FimpMessage).GetBoolArrayValue, []bool{true, false})
	checkGetter(t, NewStrMapMessage(i, s, map[string]string{"k": "v"}, nil, nil, nil), (*FimpMessage).GetStrMapValue, map[string]string{"k": "v"})
	checkGetter(t, NewIntMapMessage(i, s, map[string]int{"k": 3}, nil, nil, nil), (*FimpMessage).GetIntMapValue, map[string]int{"k": 3})
	checkGetter(t, NewFloatMapMessage(i, s, map[string]float64{"k": 3.5}, nil, nil, nil), (*FimpMessage).GetFloatMapValue, map[string]float64{"k": 3.5})
	checkGetter(t, NewBoolMapMessage(i, s, map[string]bool{"k": true}, nil, nil, nil), (*FimpMessage).GetBoolMapValue, map[string]bool{"k": true})
}

func TestFimpMessage_GetIntValueVariants(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		want    int
		wantErr bool
	}{
		{"int64", int64(5), 5, false},
		{"int", 6, 6, false},
		{"uint", uint(7), 7, false},
		{"uint overflow", uint(math.MaxUint), 0, true},
		{"string", "1", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := (&FimpMessage{Value: tt.value}).GetIntValue()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestFimpMessage_Int64CollectionsConverted(t *testing.T) {
	arr, err := (&FimpMessage{Value: []int64{1, 2}}).GetIntArrayValue()
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2}, arr)

	m, err := (&FimpMessage{Value: map[string]int64{"a": 1}}).GetIntMapValue()
	require.NoError(t, err)
	assert.Equal(t, map[string]int{"a": 1}, m)
}

func TestFimpMessage_Binary(t *testing.T) {
	msg := NewBinaryMessage("evt.test.report", "test", []byte("hi"), nil, nil, nil)
	assert.Equal(t, fimptype.VTypeBinary, msg.ValueType)
	assert.Equal(t, "aGk=", roundTrip(t, msg).Value)
}

func TestFimpMessage_ObjectAndMisc(t *testing.T) {
	msg := NewObjectMessage("evt.test.report", "test", map[string]int{"a": 1}, nil, nil, nil)
	parsed := roundTrip(t, msg)
	assert.JSONEq(t, `{"a":1}`, string(parsed.GetRawObjectValue()))

	msg.SetValue(2.5, fimptype.VTypeFloat)
	assert.Equal(t, fimptype.VTypeFloat, msg.ValueType)
	assert.Equal(t, 2.5, msg.Value)

	now := time.Now()
	assert.WithinDuration(t, now, NewNullMessage("a", "b", nil, nil, nil).GetCreationTime(), time.Second)
	assert.True(t, (&FimpMessage{CreationTime: "bad"}).GetCreationTime().IsZero())
}

func TestFimpMessage_Str(t *testing.T) {
	tests := []struct {
		name string
		msg  FimpMessage
		want string
	}{
		{"no service", FimpMessage{Source: "src", Interface: "evt.a.b", Value: 1}, "src evt.a.b 1"},
		{"plain", FimpMessage{Source: "src", Service: "svc", Interface: "evt.a.b", Value: "v"}, "src svc evt.a.b v"},
		{"plain float", FimpMessage{Source: "src", Service: "svc", Interface: "evt.a.b", Value: 1.256}, "src svc evt.a.b 1.26"},
		{"storage", FimpMessage{Source: "src", Service: "svc", Interface: "evt.a.b", Value: "v", Storage: &Storage{Strategy: StorageStrategyAggregate, SubValue: "x"}}, "src svc evt.a.b v agg x"},
		{"storage float", FimpMessage{Source: "src", Service: "svc", Interface: "evt.a.b", Value: 1.5, Storage: &Storage{Strategy: StorageStrategySkip}}, "src svc evt.a.b 1.50 ski"},
		{"storage none", FimpMessage{Source: "src", Service: "svc", Interface: "evt.a.b", Value: 1, Storage: &Storage{Strategy: StorageStrategyNone}}, "src svc evt.a.b 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.msg.Str())
		})
	}
}

func TestProps_GetTimestampValue(t *testing.T) {
	p := Props{"ok": "2024-01-02T03:04:05Z", "bad": "nope"}

	ts, found, err := p.GetTimestampValue("ok")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, 2024, ts.Year())

	_, found, err = p.GetTimestampValue("bad")
	assert.True(t, found)
	assert.Error(t, err)

	_, found, err = p.GetTimestampValue("missing")
	assert.False(t, found)
	assert.NoError(t, err)
}

func TestProps_InvalidAndMissing(t *testing.T) {
	p := Props{"k": "x"}

	_, found, err := p.GetIntValue("k")
	assert.True(t, found)
	assert.Error(t, err)

	_, found, err = p.GetFloatValue("k")
	assert.True(t, found)
	assert.Error(t, err)

	_, found, err = p.GetBoolValue("k")
	assert.True(t, found)
	assert.Error(t, err)

	_, found, err = p.GetIntValue("none")
	assert.False(t, found)
	assert.NoError(t, err)

	_, found, err = p.GetFloatValue("none")
	assert.False(t, found)
	assert.NoError(t, err)

	_, found, err = p.GetBoolValue("none")
	assert.False(t, found)
	assert.NoError(t, err)

	_, found = p.GetStringValue("none")
	assert.False(t, found)
}

func TestNewMessageFromBytes_Errors(t *testing.T) {
	tests := map[string]string{
		"missing type":   `{"serv":"s","val_t":"string","val":"x"}`,
		"unknown val_t":  `{"type":"t","serv":"s","val_t":"weird","val":"x"}`,
		"wrong string":   `{"type":"t","serv":"s","val_t":"string","val":1}`,
		"bad int array":  `{"type":"t","serv":"s","val_t":"int_array","val":"x"}`,
		"bad str map":    `{"type":"t","serv":"s","val_t":"str_map","val":{"a":"\u12"}}`,
		"bad int map":    `{"type":"t","serv":"s","val_t":"int_map","val":{"a":"x"}}`,
		"bad float map":  `{"type":"t","serv":"s","val_t":"float_map","val":{"a":"x"}}`,
		"bad bool map":   `{"type":"t","serv":"s","val_t":"bool_map","val":{"a":"x"}}`,
		"bad props":      `{"type":"t","serv":"s","val_t":"null","props":[1]}`,
		"bad storage":    `{"type":"t","serv":"s","val_t":"null","storage":[1]}`,
		"bad tags":       `{"type":"t","serv":"s","val_t":"null","tags":{"a":1}}`,
		"missing object": `{"type":"t","serv":"s","val_t":"object"}`,
	}

	for name, raw := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := NewMessageFromBytes([]byte(raw))
			assert.Error(t, err)
			assert.Nil(t, msg)
		})
	}
}

func TestNewMessageFromBytes_SkipsInvalidArrayItems(t *testing.T) {
	tests := map[string]struct {
		valT, val string
		want      any
	}{
		"bool":  {"bool_array", `[true,"x"]`, []bool{true}},
		"str":   {"str_array", `["a","\u12"]`, []string{"a"}},
		"int":   {"int_array", `[1,"x"]`, []int{1}},
		"float": {"float_array", `[1.5,"x"]`, []float64{1.5}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := NewMessageFromBytes([]byte(`{"type":"t","serv":"s","val_t":"` + tt.valT + `","val":` + tt.val + `}`))
			require.NoError(t, err)
			assert.Equal(t, tt.want, msg.Value)
		})
	}
}

func TestNewMessageFromBytes_PropsStorageTagsAndMissingServ(t *testing.T) {
	msg, err := NewMessageFromBytes([]byte(`{"type":"t","val_t":"null","props":{"a":"b"},"storage":{"strategy":"skip"},"tags":["x"],"resp_to":"r","src":"me","topic":"tp","ver":"1"}`))
	require.NoError(t, err)
	assert.Equal(t, Props{"a": "b"}, msg.Properties)
	assert.Equal(t, StorageStrategySkip, msg.Storage.Strategy)
	assert.Equal(t, Tags{"x"}, msg.Tags)
	assert.Equal(t, "r", msg.ResponseToTopic)
	assert.Equal(t, fimptype.ResourceNameT("me"), msg.Source)
	assert.Empty(t, msg.Service)
}

func TestFimpMessage_WithHelpers(t *testing.T) {
	msg := NewNullMessage("cmd.a.b", "s", nil, nil, nil).WithProperty("k", "v").WithTag("t").WithStorageStrategy(StorageStrategySplit, "sv")
	assert.Equal(t, Props{"k": "v"}, msg.Properties)
	assert.Equal(t, Tags{"t"}, msg.Tags)
	assert.Equal(t, &Storage{Strategy: StorageStrategySplit, SubValue: "sv"}, msg.Storage)
}
