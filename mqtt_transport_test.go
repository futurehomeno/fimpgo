package fimpgo

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
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

var isCorrect = make(map[int]bool)

func TestMqttTransport_Publish(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
		return
	}

	t.Log("Connected")

	mqtt.SetMessageHandler(onMsg)
	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
		return
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)

	result := <-msgChan

	if result != 1 {
		t.Error("Wrong message")
	}

	mqtt.Stop()
}

func TestMqttTransport_PublishStopPublish(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)
	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
		return
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)

	result := <-msgChan
	if result != 1 {
		t.Error("Wrong message")
	}

	mqtt.Stop()
	time.Sleep(100 * time.Millisecond)

	mqtt = NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err = mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	time.Sleep(100 * time.Millisecond)
	mqtt.Stop()
}

func TestMqttTransport_PublishSync(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}

	for range 10 {
		err = mqtt.PublishSync(&adr, msg)
		if err != nil {
			log.Info("Publish failed err:", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	mqtt.Stop()
}

func TestMqttTransport_SubUnsub(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	mqtt.SetMessageHandler(onMsg)
	if err := mqtt.Subscribe("pt:j1/mt:evt/#"); err != nil {
		t.Fatal("Subscribe err:", err)
		return
	}

	// unsubscribe and send message, shall not receive it
	mqtt.Unsubscribe("pt:j1/mt:evt/#")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.PublishSync(&adr, msg)

	select {
	case <-msgChan:
		t.Error("Should not receive msg")
	case <-time.After(2 * time.Second):
	}

	mqtt.Stop()
}

// TODO: Fix, awsiot.private.key is not available in the repo
/*func TestMqttTransport_PublishTls(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// for test replace XYZ with actual AWS IoT core address and ABC with actual clientid
	mqtt := NewMqttTransport("ssl://a1ds8ixdqbiw53-ats.iot.eu-central-1.amazonaws.com:443", "00000000alexdevtest", "", "", false, 1, 1)

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

	t.Log("Connected")

	mqtt.SetMessageHandler(onMsg)
	time.Sleep(100 * time.Millisecond)

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)

	result := <-msgChan

	if result != 1 {
		t.Error("Wrong message")
	}

	mqtt.Stop()
}*/

// TODO: Fix, awsiot.private.key is not available in the repo
/*func TestMqttTransport_PublishTls_2(t *testing.T) {
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

	t.Log("Connected")
	mqtt.SetMessageHandler(onMsg)
	time.Sleep(100 * time.Millisecond)

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)

	result := <-msgChan
	if result != 1 {
		t.Error("Wrong message")
	}

	mqtt.Stop()
}*/

func TestMqttTransport_TestChannels(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
		return
	}

	t.Log("Connected")
	time.Sleep(100 * time.Millisecond)

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	mqtt.RegisterChannel("chan1", chan1)
	mqtt.RegisterChannel("chan2", chan2)
	isCorrect[1] = false
	isCorrect[2] = false
	go func(msgChan MessageCh) {
		newMsg := <-chan1
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[1] = true
		}
	}(chan1)
	go func(msgChan MessageCh) {
		newMsg := <-chan2
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[2] = true
		}
	}(chan2)

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)
	time.Sleep(100 * time.Millisecond)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	if isCorrect[1] && isCorrect[2] {
		t.Log("Channel test - OK")
	} else {
		t.Error("Wrong result")
		t.Fail()
	}
}

func TestMqttTransport_TestResponder(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest-1", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
		return
	}

	t.Log("Connected")
	time.Sleep(100 * time.Millisecond)

	if err := mqtt.Subscribe("#"); err != nil {
		t.Fatal("Subscribe err:", err)
	}

	mqtt2 := NewMqttTransport("tcp://localhost:1883", "fimpgotest-2", "", "", true, 1, 1)
	err = mqtt2.Start()

	time.Sleep(100 * time.Millisecond)
	mqtt2.Subscribe("pt:j1c1/mt:rsp/rt:app/rn:response_tester/ad:1")

	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}
	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	mqtt.RegisterChannel("chan1", chan1)
	mqtt2.RegisterChannel("chan2", chan2)
	// responder
	go func(msgChan MessageCh) {
		for {
			newMsg := <-chan1
			if newMsg.Payload.Service == "tester" {
				if err := mqtt.RespondToRequest(newMsg.Payload, NewFloatMessage("evt.test.response", "test_responder", 35.5, nil, nil, nil)); err != nil {
					t.Error("Error responding to request:", err)
					return
				}
			}
		}
	}(chan1)

	var isResponseReceived bool
	go func(msgChan MessageCh) {
		for {
			newMsg := <-chan2
			if newMsg.Payload.Service == "test_responder" && newMsg.Topic == "pt:j1c1/mt:rsp/rt:app/rn:response_tester/ad:1" {
				isResponseReceived = true
			}
		}

	}(chan2)

	msg := NewFloatMessage("cmd.test.get_response", "tester", float64(35.5), nil, nil, nil)
	msg.ResponseToTopic = "pt:j1c1/mt:rsp/rt:app/rn:response_tester/ad:1"
	adr := Address{PayloadType: DefaultPayload, MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "test", ResourceAddress: "1"}
	mqtt.Publish(&adr, msg)
	time.Sleep(100 * time.Millisecond)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	mqtt.Unsubscribe("#")
	if !isResponseReceived {
		t.Error("Wrong result")
		t.Fail()
	}
}

func TestMqttTransport_TestChannelsWithFilters(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
		return
	}

	t.Log("Connected")
	time.Sleep(100 * time.Millisecond)

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

	isCorrect[1] = false
	isCorrect[2] = false
	isCorrect[3] = false
	isCorrect[4] = true
	isCorrect[5] = false
	go func(msgChan MessageCh) {
		newMsg := <-msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[1] = true
		}
	}(chan1)
	go func(msgChan MessageCh) {
		newMsg := <-msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[2] = true
		}
	}(chan2)

	go func(msgChan MessageCh) {
		newMsg := <-msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[3] = true
		}
	}(chan3)
	// Negative test
	go func(msgChan MessageCh) {
		<-msgChan
		isCorrect[4] = false
	}(chan4)

	go func(msgChan MessageCh) {
		<-msgChan
		isCorrect[5] = true
	}(chan5)

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)
	time.Sleep(100 * time.Millisecond)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	mqtt.UnregisterChannel("chan3")
	mqtt.UnregisterChannel("chan4")
	mqtt.UnregisterChannel("chan5")
	if isCorrect[1] && isCorrect[2] && isCorrect[3] && isCorrect[4] && isCorrect[5] {
		t.Log("Channel test - OK")
	} else {
		t.Error("Wrong result")
		t.Log(isCorrect)
		t.Fail()
	}
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
