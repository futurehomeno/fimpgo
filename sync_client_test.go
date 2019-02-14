package fimpgo

import (
	"sync/atomic"
	"testing"
	log "github.com/sirupsen/logrus"
	"time"
)

func TestSyncClient_Connect(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	inboundChan := make(MessageCh,20)
	// starting responder
	go func (msgChanS MessageCh) {
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report"{
				t.Log("Responde . New message. uid = ",msg.Payload.UID)
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(40.0), nil, nil, msg.Payload)
				t.Log("Correlation id = ",responseMsg.CorrelationID)
				mqtt.Publish(&adr,responseMsg)
			}

		}

	}(inboundChan)
	mqtt.Subscribe("pt:j1/mt:cmd/rt:app/rn:testapp/ad:1")
	mqtt.RegisterChannel("test",inboundChan)

	// Actual test
	syncClient := NewSyncClientV2(nil,20,20)
	syncClient.Connect("tcp://localhost:1883","fimpgotest2","","",true,1,1)
	syncClient.AddSubscription("pt:j1/mt:evt/rt:app/rn:testapp/ad:1")
	var counter int32
	for it:=0 ;it<100;it++ {
		i := it
		go func() {
			t.Log("Iteration = ",i)
			adr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
			msg := NewFloatMessage("cmd.sensor.get_report", "temp_sensor", float64(35.5), nil, nil, nil)
			response,err := syncClient.SendFimp(adr.Serialize(),msg,10)
			if err != nil {
				t.Error("Error",err)
				t.Fail()
			}
			val , _ := response.GetFloatValue()
			if val != 40.0 {
				t.Error("Wong result")
				t.Fail()
			}
			atomic.AddInt32(&counter,1)
			t.Log("Iteration Done = ",i)
		}()
	}

	for 100 >counter {
		time.Sleep(1 * time.Second)
	}


	syncClient.Stop()
	if counter!=100 {
		t.Error("Wong counter value")
		t.Fail()
	}
	t.Log("SyncClientConnect test - OK")

}




func TestSyncClient_SendFimp(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	inboundChan := make(MessageCh)
	// starting responder
	go func (msgChanS MessageCh) {
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report"{
				t.Log("Responde . New message. uid = ",msg.Payload.UID)
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(40.0), nil, nil, msg.Payload)
				t.Log("Correlation id = ",responseMsg.CorrelationID)
				mqtt.Publish(&adr,responseMsg)
			}

		}

	}(inboundChan)
	mqtt.RegisterChannel("test",inboundChan)
	// Actual test
	syncClient := NewSyncClient(mqtt)
	syncClient.AddSubscription("#")
	counter := 0
	for i:=0 ;i<5;i++ {
		t.Log("Iteration = ",i)
		adr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
		msg := NewFloatMessage("cmd.sensor.get_report", "temp_sensor", float64(35.5), nil, nil, nil)
		response,err := syncClient.SendFimp(adr.Serialize(),msg,5)
		if err != nil {
			t.Error("Error",err)
			t.Fail()
		}
		val , _ := response.GetFloatValue()
		if val != 40.0 {
			t.Error("Wong result")
			t.Fail()
		}
		counter++

	}
	syncClient.Stop()
	if counter!=5 {
		t.Error("Wong counter value")
		t.Fail()
	}
	t.Log("SyncClient test - OK")
}

func TestSyncClient_SendFimpWithTopicResponse(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883","fimpgotest","","",true,1,1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ",err)
	}
	inboundChan := make(MessageCh)
	// starting message responder
	go func (msgChanS MessageCh) {
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report"{
				t.Log("Responde . New message. uid = ",msg.Payload.UID)
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(40.0), nil, nil, nil)
				t.Log("Correlation id = ",responseMsg.CorrelationID)
				mqtt.Publish(&adr,responseMsg)
			}

		}

	}(inboundChan)
	mqtt.RegisterChannel("test",inboundChan)
	// Actual test
	syncClient := NewSyncClient(mqtt)
	syncClient.AddSubscription("#")
	counter := 0
	for i:=0 ;i<5;i++ {
		t.Log("Iteration = ",i)
		reqAddr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
		respAddr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}

		msg := NewFloatMessage("cmd.sensor.get_report", "temp_sensor", float64(35.5), nil, nil, nil)
		response,err := syncClient.SendFimpWithTopicResponse(reqAddr.Serialize(),msg,respAddr.Serialize(),"temp_sensor","evt.sensor.report",5)
		if err != nil {
			t.Error("Error",err)
			t.Fail()
		}
		val , _ := response.GetFloatValue()
		if val != 40.0 {
			t.Error("Wong result")
			t.Fail()
		}
		counter++

	}
	syncClient.Stop()
	if counter!=5 {
		t.Error("Wong counter value")
		t.Fail()
	}
	t.Log("SyncClient test - OK")

}
