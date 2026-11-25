package fimpgo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/futurehomeno/fimpgo/utils"
)

func newCovPool(t *testing.T, maxSize int) *MqttConnectionPool {
	t.Helper()

	pool := NewMqttConnectionPool(0, 1, maxSize, 0, MqttConnectionConfigs{ServerURI: covBroker, CleanSession: true, SubQos: 1, PubQos: 1}, "cov"+uuid.NewString()[:8])
	t.Cleanup(pool.Stop)

	return pool
}

func TestMqttConnectionPool_Defaults(t *testing.T) {
	pool := NewMqttConnectionPool(50, 0, 0, 0, MqttConnectionConfigs{}, "x")
	defer pool.Stop()

	assert.Equal(t, 0, pool.initSize)
	assert.Equal(t, 2, pool.size)
	assert.Equal(t, 20, pool.maxSize)
	assert.Equal(t, 30*time.Second, pool.maxIdleAge)

	pool.Start()
	pool.Stop()
	pool.Stop()
}

func TestMqttConnectionPool_BorrowReturnCounters(t *testing.T) {
	pool := newCovPool(t, 2)

	assert.False(t, pool.IsConnected(1))

	id1, c1, err := pool.BorrowConnection(nil)
	require.NoError(t, err)
	assert.True(t, pool.IsConnected(id1))
	assert.Equal(t, 1, pool.TotalConnections())
	assert.Equal(t, 0, pool.IdleConnections())

	id2, _, err := pool.BorrowConnection(nil)
	require.NoError(t, err)
	assert.Equal(t, 2, pool.TotalConnections())

	_, c3, err := pool.BorrowConnection(nil)
	assert.ErrorContains(t, err, "too many connections")
	assert.Nil(t, c3)

	require.NoError(t, c1.Subscribe("pt:j1/mt:evt/rt:app/rn:covpool/ad:1"))
	pool.ReturnConnection(id1)
	pool.ReturnConnection(999)
	assert.Equal(t, 1, pool.IdleConnections())
	assert.Empty(t, c1.subs)

	idAgain, _, err := pool.BorrowConnection(nil)
	require.NoError(t, err)
	assert.Equal(t, id1, idAgain)

	pool.ReturnConnection(id2)
}

func TestMqttConnectionPool_BorrowSkipsDisconnectedIdle(t *testing.T) {
	pool := newCovPool(t, 3)

	id, conn, err := pool.BorrowConnection(nil)
	require.NoError(t, err)
	pool.ReturnConnection(id)
	conn.client.Disconnect(0)
	require.Eventually(t, func() bool { return !conn.client.IsConnected() }, 5*time.Second, time.Millisecond)

	id2, _, err := pool.BorrowConnection(nil)
	require.NoError(t, err)
	assert.NotEqual(t, id, id2)
}

func TestMqttConnectionPool_CreateWithBadTLSConfig(t *testing.T) {
	pool := NewMqttConnectionPool(0, 1, 1, 0, MqttConnectionConfigs{PrivateKeyFileName: "k", CertFileName: "c", CertDir: t.TempDir()}, "x")
	defer pool.Stop()

	_, conn, err := pool.BorrowConnection(nil)
	assert.ErrorContains(t, err, "not created")
	assert.Nil(t, conn)
}

func TestMqttConnectionPool_CleanupRemovesStaleIdle(t *testing.T) {
	pool := newCovPool(t, 5)
	pool.Stop()

	ids := make([]int, 3)

	for i := range ids {
		id, _, err := pool.BorrowConnection(nil)
		require.NoError(t, err)

		ids[i] = id
	}

	for _, id := range ids {
		pool.ReturnConnection(id)
	}

	require.Equal(t, 3, pool.TotalConnections())
	pool.connPool[ids[0]].idleSince = time.Now().Add(-time.Hour)
	pool.connPool[ids[1]].idleSince = time.Now().Add(-time.Hour)

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	pool.mux.Lock()
	pool.stopChan = make(chan struct{})
	pool.poolCheckTick = ticker
	pool.mux.Unlock()
	pool.isStarted.Store(true)

	go pool.cleanupProcess()

	require.Eventually(t, func() bool { return pool.TotalConnections() == 1 }, 5*time.Second, 5*time.Millisecond)
}

func TestSyncClient_NotConnected(t *testing.T) {
	sc := NewSyncClientV2(nil)
	sc.SetGlobalPrefix("site")
	sc.Stop()

	assert.EqualError(t, sc.AddSubscription("a"), "not connected")
	assert.EqualError(t, sc.RemoveSubscription("a"), "not connected")

	_, err := sc.SendFimp("pt:j1/mt:cmd/rt:app/rn:n/ad:1", NewNullMessage("cmd.a.b", "s", nil, nil, nil), 1)
	assert.EqualError(t, err, "not connected")
}

func TestSyncClient_ConnectSubscribeAndReqResp(t *testing.T) {
	const (
		cmdTopic = "pt:j1/mt:cmd/rt:app/rn:covsync/ad:1"
		rspTopic = "pt:j1/mt:rsp/rt:app/rn:covsync/ad:1"
	)

	responder := startCovTransport(t)
	require.NoError(t, responder.Subscribe(cmdTopic))
	responder.SetMessageHandler(func(_ string, _ *Address, msg *FimpMessage, _ []byte) {
		if msg.Interface == "cmd.cov.get" {
			assert.NoError(t, responder.PublishToTopic(rspTopic, NewStringMessage("evt.cov.report", "cov", "pong", nil, nil, msg)))
		}
	})

	sc := NewSyncClient(nil)
	require.NoError(t, sc.Connect(covBroker, "cov_"+uuid.NewString()[:8], "", "", true, 1, 1, nil))
	require.NoError(t, sc.Connect(covBroker, "ignored", "", "", true, 1, 1, nil))
	defer sc.Stop()

	require.NoError(t, sc.AddSubscription(rspTopic))
	require.NoError(t, sc.RemoveSubscription(rspTopic))

	rsp, err := sc.SendReqRespFimp(cmdTopic, rspTopic, NewNullMessage("cmd.cov.get", "cov", nil, nil, nil), 5, true)
	require.NoError(t, err)
	assert.Equal(t, "pong", rsp.Value)

	_, err = sc.SendReqRespFimp(cmdTopic, rspTopic, NewNullMessage("cmd.cov.silent", "cov", nil, nil, nil), 1, true)
	assert.ErrorIs(t, err, utils.ErrTimeout)
}

func TestSyncClient_ConnPool(t *testing.T) {
	const (
		cmdTopic = "pt:j1/mt:cmd/rt:app/rn:covpoolsync/ad:1"
		rspTopic = "pt:j1/mt:rsp/rt:app/rn:covpoolsync/ad:1"
	)

	responder := startCovTransport(t)
	require.NoError(t, responder.Subscribe("site/"+cmdTopic))
	responder.SetMessageHandler(func(_ string, _ *Address, msg *FimpMessage, _ []byte) {
		assert.NoError(t, responder.PublishToTopic("site/"+rspTopic, NewStringMessage("evt.cov.report", "cov", "pong", nil, nil, msg)))
	})

	pool := newCovPool(t, 2)
	sc := &SyncClient{mqttConnPool: pool, isConnPoolEnabled: true, globalPrefix: "site"}

	rsp, err := sc.SendReqRespFimp(cmdTopic, rspTopic, NewNullMessage("cmd.cov.get", "cov", nil, nil, nil), 5, true)
	require.NoError(t, err)
	assert.Equal(t, "pong", rsp.Value)
	assert.Equal(t, 1, pool.IdleConnections())

	pool.Stop()
	pool.maxSize = 0
	_, err = sc.SendFimp(cmdTopic, NewNullMessage("cmd.cov.get", "cov", nil, nil, nil), 1)
	assert.ErrorContains(t, err, "too many connections")
}
