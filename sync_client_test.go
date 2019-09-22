package fimpgo

import (
	log "github.com/sirupsen/logrus"
	"sync/atomic"
	"testing"
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
	iterations := 1000
	for it:=0 ;it<iterations;it++ {
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

	for int32(iterations) >counter {
		time.Sleep(1 * time.Second)
	}


	syncClient.Stop()
	if counter!=int32(iterations) {
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

func TestNewSyncClientV3(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	connConfig := MqttConnectionConfigs{
		ServerURI:           "tcp://localhost:1883",
		CleanSession:        true,
		SubQos:              1,
		PubQos:              1,
	}

	connPool := NewMqttConnectionPool(0,2,10,connConfig,"pool_test_")

	_ ,responderConn,_ := connPool.GetConnection()
	responderConn.SetMessageHandler(func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte) {
		log.Info("New mqtt msg ")
		val,_ := iotMsg.GetStringValue()
		response := NewStringMessage("evt.test.response","tester",val,nil,nil,iotMsg)
		go responderConn.RespondToRequest(iotMsg,response)
	})
	responderConn.Subscribe("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1")

    //-----------------------T-E-S-T----------------------------------------------------
	syncC := NewSyncClientV3(connPool,0,0)

	msg1 := NewStringMessage("cmd.test.get_response","tester","test-1",nil,nil,nil)
	msg1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"

	msg2 := NewStringMessage("cmd.test.get_response","tester","test-2",nil,nil,nil)
	msg2.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:2"

	msg3 := NewStringMessage("cmd.test.get_response","tester","test-3",nil,nil,nil)
	msg3.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
    var response3 *FimpMessage
	readyCh := make(chan bool)
	go func() {
		log.Info("----Response 3 Start ")
		response3 , _ = syncC.SendReqRespFimp("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1","pt:j1/mt:rsp/rt:app/rn:goland/ad:1",msg3,5,true)
		log.Info("----Response 3 End")
		readyCh <- true
	}()
	log.Info("----Response 1 Start")
	response , _ := syncC.SendReqRespFimp("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1","pt:j1/mt:rsp/rt:app/rn:goland/ad:1",msg1,5,true)
	log.Info("----Response 1")
	response2 , _ := syncC.SendReqRespFimp("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1","pt:j1/mt:rsp/rt:app/rn:goland/ad:2",msg2,5,true)
	log.Info("----Response 2")
    // waiting response from goroutine
	<-readyCh
	respVal , _ := response.GetStringValue()
	respVal2 , _ := response2.GetStringValue()
	respVal3 , _ := response3.GetStringValue()

	if respVal == "test-1" && respVal2 == "test-2" && respVal3 == "test-3" {
		t.Log("SUCCESS")
	}else {
		t.Error("Wrong response")
		t.Fail()
	}


}