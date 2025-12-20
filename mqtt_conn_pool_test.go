package fimpgo

import (
	"sync/atomic"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMqttConnectionPool_GetConnection(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	template := MqttConnectionConfigs{
		ServerURI:    "tcp://127.0.0.1:1883",
		CleanSession: true,
		SubQos:       1,
		PubQos:       1,
	}

	pool := NewMqttConnectionPool(0, 1, 10, 5*time.Second, template, "pool_test_")
	pool.Start()
	idt, _, err := pool.BorrowConnection()
	require.NoError(t, err)
	pool.ReturnConnection(idt)

	id1, conn1, err := pool.BorrowConnection()
	require.NoError(t, err)
	id2, conn2, err := pool.BorrowConnection()
	require.NoError(t, err)
	id3, responderConn, err := pool.BorrowConnection()
	require.NoError(t, err)

	msg1 := NewStringMessage("cmd.test.get_response", "tester", "test-1", nil, nil, nil)
	msg1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
	msg1_1 := NewStringMessage("cmd.test.get_response", "tester", "test-3", nil, nil, nil)
	msg1_1.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:1"
	msg2 := NewStringMessage("cmd.test.get_response", "tester", "test-2", nil, nil, nil)
	msg2.ResponseToTopic = "pt:j1/mt:rsp/rt:app/rn:goland/ad:2"

	var isResp1 atomic.Bool
	var isResp2 atomic.Bool
	var isResp3 atomic.Bool

	responderConn.SetMessageHandler(func(topic string, addr *Address, iotMsg *FimpMessage, rawPayload []byte) {
		val, _ := iotMsg.GetStringValue()
		switch val {
		case "test-1":
			isResp1.Store(true)
		case "test-2":
			isResp2.Store(true)
		case "test-3":
			isResp3.Store(true)
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

	_, con1_1, err := pool.BorrowConnection()
	require.NoError(t, err)

	if err := con1_1.PublishToTopic("pt:j1/mt:cmd/rt:app/rn:conn_pool_tester/ad:1", msg1_1); err != nil {
		t.Fatal("Publish error:", err)
	}

	var i int
	for {
		if i > 3 {
			log.Errorf("Failed %t %t %t", isResp1.Load(), isResp2.Load(), isResp3.Load())
			t.Fail()
		}

		time.Sleep(100 * time.Millisecond)
		if isResp1.Load() && isResp2.Load() && isResp3.Load() {
			break
		}
		i++
	}

	pool.ReturnConnection(id1)
	pool.ReturnConnection(id2)
	pool.ReturnConnection(id3)
	pool.Stop()
}
