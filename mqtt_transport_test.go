package fimpgo

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var msgChan = make(chan int)

func onMsg(topic string, addr *Address, iotMsg *FimpMessage, rawMessage []byte) {
	log.Infof("New msg %s val=%v", topic, iotMsg.Value)
	if addr.ServiceName == "temp_sensor" && addr.ServiceAddress == "300" {
		msgChan <- 1
	} else {
		msgChan <- 2
	}
}

func TestMqttTransport_Publish(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)
	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	result := <-msgChan

	if result != 1 {
		t.Error("Wrong message result=", result)
	}

	mqtt.Stop()
}

func TestMqttTransport_PublishStopPublish(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)
	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	result := <-msgChan
	if result != 1 {
		t.Errorf("Wrong message result=%d", result)
	}

	mqtt.Stop()

	mqtt = NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err = mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.Stop()
}

func TestMqttTransport_PublishSync(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	var cnt atomic.Int64

	mqtt.SetMessageHandler(func(topic string, addr *Address, iotMsg *FimpMessage, rawMessage []byte) {
		if addr.ServiceName != "temp_sensor" {
			return
		}

		val, err := iotMsg.GetIntValue()
		if err != nil {
			log.Warnf("Wrong value %t err: %v", iotMsg.Value, err)
		} else {
			cnt.Add(val)
		}
	})

	if err := mqtt.Subscribe("pt:j1/mt:evt/#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewIntMessage("evt.sensor.report", "temp_sensor", int64(35), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}

	expVal := int64(0)
	for range 10 {
		msg.Value = int64(rand.Intn(100))
		expVal += msg.Value.(int64)

		err = mqtt.PublishSync(&adr, msg)
		if err != nil {
			log.Info("Publish failed err:", err)
		}
	}

	time.Sleep(200 * time.Millisecond)
	mqtt.Stop()

	assert.Equal(t, expVal, cnt.Load())
}

func TestMqttTransport_SubUnsub(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)
	if err := mqtt.Subscribe("pt:j1/mt:evt/#"); err != nil {
		t.Fatal("Subscribe err:", err)
		return
	}

	// unsubscribe and send message, shall not receive it
	err = mqtt.Unsubscribe("pt:j1/mt:evt/#")
	if err != nil {
		t.Fatal("Unsubscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.PublishSync(&adr, msg)
	if err != nil {
		t.Error("PublishSync err:", err)
		t.FailNow()
	}

	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	select {
	case <-msgChan:
		t.Error("Should not receive msg")
	case <-time.After(2 * time.Second):
	}

	mqtt.Stop()
}

// TODO: Fix, awsiot.private.key is not available in the repo
func TestMqttTransport_PublishTls(t *testing.T) {
	t.Skip()
	log.SetLevel(log.DebugLevel)
	// for test replace XYZ with actual AWS IoT core address and ABC with actual clientid
	mqtt := NewMqttTransport("ssl://a1ds8ixdqbiw53-ats.iot.eu-central-1.amazonaws.com:443", "00000000alexdevtest", "", "", false, 1, 1, nil)

	// for test enter valid site-id
	mqtt.SetGlobalTopicPrefix("331D092F-4685-4CC9-8337-2598E6F5D8D5")
	// for test place certificate and key into certs folder
	err := mqtt.ConfigureTls("awsiot.private.key", "awsiot.crt", "./certs", true)

	if err != nil {
		t.Fatal("Configure TLS err", err)
	}

	err = mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	result := <-msgChan

	if result != 1 {
		t.Error("Wrong message")
	}

	mqtt.Stop()
}

// TODO: Fix, awsiot.private.key is not available in the repo
func TestMqttTransport_PublishTls_2(t *testing.T) {
	t.Skip()
	connConfig := MqttConnectionConfigs{
		ServerURI:          "ssl://a1ds8ixdqbiw53-ats.iot.eu-central-1.amazonaws.com:443",
		ClientID:           "00000000alexdevtest",
		CleanSession:       true,
		SubQos:             1,
		PubQos:             1,
		CertDir:            "./certs",
		PrivateKeyFileName: "awsiot.private.key",
		CertFileName:       "awsiot.crt",
	}

	log.SetLevel(log.DebugLevel)
	// for test replace XYZ with actual AWS IoT core address and ABC with actual clientid
	mqtt := NewMqttTransportFromConfigs(connConfig)

	// for test enter valid site-id
	mqtt.SetGlobalTopicPrefix("331D092F-4685-4CC9-8337-2598E6F5D8D5")
	// for test place certificate and key into certs folder
	err := mqtt.ConfigureTls("awsiot.private.key", "awsiot.crt", "./certs", true)

	if err != nil {
		t.Fatal("Configure TLS err:", err)
	}

	err = mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	result := <-msgChan
	if result != 1 {
		t.Error("Wrong message")
	}

	mqtt.Stop()
}

func TestMqttTransport_TestChannels(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	mqtt.RegisterChannel("chan1", chan1)
	mqtt.RegisterChannel("chan2", chan2)
	correctMsg := make(chan int, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(msgChan MessageCh) {
		wg.Done()
		newMsg := <-chan1
		if newMsg.Payload.Service == "temp_sensor" {
			correctMsg <- 1
		}
	}(chan1)

	go func(msgChan MessageCh) {
		wg.Done()
		newMsg := <-chan2
		if newMsg.Payload.Service == "temp_sensor" {
			correctMsg <- 2
		}
	}(chan2)

	wg.Wait()

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Generate err:", err)
	}

	expVals := map[int]bool{1: true, 2: true}

	for len(expVals) > 0 {
		select {
		case val := <-correctMsg:
			delete(expVals, val)

		case <-time.After(2 * time.Second):
			t.Fatal("Message not received within timeout missing:", expVals)
		}
	}

	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	mqtt.Stop()
}

func TestMqttTransport_TestResponder(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest-1", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	mqtt2 := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest-2", "", "", true, 1, 1, nil)
	err = mqtt2.Start()

	if err != nil {
		t.Fatal("Start MQTT 2 err:", err)
	}

	if err := mqtt2.Subscribe("pt:j1c1/mt:rsp/rt:app/rn:response_tester/ad:1"); err != nil {
		t.Fatal("Subscribe response_tester err:", err)
	}

	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	mqtt.RegisterChannel("chan1", chan1)
	mqtt2.RegisterChannel("chan2", chan2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(msgChan MessageCh) {
		wg.Done()
		newMsg := <-chan1
		if newMsg.Payload.Service == "tester" {
			if err := mqtt.RespondToRequest(newMsg.Payload, NewFloatMessage("evt.test.response", "test_responder", 35.5, nil, nil, nil)); err != nil {
				t.Error("Error responding to request:", err)
			}
			return
		}
	}(chan1)

	rspReceived := make(chan struct{})

	go func(msgChan MessageCh) {
		wg.Done()
		newMsg := <-chan2
		if newMsg.Payload.Service == "test_responder" && newMsg.Topic == "pt:j1c1/mt:rsp/rt:app/rn:response_tester/ad:1" {
			close(rspReceived)
			return
		} else {
			t.Error("Wrong response message received :", newMsg)
		}
	}(chan2)

	wg.Wait()

	msg := NewFloatMessage("cmd.test.get_response", "tester", float64(35.5), nil, nil, nil)
	msg.ResponseToTopic = "pt:j1c1/mt:rsp/rt:app/rn:response_tester/ad:1"
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "test", ResourceAddress: "1"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	select {
	case <-rspReceived:

	case <-time.After(3 * time.Second):
		t.Error("Response not received within timeout")
		t.Fail()
	}

	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	err = mqtt.Unsubscribe("#")
	if err != nil {
		t.Fatal("AddSubscription err:", err)
	}
	mqtt.Stop()
}

func TestMqttTransport_TestChannelsWithFilters(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1, nil)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
		return
	}

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
		return
	}

	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	chan3 := make(MessageCh)
	chan4 := make(MessageCh)
	chan5 := make(MessageCh)
	mqtt.RegisterChannel("chan1", chan1)
	mqtt.RegisterChannel("chan2", chan2)
	mqtt.RegisterChannelWithFilter("chan3", chan3, FimpFilter{
		Topic:     "pt:j1/mt:evt/#",
		Service:   "temp_sensor",
		Interface: "evt.sensor.report",
	})

	mqtt.RegisterChannelWithFilter("chan4", chan4, FimpFilter{
		Topic:     "pt:j1/mt:cmd/#",
		Service:   "temp_sensor",
		Interface: "cmd.sensor.report",
	})

	testFilterFunc := func(topic string, addr *Address, iotMsg *FimpMessage) bool {
		return iotMsg.Type == "evt.sensor.report"
	}

	mqtt.RegisterChannelWithFilterFunc("chan5", chan5, testFilterFunc)

	var startedWg sync.WaitGroup
	startedWg.Add(5)
	correctMsg := make(chan int, 2)

	go func(msgChan MessageCh) {
		startedWg.Done()
		newMsg := <-msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			correctMsg <- 1
		}
	}(chan1)

	go func(msgChan MessageCh) {
		startedWg.Done()
		newMsg := <-msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			correctMsg <- 2
		}
	}(chan2)

	go func(msgChan MessageCh) {
		startedWg.Done()
		newMsg := <-msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			correctMsg <- 3
		}
	}(chan3)

	go func(msgChan MessageCh) {
		startedWg.Done()
		<-msgChan
		correctMsg <- 4
	}(chan4)

	go func(msgChan MessageCh) {
		startedWg.Done()
		<-msgChan
		correctMsg <- 5
	}(chan5)

	startedWg.Wait()

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	err = mqtt.Publish(&adr, msg)
	if err != nil {
		t.Fatal("Publish err:", err)
	}

	expVals := map[int]bool{1: true, 2: true, 3: true, 5: true}

	for len(expVals) > 0 {
		select {
		case val := <-correctMsg:
			if val == 4 {
				t.Error("Should not receive msg on chan4")
				t.Fail()
				return
			}

			delete(expVals, val)

		case <-time.After(2 * time.Second):
			t.Fatal("Message not received within timeout missing:", expVals)
		}
	}

	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	mqtt.UnregisterChannel("chan3")
	mqtt.UnregisterChannel("chan4")
	mqtt.UnregisterChannel("chan5")
	mqtt.Stop()

}

func TestAddGlobalPrefixToTopic(t *testing.T) {
	result := AddGlobalPrefixToTopic("12345", "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if result != "12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}

	result = AddGlobalPrefixToTopic("12345", "/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if result != "12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}

	result = AddGlobalPrefixToTopic("", "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if result != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}
}

func TestDetachGlobalPrefixFromTopic(t *testing.T) {
	globalPrefix, topic := DetachGlobalPrefixFromTopic("12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if globalPrefix != "12345" || topic != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}

	globalPrefix, topic = DetachGlobalPrefixFromTopic("ABC/12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if globalPrefix != "ABC/12345" || topic != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}
}
