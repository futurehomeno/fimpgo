package fimpgo

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

//var msgChan = make(chan int)

func TestNewMqttConnectionPool(t *testing.T) {

}

func TestMqttConnectionPool_GetConnection(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	template := MqttConnectionConfigs{
		ServerURI:           "tcp://localhost:1883",
		CleanSession:        true,
		SubQos:              1,
		PubQos:              1,
	}

	pool := NewMqttConnectionPool(0,2,10,template,"pool_test_")
	id1 ,conn1,_ := pool.GetConnection()
	id2 ,conn2,_ := pool.GetConnection()
	id3 ,responderConn,_ := pool.GetConnection()
	log.Infof("Connection ids %d %d %d",id1,id2,id3)
	msg1 := NewStringMessage("cmd.test.get_response","tester","test-1",nil,nil,nil)
	msg1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
	msg1_1 := NewStringMessage("cmd.test.get_response","tester","test-3",nil,nil,nil)
	msg1_1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
	msg2 := NewStringMessage("cmd.test.get_response","tester","test-2",nil,nil,nil)
	msg2.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:2"

	isResp1 := false
    isResp2 := false
    isResp3 := false
	responderConn.SetMessageHandler(func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte) {
		log.Info("New mqtt msg ")
		val,_ := iotMsg.GetStringValue()
		switch val {
		case "test-1":
			isResp1 = true
		case "test-2":
			isResp2 = true
		case "test-3":
			isResp3 = true
		}
		log.Info("Request ",val)
		response := NewStringMessage("evt.test.response","tester",val,nil,nil,iotMsg)
		responderConn.RespondToRequest(iotMsg,response)
	})

	responderConn.Subscribe("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1")

	conn1.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1",msg1)
	conn2.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1",msg2)
	time.Sleep(time.Second*1)
    pool.ReturnConnection(id1)
	id1_1 ,con1_1,_ := pool.GetConnection()
	con1_1.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1",msg1_1)
	log.Infof("Connection ids %d ",id1_1)

	var i int
    for {
    	log.Infof("Wait %d",i)
    	if i > 3 {
    		log.Error("Failed")
    		t.FailNow()
		}
    	time.Sleep(time.Second*1)
    	if isResp1 && isResp2 && isResp3{
    		break
		}
    	i++

	}
	time.Sleep(time.Second*20)
    conn1.Stop()
    conn2.Stop()
    responderConn.Stop()
    pool.ReturnConnection(id1)
    pool.ReturnConnection(id2)
    pool.ReturnConnection(id3)
	t.Log("All good . Test passed ")

}