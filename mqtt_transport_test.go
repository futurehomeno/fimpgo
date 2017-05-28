package fimpgo

import (
	"testing"
	"time"
	log "github.com/Sirupsen/logrus"
)

var msgChan = make(chan int)

func onMsg(topic string, addr *Address, iotMsg *FimpMessage){
	if addr.ServiceName == "temp_sensor" && addr.ServiceAddress == "300"{
		msgChan <- 1
	}else {
		msgChan <- 2
	}
}


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

