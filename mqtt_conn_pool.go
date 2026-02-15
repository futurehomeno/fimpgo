package fimpgo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

type connection struct {
	mqConnection *MqttTransport
	isIdle       bool
	startedAt    time.Time
	idleSince    time.Time
}

// Connection pool starts at initSize connections and can grow up to maxSize , the pool can shrink back to size defined in "size" variable

type MqttConnectionPool struct {
	mux            sync.RWMutex
	connTemplate   MqttConnectionConfigs
	connPool       map[int]*connection
	nextID         uint64
	clientIdPrefix string
	initSize       int           // init size
	size           int           // normal size
	maxSize        int           // max size
	maxIdleAge     time.Duration // Defines how long idle connection can stay in the pool before it gets destroyed
	poolCheckTick  *time.Ticker
	isStarted      atomic.Bool
}

func NewMqttConnectionPool(initSize, size, maxSize int, maxAge time.Duration, connTemplate MqttConnectionConfigs, clientIdPrefix string) *MqttConnectionPool {
	if maxAge == 0 {
		maxAge = 30 * time.Second // age in seconds
	}
	if maxSize == 0 {
		maxSize = 20
	}
	if size == 0 {
		size = 2
	}
	if initSize > maxSize {
		initSize = 0
	}
	pool := &MqttConnectionPool{connPool: make(map[int]*connection), connTemplate: connTemplate, clientIdPrefix: clientIdPrefix, initSize: initSize, size: size, maxSize: maxSize, maxIdleAge: maxAge}
	pool.Start()
	return pool
}

func (cp *MqttConnectionPool) Start() {
	cp.mux.Lock()
	if !cp.isStarted.Load() {
		cp.isStarted.Store(true)
		cp.poolCheckTick = time.NewTicker(10 * time.Second)
		go cp.cleanupProcess()
	}
	cp.mux.Unlock()
}

func (cp *MqttConnectionPool) Stop() {
	cp.mux.Lock()
	if cp.poolCheckTick != nil {
		cp.poolCheckTick.Stop()
	}
	for i := range cp.connPool {
		conn := cp.connectionByID(i)
		if conn != nil {
			conn.Stop()
			delete(cp.connPool, i)
		}
	}
	cp.mux.Unlock()

	cp.isStarted.Store(false)
}

func (cp *MqttConnectionPool) IsConnected(poolID int) bool {
	ret := false
	cp.mux.RLock()
	conn, ok := cp.connPool[poolID]

	if ok {
		ret = conn.mqConnection.IsConnected()
	}
	cp.mux.RUnlock()

	return ret
}

func (cp *MqttConnectionPool) TotalConnections() int {
	cp.mux.RLock()
	size := len(cp.connPool)
	cp.mux.RUnlock()
	return size
}

func (cp *MqttConnectionPool) IdleConnections() int {
	var size int
	cp.mux.RLock()
	for i := range cp.connPool {
		if cp.connPool[i].isIdle {
			size++
		}
	}
	cp.mux.RUnlock()
	return size
}

func (cp *MqttConnectionPool) createConnection(errHandler func(error)) (int, error) {
	if len(cp.connPool) >= cp.maxSize {
		return 0, fmt.Errorf("too many connections=%d", len(cp.connPool))
	}

	connID := cp.connID()
	conf := cp.connTemplate
	conf.ClientID = fmt.Sprintf("%s_%d", cp.clientIdPrefix, connID)
	newConnection := NewMqttTransportFromConfigs(conf, errHandler)

	if newConnection == nil {
		return 0, fmt.Errorf("create TLS connection not created")
	}

	err := newConnection.Start(10 * time.Second)

	if err == nil {
		cp.connPool[connID] = &connection{
			mqConnection: newConnection,
			isIdle:       false,
			startedAt:    time.Now(),
		}
	}

	return connID, err
}

// BorrowConnection returns first available connection from the pool or creates new connection
func (cp *MqttConnectionPool) BorrowConnection(errHandler func(error)) (int, *MqttTransport, error) {
	cp.mux.Lock()
	defer cp.mux.Unlock()

	for i := range cp.connPool {
		if cp.connPool[i].isIdle {
			if cp.connPool[i].mqConnection.client.IsConnected() {
				cp.connPool[i].isIdle = false
				return i, cp.connPool[i].mqConnection, nil
			} else {
				break
			}
		}
	}

	connId, err := cp.createConnection(errHandler)
	return connId, cp.connectionByID(connId), err
}

// ReturnConnection returns connection to pool by setting inUse status to false
func (cp *MqttConnectionPool) ReturnConnection(connId int) {
	cp.mux.Lock()
	defer cp.mux.Unlock()

	con, ok := cp.connPool[connId]
	if !ok {
		log.Debugf("No connection=%d to return", connId)
		return
	}

	err := con.mqConnection.UnsubscribeAll()
	if err != nil {
		log.Warnf("[fimpgo] UnsubscribeAll err: %v", err)
	}
	con.isIdle = true
	con.idleSince = time.Now()
}

// getConnectionById returns connection from pool or creates new connection
func (cp *MqttConnectionPool) connectionByID(connId int) *MqttTransport {
	conn, ok := cp.connPool[connId]
	if ok {
		return conn.mqConnection
	}
	return nil
}

func (cp *MqttConnectionPool) connID() int {
	return int(atomic.AddUint64(&cp.nextID, 1))
}

func (cp *MqttConnectionPool) cleanupProcess() {
	for {
		<-cp.poolCheckTick.C

		if !cp.isStarted.Load() {
			break
		}

		cp.mux.Lock()
		if len(cp.connPool) > cp.size {
			for i := range cp.connPool {
				if !cp.connPool[i].isIdle {
					continue
				}

				if (time.Since(cp.connPool[i].idleSince) > cp.maxIdleAge) && (len(cp.connPool) > cp.size) {
					conn := cp.connectionByID(i)
					if conn != nil {
						conn.Stop()
						delete(cp.connPool, i) // it is safe to delete map element in the loop
					}
				}
			}
		}
		cp.mux.Unlock()
	}

	cp.poolCheckTick.Stop()
	cp.poolCheckTick = nil
}
