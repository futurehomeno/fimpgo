package fimpgo

import (
	"errors"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/futurehomeno/fimpgo/fimptype"
)

const covBroker = "tcp://127.0.0.1:11883"

func startCovTransport(t *testing.T) *MqttTransport {
	t.Helper()

	mqtt := NewMqttTransport(covBroker, "cov_"+uuid.NewString()[:8], "", "", true, 1, 1, nil)
	require.NoError(t, mqtt.Start(5*time.Second))
	t.Cleanup(mqtt.Stop)

	return mqtt
}

func recvMsg(t *testing.T, ch MessageCh) *Message {
	t.Helper()

	select {
	case m := <-ch:
		return m
	case <-time.After(5 * time.Second):
		t.Fatal("message not received")
		return nil
	}
}

func TestMqttTransport_PublishVariants(t *testing.T) {
	mqtt := startCovTransport(t)
	mqtt.SetGlobalTopicPrefix(" site1 ")
	mqtt.SetDefaultSource("defsrc")

	name := fimptype.ResourceNameT("cov" + uuid.NewString()[:8])
	topic := func(pt string) string { return "pt:" + pt + "/mt:evt/rt:app/rn:" + name.Str() + "/ad:1" }
	require.NoError(t, mqtt.Subscribe(topic("j1")))
	require.NoError(t, mqtt.Subscribe(topic("j1c1")))

	ch := make(MessageCh, 10)
	mqtt.RegisterChannel("cov", ch)

	addr := func(pt string) *Address {
		return &Address{PayloadType: pt, MsgType: fimptype.MsgTypeEvt, ResourceType: fimptype.ResourceTypeApp, ResourceName: name, ResourceAddress: "1"}
	}
	newMsg := func() *FimpMessage { return NewStringMessage("evt.cov.report", "cov", "v", nil, nil, nil) }
	raw := func() []byte {
		b, err := newMsg().SerializeToJson()
		require.NoError(t, err)
		return b
	}

	tests := []struct {
		name string
		send func() error
		pt   string
		src  fimptype.ResourceNameT
	}{
		{"Publish", func() error { return mqtt.Publish(addr(""), newMsg()) }, "j1", "defsrc"},
		{"Publish compressed", func() error { return mqtt.Publish(addr(CompressedJsonPayload), newMsg()) }, "j1c1", "defsrc"},
		{"PublishSync", func() error { return mqtt.PublishSync(addr(DefaultPayload), newMsg()) }, "j1", "defsrc"},
		{"PublishSync compressed", func() error { return mqtt.PublishSync(addr(CompressedJsonPayload), newMsg()) }, "j1c1", "defsrc"},
		{"PublishSync unknown payload falls back", func() error { return mqtt.PublishSync(addr("zz"), newMsg()) }, "j1", "defsrc"},
		{"PublishToTopic", func() error { return mqtt.PublishToTopic(topic("j1"), newMsg()) }, "j1", "defsrc"},
		{"PublishToTopic compressed", func() error { return mqtt.PublishToTopic(topic("j1c1"), newMsg()) }, "j1c1", "defsrc"},
		{"RespondToRequest", func() error {
			req := newMsg()
			req.ResponseToTopic = topic("j1")
			return mqtt.RespondToRequest(req, newMsg())
		}, "j1", "defsrc"},
		{"PublishRaw", func() error { mqtt.PublishRaw("site1/"+topic("j1"), raw()); return nil }, "j1", ""},
		{"PublishRawSync", func() error { return mqtt.PublishRawSync("site1/"+topic("j1"), raw()) }, "j1", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.send())

			m := recvMsg(t, ch)
			assert.Equal(t, tt.pt, m.Addr.PayloadType)
			assert.Equal(t, tt.src, m.Payload.Source)
			assert.Equal(t, "evt.cov.report", m.Payload.Interface)
			assert.Equal(t, topic(tt.pt), m.Topic)
		})
	}

	require.NoError(t, mqtt.Unsubscribe(topic("j1")))
	require.NoError(t, mqtt.UnsubscribeAll())
}

func TestMqttTransport_PublishErrors(t *testing.T) {
	mqtt := NewMqttTransport(covBroker, "cov_unused", "", "", true, 1, 1, nil)
	addr := &Address{PayloadType: "zz", MsgType: fimptype.MsgTypeEvt, ResourceType: fimptype.ResourceTypeApp, ResourceName: "n", ResourceAddress: "1"}
	bad := NewObjectMessage("evt.a.b", "s", make(chan int), nil, nil, nil)

	assert.ErrorContains(t, mqtt.Publish(addr, NewNullMessage("evt.a.b", "s", nil, nil, nil)), "unknown payload type")
	assert.Error(t, mqtt.Publish(&Address{}, bad))
	assert.Error(t, mqtt.PublishSync(addr, bad))
	assert.Error(t, mqtt.PublishToTopic("pt:j1/mt:evt/rt:app/rn:n/ad:1", bad))
	assert.ErrorContains(t, mqtt.RespondToRequest(NewNullMessage("a", "b", nil, nil, nil), bad), "empty response topic")
	assert.Error(t, mqtt.PublishRawSync("t", nil))
}

func TestMqttTransport_OfflineSubscribeErrors(t *testing.T) {
	mqtt := NewMqttTransport(covBroker, "cov_offline", "", "", true, 1, 1, nil)

	assert.NoError(t, mqtt.Subscribe(" "))
	assert.Error(t, mqtt.Subscribe("a/b"))
	assert.Error(t, mqtt.Unsubscribe("a/b"))

	mqtt.subs["a/b"] = 1
	assert.Error(t, mqtt.UnsubscribeAll())
}

func TestMqttTransport_ConnectionLostHandler(t *testing.T) {
	mqtt := NewMqttTransport(covBroker, "cov_lost", "", "", true, 1, 1, nil)
	mqtt.connState.Init()
	mqtt.connState.OnConnect()
	require.True(t, mqtt.IsConnected())

	mqtt.onConnectionLost(mqtt.client, errors.New("boom"))
	assert.False(t, mqtt.IsConnected())

	var got error

	mqtt.SetOnConnectionLostHandler(func(_ MQTT.Client, err error) { got = err })
	mqtt.onConnectionLost(mqtt.client, errors.New("boom2"))
	assert.EqualError(t, got, "boom2")
}

func TestMqttTransport_QueueOverflowStopsTransport(t *testing.T) {
	var reported error

	mqtt := NewMqttTransport(covBroker, "cov_overflow", "", "", true, 1, 1, func(err error) { reported = err })
	mqtt.mainQueue = make(chan MQTT.Message)
	mqtt.connState.Init()

	msg := fakeMessage{topic: "t"}

	// onMessage stops the transport after more than 20 consecutive overflows, but only while connected.
	for range 22 {
		mqtt.onMessage(nil, msg)
	}

	assert.NoError(t, reported)

	mqtt.connState.OnConnect()
	mqtt.onMessage(nil, msg)

	assert.EqualError(t, reported, "main msg queue stuck")
	assert.False(t, mqtt.IsConnected())
}

func TestMqttTransport_HandleIncomingMessage(t *testing.T) {
	const (
		plainTopic = "pt:j1/mt:evt/rt:app/rn:n/ad:1"
		compTopic  = "pt:j1c1/mt:evt/rt:app/rn:n/ad:1"
	)

	payload, err := NewStringMessage("evt.a.b", "svc", "v", nil, nil, nil).SerializeToJson()
	require.NoError(t, err)

	compressed, err := NewMsgCompressor("", "").CompressBinMsg(payload)
	require.NoError(t, err)

	tests := []struct {
		name   string
		prefix string
		topic  string
		data   []byte
		noComp bool
		want   bool
	}{
		{"plain", "", plainTopic, payload, false, true},
		{"prefixed", "site", "site/" + plainTopic, payload, false, true},
		{"compressed", "", compTopic, compressed, false, true},
		{"compressed without compressor", "", compTopic, compressed, true, false},
		{"corrupt compressed", "", compTopic, []byte("junk"), false, false},
		{"unknown payload type", "", "pt:zz/mt:evt/rt:app/rn:n/ad:1", payload, false, false},
		{"bad topic", "", "pt:j1/bad:x:y", payload, false, false},
		{"prefix without payload type", "site", "site/nothing", payload, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mqtt := NewMqttTransport(covBroker, "cov_incoming", "", "", true, 1, 1, nil)
			mqtt.SetGlobalTopicPrefix(tt.prefix)

			if tt.noComp {
				mqtt.compressor = nil
			}

			handled := 0
			mqtt.SetMessageHandler(func(string, *Address, *FimpMessage, []byte) { handled++ })
			mqtt.handleIncomingMessage(fakeMessage{topic: tt.topic, payload: tt.data})

			assert.Equal(t, tt.want, handled == 1)
		})
	}
}

func TestMqttTransport_ChannelFiltering(t *testing.T) {
	mqtt := NewMqttTransport(covBroker, "cov_filters", "", "", true, 1, 1, nil)
	payload, err := NewStringMessage("evt.a.b", "svc", "v", nil, nil, nil).SerializeToJson()
	require.NoError(t, err)

	all, byFilter, byFunc, blocked, panicky := make(MessageCh, 1), make(MessageCh, 1), make(MessageCh, 1), make(MessageCh, 1), make(MessageCh, 1)

	mqtt.RegisterChannel("all", all)
	mqtt.RegisterChannelWithFilter("filter", byFilter, FimpFilter{Topic: "pt:j1/mt:evt/#", Service: "*", Interface: "evt.a.b"})
	mqtt.RegisterChannelWithFilterFunc("func", byFunc, func(string, *Address, *FimpMessage) bool { return true })
	mqtt.RegisterChannelWithFilter("blocked", blocked, FimpFilter{Topic: "pt:j1/mt:evt/#", Service: "other", Interface: "*"})
	mqtt.RegisterChannelWithFilterFunc("panicky", panicky, func(string, *Address, *FimpMessage) bool { panic("boom") })

	mqtt.handleIncomingMessage(fakeMessage{topic: "pt:j1/mt:evt/rt:app/rn:n/ad:1", payload: payload})

	assert.Len(t, all, 1)
	assert.Len(t, byFilter, 1)
	assert.Len(t, byFunc, 1)
	assert.Empty(t, blocked)
	assert.Empty(t, panicky)

	mqtt.UnregisterChannel("filter")
	assert.NotContains(t, mqtt.subFilters, "filter")
	assert.NotContains(t, mqtt.subChannels, "filter")

	slow := NewMqttTransport(covBroker, "cov_slow", "", "", true, 1, 1, nil)
	slow.RegisterChannel("slow", make(MessageCh))
	slow.SetReceiveChTimeout(0)
	slow.handleIncomingMessage(fakeMessage{topic: "pt:j1/mt:evt/rt:app/rn:n/ad:1", payload: payload})
}

func TestMqttTransport_IsChannelInterestedBinaryPayload(t *testing.T) {
	mqtt := NewMqttTransport(covBroker, "cov_binary", "", "", true, 1, 1, nil)
	mqtt.RegisterChannelWithFilter("c", make(MessageCh), FimpFilter{Topic: "pt:j1c1/#"})

	assert.True(t, mqtt.isChannelInterested("c", "pt:j1c1/mt:evt", nil, nil))
	assert.False(t, mqtt.isChannelInterested("c", "pt:j1/mt:evt", nil, nil))
}

func TestNewMqttTransportFromConfigs(t *testing.T) {
	mqtt := NewMqttTransportFromConfigs(MqttConnectionConfigs{
		ServerURI: covBroker, ClientID: "cov_cfg", GlobalTopicPrefix: " pfx ", StartFailRetryCount: 3, ReceiveChTimeout: 7, MainQueueSize: 5,
	}, nil)
	require.NotNil(t, mqtt)

	assert.Equal(t, "pfx", mqtt.globalTopicPrefix())
	assert.Equal(t, 3, mqtt.startFailRetryCount)
	assert.EqualValues(t, 7, mqtt.receiveChTimeout.Load())
	assert.Equal(t, 5, cap(mqtt.mainQueue))

	mqtt.SetStartAutoRetryCount(4)
	mqtt.SetCertDir("/x")
	assert.Equal(t, 4, mqtt.startFailRetryCount)
	assert.Equal(t, "/x", mqtt.certDir)

	assert.Nil(t, NewMqttTransportFromConfigs(MqttConnectionConfigs{
		ServerURI: covBroker, PrivateKeyFileName: "k", CertFileName: "c", CertDir: t.TempDir(),
	}, nil))
}
