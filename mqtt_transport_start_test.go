package fimpgo

import (
	"context"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/eclipse/paho.mqtt.golang/packets"
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

func listen(t *testing.T, addr string) net.Listener {
	t.Helper()

	var lc net.ListenConfig

	ln, err := lc.Listen(t.Context(), "tcp", addr)
	require.NoError(t, err)

	t.Cleanup(func() { _ = ln.Close() })

	return ln
}

func serveBroker(ln net.Listener, subackDelay time.Duration) {
	for {
		c, err := ln.Accept()
		if err != nil {
			return
		}

		go func() {
			defer func() { _ = c.Close() }()

			for {
				p, err := packets.ReadPacket(c)
				if err != nil {
					return
				}

				switch p := p.(type) {
				case *packets.ConnectPacket:
					_ = packets.NewControlPacket(packets.Connack).Write(c)
				case *packets.SubscribePacket:
					time.Sleep(subackDelay)

					ack := packets.SubackPacket{FixedHeader: packets.FixedHeader{MessageType: packets.Suback}, MessageID: p.MessageID, ReturnCodes: p.Qoss}
					_ = ack.Write(c)
				}
			}
		}()
	}
}

func TestMqttTransport_StartWaitsOnPendingConnect(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	mqtt := NewMqttTransport("tcp://127.0.0.1:1", "test_pending_connect", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(2)

	assert.ErrorIs(t, mqtt.Start(100*time.Millisecond), utils.ErrTimeout)
	assert.False(t, mqtt.IsConnected())
	assert.Eventually(t, func() bool { return !mqtt.client.IsConnected() }, time.Second, 10*time.Millisecond)

	for _, e := range hook.AllEntries() {
		assert.NotContains(t, e.Message, "status can only transition")
	}
}

func TestMqttTransport_StartReturnsWhenPendingConnectSucceeds(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_late_connect", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(3)

	defer mqtt.Stop()

	// The kernel completes the TCP handshake, so the attempt stays pending until the broker starts accepting.
	time.AfterFunc(600*time.Millisecond, func() { serveBroker(ln, 0) })

	start := time.Now()

	require.NoError(t, mqtt.Start(300*time.Millisecond))
	assert.Less(t, time.Since(start), 2*time.Second)
}

func TestMqttTransport_StartReconnectsWithoutConnectRetry(t *testing.T) {
	hook := logtest.NewGlobal()
	defer hook.Reset()

	ctx := t.Context()
	ln := listen(t, "127.0.0.1:0")

	var accepts atomic.Int32

	// Never answering CONNACK makes each attempt fail after paho's connect timeout.
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

	assert.ErrorIs(t, mqtt.Start(1500*time.Millisecond), os.ErrDeadlineExceeded)
	assert.EqualValues(t, 2, accepts.Load())
	assert.True(t, slices.ContainsFunc(hook.AllEntries(), func(e *log.Entry) bool {
		return strings.Contains(e.Message, "connect attempt 1 failed")
	}))
}

func TestMqttTransport_StartReturnsLastConnectErrorWithinBudget(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}

			if _, err := packets.ReadPacket(c); err == nil {
				ack := packets.ConnackPacket{FixedHeader: packets.FixedHeader{MessageType: packets.Connack}, ReturnCode: packets.ErrRefusedNotAuthorised}
				_ = ack.Write(c)
			}

			_ = c.Close()
		}
	}()

	opts := MQTT.NewClientOptions().AddBroker("tcp://" + ln.Addr().String()).SetProtocolVersion(4)
	mqtt := NewMqttTransportFromConnection(MQTT.NewClient(opts), 1, 1)
	mqtt.SetStartAutoRetryCount(2)

	start := time.Now()

	// The budget is 2*300ms plus the 1s backoff between the two attempts.
	assert.ErrorIs(t, mqtt.Start(300*time.Millisecond), packets.ErrorRefusedNotAuthorised)
	assert.Less(t, time.Since(start), 3*time.Second)
}

func TestMqttTransport_StartWithOutOfRangeRetryCount(t *testing.T) {
	for _, count := range []int{0, 100000} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			ln := listen(t, "127.0.0.1:0")

			go serveBroker(ln, 0)

			mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_retry_count", "", "", true, 1, 1, nil)
			mqtt.SetStartAutoRetryCount(count)

			defer mqtt.Stop()

			require.NoError(t, mqtt.Start(time.Second))
		})
	}
}

func TestMqttTransport_StartConnectsAfterFailedStart(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")
	addr := ln.Addr().String()
	require.NoError(t, ln.Close())

	mqtt := NewMqttTransportFromConnection(nil, 1, 1)
	opts := defaultClientOptions("tcp://"+addr, "test_restart", "", "", true, mqtt.onConnectionLost)
	mqtt.client = MQTT.NewClient(opts.SetOnConnectHandler(mqtt.onConnect).SetConnectRetryInterval(500 * time.Millisecond))
	mqtt.SetStartAutoRetryCount(1)

	defer mqtt.Stop()

	require.ErrorIs(t, mqtt.Start(300*time.Millisecond), utils.ErrTimeout)

	go serveBroker(listen(t, addr), 0)

	require.NoError(t, mqtt.Start(2*time.Second))
}

func TestMqttTransport_StartWhileStarted(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go serveBroker(ln, 0)

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_start_twice", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(1)

	defer mqtt.Stop()

	require.NoError(t, mqtt.Start(time.Second))

	returned := make(chan error, 1)
	go func() { returned <- mqtt.Start(time.Second) }()

	select {
	case err := <-returned:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("second Start hung")
	}
}

// The check-then-Init window is microseconds wide, so the race is retried on fresh transports.
func TestMqttTransport_ConcurrentStart(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go serveBroker(ln, 0)

	for i := range 50 {
		mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), fmt.Sprintf("test_concurrent_start_%d", i), "", "", true, 1, 1, nil)
		mqtt.SetStartAutoRetryCount(1)

		gate, errs := make(chan struct{}), make(chan error, 20)
		for range cap(errs) {
			go func() {
				<-gate
				errs <- mqtt.Start(2 * time.Second)
			}()
		}

		close(gate)

		for range cap(errs) {
			require.NoError(t, <-errs)
		}

		requireReturns(t, mqtt.Stop, "Stop")
	}
}

func TestMqttTransport_StartWhileStopWaitsForReader(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go serveBroker(ln, 0)

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_start_during_stop", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(1)

	handling, release := make(chan struct{}), make(chan struct{})
	mqtt.SetMessageHandler(func(string, *Address, *FimpMessage, []byte) {
		close(handling)
		<-release
	})

	require.NoError(t, mqtt.Start(time.Second))

	mqtt.mainQueue <- fakeMessage{topic: "pt:j1/mt:evt/rt:app/rn:test/ad:1", payload: []byte(`{"serv":"test","type":"evt.test.report","val_t":"null"}`)}
	<-handling

	stopped := make(chan struct{})
	go func() {
		mqtt.Stop()
		close(stopped)
	}()

	require.Eventually(t, func() bool {
		select {
		case <-mqtt.connState.DoneC():
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

	started := make(chan error, 1)
	go func() { started <- mqtt.Start(time.Second) }()

	time.Sleep(200 * time.Millisecond)
	close(release)

	requireReturns(t, func() { <-stopped }, "Stop")
	var err error
	requireReturns(t, func() { err = <-started }, "Start")
	require.NoError(t, err)
	requireReturns(t, mqtt.Stop, "second Stop")
}

func TestMqttTransport_StartRightAfterStop(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go serveBroker(ln, 0)

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_restart_after_stop", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(1)

	for range 20 {
		require.NoError(t, mqtt.Start(time.Second))
		mqtt.Stop()
	}
}

func requireReturns(t *testing.T, f func(), name string) {
	t.Helper()

	returned := make(chan struct{})
	go func() {
		f()
		close(returned)
	}()

	select {
	case <-returned:
	case <-time.After(5 * time.Second):
		t.Fatal(name + " hung")
	}
}

func TestMqttTransport_FailedStartStaysDisconnected(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go serveBroker(ln, 600*time.Millisecond)

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_failed_start", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(1)
	mqtt.subs["pt:j1/mt:evt/rt:app/rn:test/ad:1"] = 1

	defer mqtt.Stop()

	require.ErrorIs(t, mqtt.Start(300*time.Millisecond), utils.ErrTimeout)
	assert.Never(t, mqtt.IsConnected, time.Second, 20*time.Millisecond)
}

func TestMqttTransport_StaleConnectDuringRestartStaysDisconnected(t *testing.T) {
	ln := listen(t, "127.0.0.1:0")

	go serveBroker(ln, 600*time.Millisecond)

	mqtt := NewMqttTransport("tcp://"+ln.Addr().String(), "test_stale_connect", "", "", true, 1, 1, nil)
	mqtt.SetStartAutoRetryCount(1)
	mqtt.subs["pt:j1/mt:evt/rt:app/rn:test/ad:1"] = 1

	defer mqtt.Stop()

	require.ErrorIs(t, mqtt.Start(300*time.Millisecond), utils.ErrTimeout)

	mqtt.connState.Init()
	assert.Never(t, mqtt.IsConnected, time.Second, 20*time.Millisecond)
}

func TestConnState_OnConnectAfterDoneStaysDisconnected(t *testing.T) {
	var c ConnStateT

	c.Init()
	c.OnDone()
	c.OnConnect(c.DoneC())

	assert.False(t, c.IsConnected())
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
		{"value of wrong type", "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:6_0", []byte(`{"type":"evt.meter.report","val_t":"float","val":"user@example.com"}`), `keys=["type" "val_t" "val"]`},
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
			assert.NotContains(t, msg, "user@example.com")
		})
	}
}

func TestPayloadKeys_CutsLongKeyBeforeCopy(t *testing.T) {
	keys := payloadKeys([]byte(`{"` + strings.Repeat("k", 100) + `":0}`))

	assert.Equal(t, []string{strings.Repeat("k", 64)}, keys)
}
