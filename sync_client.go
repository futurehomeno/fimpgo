package fimpgo

import (
	"fmt"
	"sync"
	"time"

	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/futurehomeno/fimpgo/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// SyncClient allows sync interaction over async channel.
type SyncClient struct {
	mqttTransportLock   sync.Mutex
	mqttTransport       *MqttTransport
	mqttConnPool        *MqttConnectionPool
	isConnPoolEnabled   bool
	transactionPoolSize int // Max transaction pool size
	inboundBufferSize   int // Inbound message channel buffer size
	stopSignalCh        chan bool
	globalPrefix        string
	mqttStarted         bool
}

// SetGlobalPrefix configures global prefix/site_id . Most be used from backend services.
func (sc *SyncClient) SetGlobalPrefix(globalPrefix string) {
	sc.globalPrefix = globalPrefix
}

func (sc *SyncClient) SetTransactionPoolSize(transactionPoolSize int) {
	sc.transactionPoolSize = transactionPoolSize
}

// NewSyncClient creates sync client using existing mqtt connection
func NewSyncClient(mqttTransport *MqttTransport) *SyncClient {
	sc := SyncClient{mqttTransport: mqttTransport}
	sc.transactionPoolSize = 20
	sc.inboundBufferSize = 10
	sc.init()
	return &sc
}

// NewSyncClientV2 creates new sync client using existing mqtt connection and configures transactionPool size and inboundBufferSize
func NewSyncClientV2(mqttTransport *MqttTransport, transactionPoolSize int, inboundBuffSize int) *SyncClient {
	sc := SyncClient{mqttTransport: mqttTransport}
	sc.transactionPoolSize = transactionPoolSize
	sc.inboundBufferSize = inboundBuffSize
	sc.init()
	return &sc
}

func (sc *SyncClient) SetConfigs(transactionPoolSize int, inboundBuffSize int) {
	if transactionPoolSize == 0 {
		transactionPoolSize = 20
	}

	if inboundBuffSize == 0 {
		inboundBuffSize = 10
	}

	sc.transactionPoolSize = transactionPoolSize
	sc.inboundBufferSize = inboundBuffSize
}

func (sc *SyncClient) init() {
	sc.stopSignalCh = make(chan bool)
}

// Connect establishes internal connection to mqtt broker and initializes mqtt
// Should be used if MqttTransport instance is not provided in constructor .
func (sc *SyncClient) Connect(serverURI string, clientID string, username string, password string, cleanSession bool, subQos byte, pubQos byte, errHandler func(error)) error {
	sc.mqttTransportLock.Lock()
	defer sc.mqttTransportLock.Unlock()

	if sc.mqttTransport == nil {
		sc.mqttTransport = NewMqttTransport(serverURI, clientID, username, password, cleanSession, subQos, pubQos, errHandler)
		err := sc.mqttTransport.Start(10 * time.Second)
		if err != nil {
			return err
		}

		sc.mqttStarted = true
	}

	return nil
}

// Stop has to be invoked to stop message listener
func (sc *SyncClient) Stop() {
	sc.mqttTransportLock.Lock()
	defer sc.mqttTransportLock.Unlock()

	if sc.mqttStarted {
		sc.mqttTransport.Stop()
	}
}

// AddSubscription has to be invoked before Send methods
func (sc *SyncClient) AddSubscription(topic string) error {
	if sc.mqttTransport == nil {
		return fmt.Errorf("not connected")
	}

	return sc.mqttTransport.Subscribe(topic)
}

// RemoveSubscription
func (sc *SyncClient) RemoveSubscription(topic string) error {
	if sc.mqttTransport == nil {
		return fmt.Errorf("not connected")
	}

	return sc.mqttTransport.Unsubscribe(topic)
}

// SendFimpWithTopicResponse send message over mqtt and awaits response from responseTopic with responseService and responseMsgType
func (sc *SyncClient) sendFimpWithTopicResponse(topic string, fimpMsg *FimpMessage, responseTopic string, responseService fimptype.ServiceTypeT, responseMsgType string, timeout int, autoSubscribe bool) (*FimpMessage, error) {
	var conId int
	var conn *MqttTransport
	var inboundCh = make(MessageCh, 10)
	var err error
	var chanName = uuid.New().String()

	defer func() {
		if conn == nil {
			return
		}

		if autoSubscribe && responseTopic != "" {
			if err := conn.Unsubscribe(responseTopic); err != nil {
				log.Error("[fimpgo] Error unsubscribing from topic:", err)
			}
		}

		conn.UnregisterChannel(chanName)
		close(inboundCh)

		if sc.isConnPoolEnabled {
			// force unset global prefix
			conn.SetGlobalTopicPrefix("")
			sc.mqttConnPool.ReturnConnection(conId)
		}
	}()

	if sc.isConnPoolEnabled {
		conId, conn, err = sc.mqttConnPool.BorrowConnection(nil)
		if err != nil {
			return nil, err
		}
	} else {
		conn = sc.mqttTransport
	}

	if conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	conn.RegisterChannel(chanName, inboundCh)

	responseChannel := sc.startResponseListener(fimpMsg, responseMsgType, responseService, responseTopic, inboundCh, timeout)

	// force the global prefix -> this is useful for per-site operations
	if sc.globalPrefix != "" {
		conn.SetGlobalTopicPrefix(sc.globalPrefix)
	}

	if autoSubscribe && responseTopic != "" {
		if err = conn.Subscribe(responseTopic); err != nil {
			return nil, fmt.Errorf("subscribe err: %w", err)
		}
	}

	if err = conn.PublishToTopic(topic, fimpMsg); err != nil {
		return nil, fmt.Errorf("publish err: %w", err)
	}

	select {
	case fimpResponse := <-responseChannel:
		return fimpResponse, nil
	case <-time.After(time.Second * time.Duration(timeout)):
		return nil, utils.ErrTimeout
	}
}

// SendReqRespFimp sends msg to topic and expects to receive response on response topic . If autoSubscribe is set to true , the system will automatically subscribe and unsubscribe from response topic.
func (sc *SyncClient) SendReqRespFimp(cmdTopic, responseTopic string, reqMsg *FimpMessage, timeout int, autoSubscribe bool) (*FimpMessage, error) {
	return sc.sendFimpWithTopicResponse(cmdTopic, reqMsg, responseTopic, "", "", timeout, autoSubscribe)
}

// SendFimp sends message over mqtt and blocks until request is received or timeout is reached .
// messages are correlated using uid->corid
func (sc *SyncClient) SendFimp(topic string, fimpMsg *FimpMessage, timeout int) (*FimpMessage, error) {
	return sc.SendFimpWithTopicResponse(topic, fimpMsg, "", "", "", timeout)
}

// SendFimpWithTopicResponse send message over mqtt and awaits response from responseTopic with responseService and responseMsgType (the method is for backward compatibility)
func (sc *SyncClient) SendFimpWithTopicResponse(topic string, fimpMsg *FimpMessage, responseTopic string, responseService fimptype.ServiceTypeT, responseMsgType string, timeout int) (*FimpMessage, error) {
	return sc.sendFimpWithTopicResponse(topic, fimpMsg, responseTopic, responseService, responseMsgType, timeout, false)
}

// startResponseListener starts response listener , it blocks callers proc until response is received or timeout.
func (sc *SyncClient) startResponseListener(requestMsg *FimpMessage, respMsgType string, respService fimptype.ServiceTypeT, respTopic string, inboundCh MessageCh, timeout int) chan *FimpMessage {
	respChan := make(chan *FimpMessage)

	go func() {
		for msg := range inboundCh {
			if (respMsgType == msg.Payload.Type && respService == msg.Payload.Service && respTopic == msg.Topic) || requestMsg.UID == msg.Payload.CorrelationID {
				select {
				case respChan <- msg.Payload:
				case <-time.After(time.Second * time.Duration(timeout)):
				}
				return
			}
		}
	}()
	return respChan
}
