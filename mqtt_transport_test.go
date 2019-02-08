package fimpgo

import (
	"testing"
	"time"
	log "github.com/sirupsen/logrus"
)

var msgChan = make(chan int)

func onMsg(topic string, addr *Address, iotMsg *FimpMessage,rawMessage []byte){
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

func TestMqttTransport_PublishTls(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// for test replace XYZ with actual AWS IoT core address and ABC with actual clientid
	mqtt := NewMqttTransport("ssl://XYZ.amazonaws.com:443","ABC","","",true,1,1)

	// for test enter valid site-id
	mqtt.SetGlobalTopicPrefix("XXX")
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

