package fimpgo

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type connection struct {
	mqConnection *MqttTransport
	isInUse      bool
	startTime    time.Time
}

type MqttConnectionPool struct {
	mux            sync.Mutex
	connTemplate   MqttConnectionConfigs
	connPool       map[int]*connection
	clientIdPrefix string
	initSize       int
	size           int
	maxSize        int
}

//TODO : clean up old connections

func NewMqttConnectionPool(initSize, size, maxSize int, connTemplate MqttConnectionConfigs, clientIdPrefix string) *MqttConnectionPool {
	return &MqttConnectionPool{connPool: make(map[int]*connection), connTemplate: connTemplate, clientIdPrefix: clientIdPrefix, initSize: initSize, size: size, maxSize: maxSize}
}

// AddConnection adds existing connection to the pool
func (cp *MqttConnectionPool) AddConnection(connId int, conn *MqttTransport) int {
	defer cp.mux.Unlock()
	cp.mux.Lock()
	if connId == 0 {
		connId = cp.genConnId()
	}
	cp.connPool[connId] = &connection{
		mqConnection: conn,
		isInUse:      true,
		startTime:    time.Now(),
	}
	return connId
}

func (cp *MqttConnectionPool) createConnection() (int, error) {
	connId := cp.genConnId()
	if len(cp.connPool) >= cp.maxSize {
		log.Error("<mq-pool> Too many connections")
		return 0, errors.New("too many connections")
	}
	conf := cp.connTemplate
	conf.ClientID = fmt.Sprintf("%s-%d", cp.clientIdPrefix, connId)
	newConnection := NewMqttTransportFromConfigs(conf)
	cp.connPool[connId] = &connection{
		mqConnection: newConnection,
		isInUse:      true,
		startTime:    time.Now(),
	}
	err := cp.connPool[connId].mqConnection.Start()
	return connId, err
}

// GetConnectionById returns first available connection from the pool or creates new connection
func (cp *MqttConnectionPool) GetConnection() (int, *MqttTransport, error) {
	defer cp.mux.Unlock()
	cp.mux.Lock()
	for i := range cp.connPool {
		if !cp.connPool[i].isInUse {
			return i, cp.connPool[i].mqConnection, nil
		}
	}
	connId, err := cp.createConnection()
	return connId, cp.GetConnectionById(connId), err
}

// GetConnectionById returns connection from pool or creates new connection
func (cp *MqttConnectionPool) GetConnectionById(connId int) *MqttTransport {
	conn, ok := cp.connPool[connId]
	if ok {
		return conn.mqConnection
	}
	return nil
}

// ReturnConnection returns connection to pool by setting inUse status to false
func (cp *MqttConnectionPool) ReturnConnection(connId int) {
	defer cp.mux.Unlock()
	cp.mux.Lock()
	con ,ok := cp.connPool[connId]
	if ok {
		log.Debugf("Connection %d returned to pool",connId)
		con.isInUse = false
	}
}

func (cp *MqttConnectionPool) genConnId() int {
	return len(cp.connPool) + 1
}


