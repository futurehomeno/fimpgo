package fimpgo

import (
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSyncClient_Connect(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	inboundChan := make(MessageCh, 20)

	if err := mqtt.Subscribe("pt:j1/mt:cmd/rt:app/rn:testapp/ad:1"); err != nil {
		t.Fatal("Subscribe error:", err)
	}

	mqtt.RegisterChannel("test", inboundChan)

	// Actual test
	syncClient := NewSyncClientV2(nil, 20, 20)
	if err := syncClient.Connect("tcp://127.0.0.1:1883", "fimpgotest2", "", "", true, 1, 1); err != nil {
		t.Fatal("Error connecting sync client to broker ", err)
	}

	if err := syncClient.AddSubscription("pt:j1/mt:evt/rt:app/rn:testapp/ad:1"); err != nil {
		t.Fatal("Error adding subscription ", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// starting responder
	go func(msgChanS MessageCh) {
		wg.Done()
		testVal := float64(0.0)
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report" {
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", testVal, nil, nil, msg.Payload)
				if err := mqtt.PublishToTopic("pt:j1/mt:evt/rt:app/rn:testapp/ad:1", responseMsg); err != nil {
					t.Error("Publish error:", err)
					t.Fail()
				}

				testVal += 0.1
			}
		}
	}(inboundChan)

	wg.Wait()

	iterations := 500
	expVal := float64(0.0)

	for range iterations {
		msg := NewNullMessage("cmd.sensor.get_report", "temp_sensor", nil, nil, nil)
		response, err := syncClient.SendFimp("pt:j1/mt:cmd/rt:app/rn:testapp/ad:1", msg, 1)
		if err != nil {
			t.Fatalf("SendFimp err %v", err)
		}
		val, err := response.GetFloatValue()
		if err != nil {
			t.Fatalf("SendFimp err %v", err)
		}

		if val != expVal {
			t.Fatalf("Wrong result exp=%.2f got=%.2f ", expVal, val)
		}

		expVal += 0.1
	}

	syncClient.Stop()
	mqtt.Stop()
}

func TestSyncClient_SendFimp(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	inboundChan := make(MessageCh)

	mqtt.RegisterChannel("test", inboundChan)
	// Actual test
	syncClient := NewSyncClient(mqtt)
	syncClient.AddSubscription("#")

	var wg sync.WaitGroup
	wg.Add(1)

	// starting responder
	go func(msgChanS MessageCh) {
		wg.Done()
		temp := float64(0.0)
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report" {
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", temp, nil, nil, msg.Payload)
				temp += 0.1
				mqtt.Publish(&adr, responseMsg)
			}
		}
	}(inboundChan)

	wg.Wait()

	expVal := float64(0.0)
	for range 5 {
		adr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
		msg := NewNullMessage("cmd.sensor.get_report", "temp_sensor", nil, nil, nil)
		response, err := syncClient.SendFimp(adr.Serialize(), msg, 2)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		val, err := response.GetFloatValue()
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		if val != expVal {
			t.Fatalf("Wrong result exp=%.2f got=%.2f ", expVal, val)
		}

		expVal += 0.1
	}

	syncClient.Stop()
	mqtt.Stop()
}

func TestSyncClient_SendFimpWithTopicResponse(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := NewMqttTransport("tcp://127.0.0.1:1883", "fimpgotest", "", "", true, 1, 1)
	err := mqtt.Start()
	if err != nil {
		t.Fatal("Start MQTT err:", err)
	}

	inboundChan := make(MessageCh)
	mqtt.RegisterChannel("test", inboundChan)
	// Actual test
	syncClient := NewSyncClient(mqtt)
	syncClient.AddSubscription("#")

	var wg sync.WaitGroup
	wg.Add(1)

	// starting message responder
	go func(msgChanS MessageCh) {
		wg.Done()
		temp := float64(0.0)
		for msg := range msgChanS {
			if msg.Payload.Type == "cmd.sensor.get_report" {
				adr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
				responseMsg := NewFloatMessage("evt.sensor.report", "temp_sensor", temp, nil, nil, nil)
				mqtt.Publish(&adr, responseMsg)
				temp += 0.1
			}
		}
	}(inboundChan)
	wg.Wait()

	expVal := float64(0.0)

	for range 5 {
		reqAddr := Address{MsgType: MsgTypeCmd, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}
		respAddr := Address{MsgType: MsgTypeEvt, ResourceType: ResourceTypeApp, ResourceName: "testapp", ResourceAddress: "1"}

		msg := NewNullMessage("cmd.sensor.get_report", "temp_sensor", nil, nil, nil)
		response, err := syncClient.SendFimpWithTopicResponse(reqAddr.Serialize(), msg, respAddr.Serialize(), "temp_sensor", "evt.sensor.report", 5)
		if err != nil {
			t.Fatalf("Error %v", err)
		}

		val, _ := response.GetFloatValue()
		if val != expVal {
			t.Fatalf("Wrong result exp=%.2f got=%.2f ", expVal, val)
		}

		expVal += 0.1
	}

	syncClient.Stop()
	mqtt.Stop()
}
