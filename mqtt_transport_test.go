package fimpgo

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)



var msgChan = make(chan int)

func onMsg(topic string, addr *Address, iotMsg *FimpMessage,rawMessage []byte){
	log.Info("New message")
	if addr.ServiceName == "temp_sensor" && addr.ServiceAddress == "300"{
		msgChan <- 1
	}else {
		msgChan <- 2
	}
}

var isCorrect = make(map[int]bool)



func TestMqttTransport_Publish(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	mqtt.SetMessageHandler(onMsg)
	mqtt.Subscribe("#")
	t.Log("Publishing message")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)

	t.Log("Waiting for new message")
	result := <- msgChan
	t.Log("Got new message")
	mqtt.Stop()
	if result != 1 {
		t.Error("Wrong message")
	}

}

func TestMqttTransport_PublishStopPublish(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	mqtt.SetMessageHandler(onMsg)
	mqtt.Subscribe("#")
	t.Log("Publishing message")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)

	t.Log("Waiting for new message")
	result := <- msgChan
	t.Log("Got new message")
	mqtt.Stop()
	if result != 1 {
		t.Error("Wrong message")
	}
	time.Sleep(time.Second*5)
	mqtt = NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err = mqtt.Start()
	t.Log("Connected 2")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	time.Sleep(time.Second*5)

	t.Log("Done")
	mqtt.Stop()

}


func TestMqttTransport_PublishSync(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	MQTT.DEBUG = log.StandardLogger()
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	t.Log("Publishing message")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}

	for i:=0;i<10;i++ {
		err = mqtt.PublishSync(&adr,msg)
		if err != nil {
			log.Info("Publish failed . Err :",)
		}else {
			log.Info("Publish success ")
		}
		time.Sleep(time.Second*5)

	}

	t.Log("Waiting for new message")
	t.Log("Got new message")
	mqtt.Stop()

}

func TestMqttTransport_SubUnsub(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	mqtt.SetMessageHandler(onMsg)
	mqtt.Subscribe("pt:j1/mt:evt/#")
	//mqtt.Subscribe("pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:temp_sensor/ad:300")
	//mqtt.Unsubscribe("pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:temp_sensor/ad:300")
	mqtt.Unsubscribe("pt:j1/mt:evt/#")
	t.Log("Publishing message")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.PublishSync(&adr,msg)

	t.Log("Waiting for new message")
	result := <- msgChan
	t.Log("Got new message")
	mqtt.Stop()
	if result != 1 {
		t.Error("Wrong message")
	}

}

func TestMqttTransport_PublishTls(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// for test replace XYZ with actual AWS IoT core address and ABC with actual clientid
	mqtt := NewMqttTransport("ssl://a1ds8ixdqbiw53-ats.iot.eu-central-1.amazonaws.com:443","00000000alexdevtest","","",false,1,1)

	// for test enter valid site-id
	mqtt.SetGlobalTopicPrefix("331D092F-4685-4CC9-8337-2598E6F5D8D5")
	// for test place certificate and key into certs folder
	err := mqtt.ConfigureTls("awsiot.private.key","awsiot.crt","./certs",true)

	if err != nil {
		t.Error("Certificate error :",err)
	}

	err = mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	mqtt.SetMessageHandler(onMsg)
	time.Sleep(time.Second*1)
	mqtt.Subscribe("#")
	t.Log("Publishing message")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)

	t.Log("Waiting for new message")
	result := <- msgChan
	t.Log("Got new message")
	mqtt.Stop()
	if result != 1 {
		t.Error("Wrong message")
	}

}

func TestMqttTransport_PublishTls_2(t *testing.T) {

	connConfig := MqttConnectionConfigs{
		ServerURI:           "ssl://a1ds8ixdqbiw53-ats.iot.eu-central-1.amazonaws.com:443",
		ClientID:            "00000000alexdevtest",
		CleanSession:        true,
		SubQos:              1,
		PubQos:              1,
		CertDir:             "./certs",
		PrivateKeyFileName:  "awsiot.private.key",
		CertFileName:        "awsiot.crt",
	}

	log.SetLevel(log.DebugLevel)
	// for test replace XYZ with actual AWS IoT core address and ABC with actual clientid
	mqtt := NewMqttTransportFromConfigs(connConfig)

	// for test enter valid site-id
	mqtt.SetGlobalTopicPrefix("331D092F-4685-4CC9-8337-2598E6F5D8D5")
	// for test place certificate and key into certs folder
	err := mqtt.ConfigureTls("awsiot.private.key","awsiot.crt","./certs",true)

	if err != nil {
		t.Error("Certificate error :",err)
	}

	err = mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}

	mqtt.SetMessageHandler(onMsg)
	time.Sleep(time.Second*1)
	mqtt.Subscribe("#")
	t.Log("Publishing message")

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)

	t.Log("Waiting for new message")
	result := <- msgChan
	t.Log("Got new message")
	mqtt.Stop()
	if result != 1 {
		t.Error("Wrong message")
	}

}


func TestMqttTransport_TestChannels(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	time.Sleep(time.Second*1)
	mqtt.Subscribe("#")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	mqtt.RegisterChannel("chan1",chan1)
	mqtt.RegisterChannel("chan2",chan2)
	isCorrect[1] = false
	isCorrect[2] = false
	go func(msgChan MessageCh) {
	   newMsg :=<- chan1
	   if newMsg.Payload.Service == "temp_sensor" {
	   		isCorrect[1] = true
	   }
	}(chan1)
	go func(msgChan MessageCh) {
		newMsg :=<- chan2
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[2] = true
		}
	}(chan2)

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)
	time.Sleep(time.Second*1)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	if isCorrect[1] && isCorrect[2] {
		t.Log("Channel test - OK")
	}else {
		t.Error("Wrong result")
		t.Fail()
	}
}

func TestMqttTransport_TestResponder(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	var isResponseReceived bool
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest-1","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	time.Sleep(time.Second*1)
	mqtt.Subscribe("#")

	mqtt2 := NewMqttTransport("tcp://localhost:1883","fimpgotest-2","","",true,1,1)
	err = mqtt2.Start()
	t.Log("Connected")
	time.Sleep(time.Second*1)
	mqtt2.Subscribe("pt:j1/mt:rsp/rt:app/rn:response_tester/ad:1")


	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	mqtt.RegisterChannel("chan1",chan1)
	mqtt2.RegisterChannel("chan2",chan2)
	// responder
	go func(msgChan MessageCh) {
		newMsg :=<- chan1
		if newMsg.Payload.Service == "tester" {
			mqtt.RespondToRequest(newMsg.Payload,NewFloatMessage("evt.test.response", "test_responder", float64(35.5), nil, nil, nil))
		}
	}(chan1)

	go func(msgChan MessageCh) {
			newMsg :=<- chan2
			t.Log("Service = "+newMsg.Payload.Service)
			if newMsg.Payload.Service == "test_responder" && newMsg.Topic == "pt:j1/mt:rsp/rt:app/rn:response_tester/ad:1" {
				isResponseReceived = true
			}

	}(chan2)

	msg := NewFloatMessage("cmd.test.get_response", "tester", float64(35.5), nil, nil, nil)
	msg.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:response_tester/ad:1"
	adr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "test", ResourceAddress: "1"}
	mqtt.Publish(&adr,msg)
	time.Sleep(time.Second*2)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	mqtt.Unsubscribe("#")
	if isResponseReceived {
		t.Log("Response received")
	}else {
		t.Error("Wrong result")
		t.Fail()
	}
}

func TestMqttTransport_TestChannelsWithFilters(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	time.Sleep(time.Second*1)
	mqtt.Subscribe("#")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	chan1 := make(MessageCh)
	chan2 := make(MessageCh)
	chan3 := make(MessageCh)
	chan4 := make(MessageCh)
	chan5 := make(MessageCh)
	mqtt.RegisterChannel("chan1",chan1)
	mqtt.RegisterChannel("chan2",chan2)
	mqtt.RegisterChannelWithFilter("chan3",chan3,FimpFilter{
		Topic:     "pt:j1/mt:evt/#",
		Service:   "temp_sensor",
		Interface: "evt.sensor.report",
	})

	mqtt.RegisterChannelWithFilter("chan4",chan4,FimpFilter{
		Topic:     "pt:j1/mt:cmd/#",
		Service:   "temp_sensor",
		Interface: "cmd.sensor.report",
	})

	testFilterFunc := func (topic string, addr *Address, iotMsg *FimpMessage) bool {
		if iotMsg.Type == "evt.sensor.report"{
			return true
		}
		return false
	}

	mqtt.RegisterChannelWithFilterFunc("chan5",chan5,testFilterFunc)

	isCorrect[1] = false
	isCorrect[2] = false
	isCorrect[3] = false
	isCorrect[4] = true
	isCorrect[5] = false
	go func(msgChan MessageCh) {
		newMsg :=<- msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[1] = true
		}
	}(chan1)
	go func(msgChan MessageCh) {
		newMsg :=<- msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[2] = true
		}
	}(chan2)

	go func(msgChan MessageCh) {
		newMsg :=<- msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect[3] = true
		}
	}(chan3)
	// Negative test 
	go func(msgChan MessageCh) {
		_=<- msgChan
		isCorrect[4] = false

	}(chan4)

	go func(msgChan MessageCh) {
		_=<- msgChan
		isCorrect[5] = true

	}(chan5)

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)
	time.Sleep(time.Second*1)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	mqtt.UnregisterChannel("chan3")
	mqtt.UnregisterChannel("chan4")
	mqtt.UnregisterChannel("chan5")
	if isCorrect[1] && isCorrect[2] && isCorrect[3] && isCorrect[4] && isCorrect[5]{
		t.Log("Channel test - OK")
	}else {
		t.Error("Wrong result")
		t.Log(isCorrect)
		t.Fail()
	}


}

func TestAddGlobalPrefixToTopic(t *testing.T) {
	result := AddGlobalPrefixToTopic("12345","pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if result != "12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}else {
		t.Log("AddGlobalPrefixToTopic test 1 - OK")
	}
	result = AddGlobalPrefixToTopic("12345","/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if result != "12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}else {
		t.Log("AddGlobalPrefixToTopic test 2 - OK")
	}
	result = AddGlobalPrefixToTopic("","pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if result != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}else {
		t.Log("AddGlobalPrefixToTopic test 3 - OK")
	}
}

func TestDetachGlobalPrefixFromTopic(t *testing.T) {
	globalPrefix,topic := DetachGlobalPrefixFromTopic("12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if globalPrefix != "12345" || topic != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}else {
		t.Log("DetachGlobalPrefixFromTopic test 1 - OK")
	}
	globalPrefix,topic = DetachGlobalPrefixFromTopic("ABC/12345/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0")
	if globalPrefix != "ABC/12345" || topic != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:49_0" {
		t.Error("Wrong topic")
	}else {
		t.Log("Result ,",globalPrefix,topic)
		t.Log("DetachGlobalPrefixFromTopic test 2 - OK")
	}
}

