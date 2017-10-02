package fimpgo

import (
	"testing"
	"time"
	log "github.com/Sirupsen/logrus"
)

var msgChan = make(chan int)

func onMsg(topic string, addr *Address, iotMsg *FimpMessage,rawMessage []byte){
	if addr.ServiceName == "temp_sensor" && addr.ServiceAddress == "300"{
		msgChan <- 1
	}else {
		msgChan <- 2
	}
}

var isCorrect1 bool
var isCorrect2 bool


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
	isCorrect1 = false
	isCorrect2 = false
	go func(msgChan MessageCh) {
	   newMsg :=<- msgChan
	   if newMsg.Payload.Service == "temp_sensor" {
	   		isCorrect1 = true
	   }
	}(chan1)
	go func(msgChan MessageCh) {
		newMsg :=<- msgChan
		if newMsg.Payload.Service == "temp_sensor" {
			isCorrect2 = true
		}
	}(chan2)

	msg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(35.5), nil, nil, nil)
	adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr,msg)
	time.Sleep(time.Second*1)
	mqtt.UnregisterChannel("chan1")
	mqtt.UnregisterChannel("chan2")
	if !isCorrect1 || !isCorrect2 {
		t.Error("Wrong result")
		t.Fail()
	}
	t.Log("Channel test - OK")

}