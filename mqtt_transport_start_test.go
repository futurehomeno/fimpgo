package fimpgo

import (
	"context"
	"fmt"
	"io"
	"net"
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

	addr := ln.Addr().String()
	require.NoError(t, ln.Close())

	mqtt := NewMqttTransport("tcp://"+addr, "test_backoff_connect", "", "", true, 1, 1, nil)
	opts := defaultClientOptions("tcp://"+addr, "test_backoff_connect", "", "", true, mqtt.onConnectionLost)
	opts.SetDefaultPublishHandler(mqtt.onMessage)
	opts.SetOnConnectHandler(mqtt.onConnect)
	opts.SetConnectRetryInterval(100 * time.Millisecond)
	mqtt.client = MQTT.NewClient(opts)
	mqtt.SetStartAutoRetryCount(3)

	defer mqtt.Stop()

	go func() {
		time.Sleep(600 * time.Millisecond)

		proxy, err := lc.Listen(ctx, "tcp", addr)
		if err != nil {
			return
		}

		context.AfterFunc(ctx, func() { _ = proxy.Close() })

		var d net.Dialer

		for {
			c, err := proxy.Accept()
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

func TestMqttTransport_LogsKeysOfUnparsablePayload(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	topic := "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:6_0"
	payload := []byte(`{"serv":"meter_elec","type":"evt.meter.report","val":12.5}`)

	mqtt := NewMqttTransport("tcp://127.0.0.1:11883", "test_unparsable", "", "", true, 1, 1, nil)
	mqtt.handleIncomingMessage(fakeMessage{topic: topic, payload: payload})

	var msg string

	for _, e := range hook.AllEntries() {
		if e.Level == log.ErrorLevel {
			msg = e.Message
		}
	}

	require.NotEmpty(t, msg)
	assert.Contains(t, msg, "topic="+topic)
	assert.Contains(t, msg, fmt.Sprintf("len=%d", len(payload)))
	assert.Contains(t, msg, "keys=[serv type val]")
}
