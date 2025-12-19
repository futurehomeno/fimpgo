package fimpgo

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestMqttConnectionPool_GetConnection(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	template := MqttConnectionConfigs{
		ServerURI:    "tcp://localhost:1883",
		CleanSession: true,
		SubQos:       1,
		PubQos:       1,
	}

	pool := NewMqttConnectionPool(0, 1, 10, 5*time.Second, template, "pool_test_")
	pool.Start()
	idt, _, _ := pool.BorrowConnection()
	pool.ReturnConnection(idt)
	id1, conn1, _ := pool.BorrowConnection()
	id2, conn2, _ := pool.BorrowConnection()
	id3, responderConn, _ := pool.BorrowConnection()
	msg1 := NewStringMessage("cmd.test.get_response", "tester", "test-1", nil, nil, nil)
	msg1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
	msg1_1 := NewStringMessage("cmd.test.get_response", "tester", "test-3", nil, nil, nil)
	msg1_1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
	msg2 := NewStringMessage("cmd.test.get_response", "tester", "test-2", nil, nil, nil)
	msg2.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:2"

	isResp1 := false
	isResp2 := false
	isResp3 := false
	responderConn.SetMessageHandler(func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte) {
		val, _ := iotMsg.GetStringValue()
		switch val {
		case "test-1":
			isResp1 = true
		case "test-2":
			isResp2 = true
		case "test-3":
			isResp3 = true
		}
		response := NewStringMessage("evt.test.response", "tester", val, nil, nil, iotMsg)
		if err := responderConn.RespondToRequest(iotMsg, response); err != nil {
			log.Error("Respond to rq err:", err)
		}
	})

	if err := responderConn.Subscribe("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1"); err != nil {
		t.Fatal("Subscribe error:", err)
	}

	if err := conn1.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1", msg1); err != nil {
		t.Fatal("Publish error:", err)
	}

	if err := conn2.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1", msg2); err != nil {
		t.Fatal("Publish error:", err)
	}

	time.Sleep(100 * time.Millisecond)
	pool.ReturnConnection(id1)

	_, con1_1, _ := pool.BorrowConnection()
	if err := con1_1.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1", msg1_1); err != nil {
		t.Fatal("Publish error:", err)
	}

	var i int
	for {
		if i > 3 {
			log.Error("Failed")
			t.FailNow()
		}
		time.Sleep(100 * time.Millisecond)
		if isResp1 && isResp2 && isResp3 {
			break
		}
		i++

	}

	pool.ReturnConnection(id1)
	pool.ReturnConnection(id2)
	pool.ReturnConnection(id3)
	time.Sleep(100 * time.Millisecond)
	pool.Stop()
}
