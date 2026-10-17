package transport

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/futurehomeno/fimpgo"
)

func decodeBatch(t *testing.T, b []byte) []fimpgo.FimpMessage {
	t.Helper()
	var msgs []fimpgo.FimpMessage
	require.NoError(t, json.Unmarshal(b, &msgs))
	return msgs
}

func TestNewBufferedStream_InvalidArgs(t *testing.T) {
	assert.Nil(t, NewBufferedStream(0, time.Second, false))
	assert.Nil(t, NewBufferedStream(1, 0, false))
}

func TestBufferedStream_FlushOnSizeToChannel(t *testing.T) {
	su := NewBufferedStream(2, time.Hour, false)
	defer su.Close()
	sink := su.ConfigureChanelSink(2)
	assert.Equal(t, sink, su.SinkChannel())

	su.EnqueueMessage("pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:temp_sensor/ad:2", fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", 1, nil, nil, nil))
	assert.Equal(t, 1, su.Size())
	su.EnqueueMessage("pt:j1/mt:cmd/rt:dev/rn:zw/ad:1/sv:out_bin_switch/ad:2", fimpgo.NewBoolMessage("cmd.binary.set", "out_bin_switch", true, nil, nil, nil))
	assert.Equal(t, 0, su.Size())

	msgs := decodeBatch(t, <-sink)
	require.Len(t, msgs, 2)
	assert.Equal(t, "/rt:dev/rn:zw/ad:1/sv:temp_sensor/ad:2", msgs[0].Topic)
	assert.Equal(t, "/rt:dev/rn:zw/ad:1/sv:out_bin_switch/ad:2", msgs[1].Topic)
}

func TestBufferedStream_FlushOnInterval(t *testing.T) {
	su := NewBufferedStream(100, 20*time.Millisecond, false)
	defer su.Close()
	sink := su.ConfigureChanelSink(1)
	su.EnqueueMessage("t", fimpgo.NewStringMessage("evt.a.b", "s", "v", nil, nil, nil))

	select {
	case b := <-sink:
		assert.Len(t, decodeBatch(t, b), 1)
	case <-time.After(2 * time.Second):
		t.Fatal("interval flush did not happen")
	}
}

func TestBufferedStream_FileSink(t *testing.T) {
	for _, compress := range []bool{false, true} {
		dir := t.TempDir()
		su := NewBufferedStream(100, time.Hour, compress)
		su.ConfigureFileSink("pfx", dir)
		su.EnqueueMessage("t", fimpgo.NewStringMessage("evt.a.b", "s", "v", nil, nil, nil))
		su.FlushBuffer()
		su.FlushBuffer()
		su.Close()
		su.Close()

		ext := "json"
		if compress {
			ext = "gz"
		}
		files, err := filepath.Glob(filepath.Join(dir, "pfx_*."+ext))
		require.NoError(t, err)
		require.Len(t, files, 1)
		body, err := os.ReadFile(files[0])
		require.NoError(t, err)
		if compress {
			body, err = fimpgo.NewMsgCompressor("", "").DecompressBinMsg(body)
			require.NoError(t, err)
		}
		assert.Len(t, decodeBatch(t, body), 1)
	}
}

func TestBufferedStream_FileSinkError(t *testing.T) {
	su := NewBufferedStream(100, time.Hour, false)
	defer su.Close()
	su.ConfigureFileSink("pfx", filepath.Join(t.TempDir(), "missing"))
	su.EnqueueMessage("t", fimpgo.NewStringMessage("evt.a.b", "s", "v", nil, nil, nil))
	su.FlushBuffer()
	assert.Equal(t, 0, su.Size())
}

func TestBufferedStream_ObjectValueAndSource(t *testing.T) {
	su := NewBufferedStream(1, time.Hour, false)
	defer su.Close()
	sink := su.ConfigureChanelSink(1)

	src := make(fimpgo.MessageCh, 1)
	su.SetSourceStream(src)
	obj, err := fimpgo.NewMessageFromBytes([]byte(`{"type":"evt.o.report","serv":"s","val_t":"object","val":{"k":"v"}}`))
	require.NoError(t, err)
	src <- &fimpgo.Message{Topic: "pt:j1/mt:evt/x", Payload: obj}
	close(src)

	msgs := decodeBatch(t, <-sink)
	require.Len(t, msgs, 1)
	assert.Equal(t, "/x", msgs[0].Topic)
	assert.Equal(t, map[string]any{"k": "v"}, msgs[0].Value)
}
