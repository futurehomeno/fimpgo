package fimpgo

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
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
	clientIdPrefix string
	initSize       int // init size
	size           int // normal size
	maxSize        int // max size
	maxIdleAge     time.Duration // Defines how long idle connection can stay in the pool before it gets destroyed
	poolCheckTick  *time.Ticker
	isActive       bool
}

func NewMqttConnectionPool(initSize, size, maxSize int,maxAge time.Duration, connTemplate MqttConnectionConfigs, clientIdPrefix string) *MqttConnectionPool {
	if maxAge == 0 {
		maxAge = 30*time.Second  // age in seconds
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
	pool := &MqttConnectionPool{connPool: make(map[int]*connection), connTemplate: connTemplate, clientIdPrefix: clientIdPrefix, initSize: initSize, size: size, maxSize: maxSize,maxIdleAge:maxAge}
	pool.Start()
	return pool
}

func (cp *MqttConnectionPool)Start() {
	if !cp.isActive {
		cp.isActive = true
		cp.poolCheckTick = time.NewTicker(10 * time.Second)
		go cp.cleanupProcess()
	}
}

func (cp *MqttConnectionPool)Stop() {
	cp.isActive = false
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

func (cp *MqttConnectionPool) createConnection() (int, error) {
	connId := cp.genConnId()
	if len(cp.connPool) >= cp.maxSize {
		log.Error("<mq-pool> Too many connections")
		return 0, errors.New("too many connections")
	}
	conf := cp.connTemplate
	conf.ClientID = fmt.Sprintf("%s_%d", cp.clientIdPrefix, connId)
	newConnection := NewMqttTransportFromConfigs(conf)
	err := newConnection.Start()
	cp.connPool[connId] = &connection{
		mqConnection: newConnection,
		isIdle:       false,
		startedAt:    time.Now(),
	}
	log.Debugf("New connection %d created . Pool size = %d",connId,len(cp.connPool))
	return connId, err
}

// getConnectionById returns first available connection from the pool or creates new connection
func (cp *MqttConnectionPool) BorrowConnection() (int, *MqttTransport, error) {
	defer cp.mux.Unlock()
	cp.mux.Lock()
	for i := range cp.connPool {
		if cp.connPool[i].isIdle {
			if cp.connPool[i].mqConnection.Client().IsConnected() {
				cp.connPool[i].isIdle = false
				return i, cp.connPool[i].mqConnection, nil
			}else {
				break
			}
		}
	}
	connId, err := cp.createConnection()
	return connId, cp.getConnectionById(connId), err
}

// ReturnConnection returns connection to pool by setting inUse status to false
func (cp *MqttConnectionPool) ReturnConnection(connId int) {
	defer cp.mux.RUnlock()
	cp.mux.RLock()
	con ,ok := cp.connPool[connId]
	if ok {
		con.mqConnection.UnsubscribeAll()
		con.isIdle = true
		con.idleSince = time.Now()
		log.Debugf("Connection %d returned to pool.",connId)
	}
}

// getConnectionById returns connection from pool or creates new connection
func (cp *MqttConnectionPool) getConnectionById(connId int) *MqttTransport {
	conn, ok := cp.connPool[connId]
	if ok {
		return conn.mqConnection
	}
	return nil
}


func (cp *MqttConnectionPool) genConnId() int {
	rand.Seed(int64(time.Now().Nanosecond()))
	for {
		id := rand.Int()
		if _,ok:=cp.connPool[id];!ok {
			return id
		}
	}
}

func (cp *MqttConnectionPool) cleanupProcess() {
	for {
		<-cp.poolCheckTick.C
		if !cp.isActive {
			break
		}
		cp.mux.Lock()
		if len(cp.connPool)>cp.size {
			for i := range cp.connPool {
				if cp.connPool[i].isIdle {
					if (time.Since(cp.connPool[i].idleSince) > (cp.maxIdleAge)) && (len(cp.connPool)>cp.size) {
						log.Debugf("<conn-pool> Destroying old connection")
						conn := cp.getConnectionById(i)
						conn.Stop()
						delete(cp.connPool,i) // it is safe to delete map element in the loop
					}else {
						//log.Debugf("<conn-pool> Nothing to clean")
					}

				}
			}
		}
		cp.mux.Unlock()
	}
}


