package fimpgo

import (
	"context"
	"fmt"
	"io"
	"net"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/futurehomeno/fimpgo/utils"
)

type fakeMessage struct {
	topic   string
	payload []byte
}

func (m fakeMessage) Duplicate() bool   { return false }
func (m fakeMessage) Qos() byte         { return 0 }
func (m fakeMessage) Retained() bool    { return false }
func (m fakeMessage) Topic() string     { return m.topic }
func (m fakeMessage) MessageID() uint16 { return 0 }
func (m fakeMessage) Payload() []byte   { return m.payload }
func (m fakeMessage) Ack()              {}

func TestMqttTransport_StartWaitsOnPendingConnect(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	mqtt := NewMqttTransport("tcp://127.0.0.1:1", "test_pending_connect", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(2)

	err := mqtt.Start(500 * time.Millisecond)

	assert.ErrorIs(t, err, utils.ErrTimeout)
	assert.Eventually(t, func() bool { return !mqtt.client.IsConnected() }, time.Second, 10*time.Millisecond)

	for _, e := range hook.AllEntries() {
		assert.NotContains(t, e.Message, "status can only transition")
	}
}

func TestMqttTransport_StartReturnsWhenPendingConnectSucceedsDuringBackoff(t *testing.T) {
	ctx := t.Context()

	var lc net.ListenConfig

	ln, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	context.AfterFunc(ctx, func() { _ = ln.Close() })

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_backoff_connect", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(3)

	defer mqtt.Stop()

	// The kernel completes the TCP handshake, so the first attempt stays pending until the proxy starts accepting.
	go func() {
		time.Sleep(600 * time.Millisecond)

		var d net.Dialer

		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}

			b, err := d.DialContext(ctx, "tcp", "127.0.0.1:11883")
			if err != nil {
				_ = c.Close()

				continue
			}

			go func() { _, _ = io.Copy(b, c) }()
			go func() { _, _ = io.Copy(c, b) }()
		}
	}()

	start := time.Now()

	require.NoError(t, mqtt.Start(300*time.Millisecond))
	assert.Less(t, time.Since(start), 2*time.Second)
}

func TestMqttTransport_StartReconnectsWhenPendingConnectFailsDuringBackoff(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	ctx := t.Context()

	var lc net.ListenConfig

	ln, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	context.AfterFunc(ctx, func() { _ = ln.Close() })

	var accepts atomic.Int32

	// Never answering CONNACK makes each attempt fail after paho's connect timeout, inside the backoff wait.
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}

			accepts.Add(1)
			context.AfterFunc(ctx, func() { _ = c.Close() })
		}
	}()

	opts := MQTT.NewClientOptions().AddBroker("tcp://" + ln.Addr().String()).SetConnectTimeout(time.Second).SetProtocolVersion(4)
	mqtt := NewMqttTransportFromConnection(MQTT.NewClient(opts), 1, 1)
	mqtt.SetStartAutoRetryCount(2)

	start := time.Now()

	assert.ErrorIs(t, mqtt.Start(300*time.Millisecond), utils.ErrTimeout)
	assert.EqualValues(t, 2, accepts.Load())
	// Reconnecting right after the failure would race paho's status reset and skip the backoff.
	assert.GreaterOrEqual(t, time.Since(start), 4*time.Second)
	assert.True(t, slices.ContainsFunc(hook.AllEntries(), func(e *log.Entry) bool {
		return strings.Contains(e.Message, "pending connect failed")
	}))
}

func TestMqttTransport_LogsKeysOfUnparsablePayload(t *testing.T) {
	payload := []byte(`{"serv":"meter_elec","type":"evt.meter.report","val":12.5}`)

	compressed, err := NewMsgCompressor("", "").CompressBinMsg(payload)
	require.NoError(t, err)

	bounded := `{"a\nb":0,"` + strings.Repeat("k", 100) + `":0`
	for i := 2; i < 20; i++ {
		bounded += fmt.Sprintf(`,"k%d":0`, i)
	}

	boundedKeys := `"a\nb" "` + strings.Repeat("k", 64) + `"`
	for i := 2; i < 16; i++ {
		boundedKeys += fmt.Sprintf(` "k%d"`, i)
	}

	tcs := []struct {
		name    string
		topic   string
		payload []byte
		keys    string
	}{
		{"plain", "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:6_0", payload, `keys=["serv" "type" "val"]`},
		{"compressed", "pt:j1c1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:6_0", compressed, `keys=["serv" "type" "val"]`},
		{"escaped and bounded", "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:6_0", []byte(bounded + "}"), "keys=[" + boundedKeys + ` "..."]`},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			hook := logtest.NewGlobal()
			defer hook.Reset()

			mqtt := NewMqttTransport("tcp://127.0.0.1:11883", "test_unparsable", "", "", true, 1, 1, nil)
			mqtt.handleIncomingMessage(fakeMessage{topic: tc.topic, payload: tc.payload})

			var msg string

			for _, e := range hook.AllEntries() {
				if e.Level == log.ErrorLevel {
					msg = e.Message
				}
			}

			require.NotEmpty(t, msg)
			assert.Contains(t, msg, "topic="+tc.topic)
			assert.Contains(t, msg, fmt.Sprintf("len=%d", len(tc.payload)))
			assert.Contains(t, msg, tc.keys)
			assert.NotContains(t, msg, "\n")
		})
	}
}
