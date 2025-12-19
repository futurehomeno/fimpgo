package fimpgo

import (
	"math/rand"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSyncClient_Connect(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Error connecting to broker ", err)
	}
	t.Log("Connected")

	inboundChan := make(MessageCh, 20)
	// starting responder
	go func(msgChanS MessageCh) {
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report" {
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(40.0), nil, nil, msg.Payload)
				if err := mqtt.PublishToTopic("pt:j1/mt:evt/rt:app/rn:testapp/ad:1", responseMsg); err != nil {
					t.Error("Publish error:", err)
					t.Fail()
				}
			}
		}
	}(inboundChan)

	if err := mqtt.Subscribe("pt:j1/mt:cmd/rt:app/rn:testapp/ad:1"); err != nil {
		t.Fatal("Subscribe error:", err)
	}

	mqtt.RegisterChannel("test", inboundChan)

	// Actual test
	syncClient := NewSyncClientV2(nil, 20, 20)
	if err := syncClient.Connect("tcp://localhost:1883", "fimpgotest2", "", "", true, 1, 1); err != nil {
		t.Fatal("Error connecting sync client to broker ", err)
	}

	if err := syncClient.AddSubscription("pt:j1/mt:evt/rt:app/rn:testapp/ad:1"); err != nil {
		t.Fatal("Error adding subscription ", err)
	}

	iterations := 1000
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	for range iterations {
		go func() {
			testVal := (rand.Intn(2800) + 1000) / 10.0
			msg := NewFloatMessage("cmd.sensor.get_report", "temp_sensor", float64(testVal), nil, nil, nil)
			response, err := syncClient.SendFimp("pt:j1/mt:cmd/rt:app/rn:testapp/ad:1", msg, 10)
			if err != nil {
				t.Error("SendFimp err", err)
				t.Fail()
			}
			val, _ := response.GetFloatValue()
			if val != float64(testVal) {
				t.Errorf("Wong result exp=%.2f got=%.2f ", float64(testVal), val)
				t.Fail()
			}
		}()

		waitgroup.Done()
	}

	waitgroup.Wait()
	syncClient.Stop()
}

func TestSyncClient_SendFimp(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Error connecting to broker ", err)
	}
	t.Log("Connected")

	inboundChan := make(MessageCh)
	// starting responder
	go func(msgChanS MessageCh) {
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report" {
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(40.0), nil, nil, msg.Payload)
				mqtt.Publish(&adr, responseMsg)
			}
		}

	}(inboundChan)
	mqtt.RegisterChannel("test", inboundChan)
	// Actual test
	syncClient := NewSyncClient(mqtt)
	syncClient.AddSubscription("#")

	for range 5 {
		testVal := (rand.Intn(2800) + 1000) / 10.0
		adr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
		msg := NewFloatMessage("cmd.sensor.get_report", "temp_sensor", float64(testVal), nil, nil, nil)
		response, err := syncClient.SendFimp(adr.Serialize(), msg, 5)
		if err != nil {
			t.Error("Error", err)
			t.Fail()
		}
		val, _ := response.GetFloatValue()
		if val != float64(testVal) {
			t.Errorf("Wong result exp=%.2f got=%.2f ", float64(testVal), val)
			t.Fail()
		}
	}

	syncClient.Stop()
}

func TestSyncClient_SendFimpWithTopicResponse(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://localhost:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Error connecting to broker ", err)
	}

	t.Log("Connected")
	inboundChan := make(MessageCh)
	// starting message responder
	go func(msgChanS MessageCh) {
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report" {
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", float64(40.0), nil, nil, nil)
				mqtt.Publish(&adr, responseMsg)
			}

		}

	}(inboundChan)
	mqtt.RegisterChannel("test", inboundChan)
	// Actual test
	syncClient := NewSyncClient(mqtt)
	syncClient.AddSubscription("#")
	counter := 0
	for i := range 5 {
		t.Log("Iteration = ", i)
		reqAddr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
		respAddr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}

		msg := NewFloatMessage("cmd.sensor.get_report", "temp_sensor", float64(35.5), nil, nil, nil)
		response, err := syncClient.SendFimpWithTopicResponse(reqAddr.Serialize(), msg, respAddr.Serialize(), "temp_sensor", "evt.sensor.report", 5)
		if err != nil {
			t.Error("Error", err)
			t.Fail()
		}
		val, _ := response.GetFloatValue()
		if val != 40.0 {
			t.Error("Wong result")
			t.Fail()
		}
		counter++

	}
	syncClient.Stop()
	if counter != 5 {
		t.Error("Wong counter value")
		t.Fail()
	}
	t.Log("SyncClient test - OK")

}
