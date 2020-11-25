package fimpgo

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// SyncClient allows sync interaction over async channel.
type SyncClient struct {
	mqttTransport         *MqttTransport
	mqttConnPool          *MqttConnectionPool
	isConnPoolEnabled     bool
	transactionPoolSize   int // Max transaction pool size
	inboundBufferSize     int // Inbound message channel buffer size
	inboundMsgChannel     MessageCh
	inboundChannelName    string
	stopSignalCh          chan bool
	isStartedUsingConnect bool
}

func (sc *SyncClient) SetTransactionPoolSize(transactionPoolSize int) {
	sc.transactionPoolSize = transactionPoolSize
}

func NewSyncClient(mqttTransport *MqttTransport) *SyncClient {
	sc := SyncClient{mqttTransport: mqttTransport}
	sc.transactionPoolSize = 20
	sc.inboundBufferSize = 10
	sc.init()
	return &sc
}

func NewSyncClientV2(mqttTransport *MqttTransport, transactionPoolSize int, inboundBuffSize int) *SyncClient {
	sc := SyncClient{mqttTransport: mqttTransport}
	sc.transactionPoolSize = transactionPoolSize
	sc.inboundBufferSize = inboundBuffSize
	sc.init()
	return &sc
}

// NewSyncClientV3 Creates new sync client either using connections pool internal connection
func NewSyncClientV3(mqttTransport *MqttTransport, connPool *MqttConnectionPool) *SyncClient {
	sc := SyncClient{mqttTransport: mqttTransport, mqttConnPool: connPool, isConnPoolEnabled: true}
	sc.transactionPoolSize = 20
	sc.inboundBufferSize = 10
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
}

func (sc *SyncClient) init() {
	sc.stopSignalCh = make(chan bool)
}

// Connect establishes internal connection to mqtt broker and initializes mqtt
// Should be used if MqttTransport instance is not provided in constructor .
func (sc *SyncClient) Connect(serverURI string, clientID string, username string, password string, cleanSession bool, subQos byte, pubQos byte) error {
	if sc.mqttTransport == nil {
		log.Info("<SyncClient> Connecting to mqtt broker")
		sc.mqttTransport = NewMqttTransport(serverURI, clientID, username, password, cleanSession, subQos, pubQos)
		err := sc.mqttTransport.Start()
		if err != nil {
			log.Error("<SyncClient> Error connecting to broker :", err)
			return err
		}
		sc.isStartedUsingConnect = true
	} else {
		log.Info("<SyncClient> Already connected")
	}

	return nil
}

// Stop has to be invoked to stop message listener
func (sc *SyncClient) Stop() {
	if sc.isStartedUsingConnect {
		sc.mqttTransport.Stop()
	}

}

// AddSubscription has to be invoked before Send methods
func (sc *SyncClient) AddSubscription(topic string) {
	if err := sc.mqttTransport.Subscribe(topic); err != nil {
		log.Error("<SyncClient> error subscribing to topic:", err)
	}
}

func (sc *SyncClient) RemoveSubscription(topic string) {
	if err := sc.mqttTransport.Unsubscribe(topic); err != nil {
		log.Error("<SyncClient> error unsubscribing from topic:", err)
	}
}

// SendFimpWithTopicResponse send message over mqtt and awaits response from responseTopic with responseService and responseMsgType
func (sc *SyncClient) sendFimpWithTopicResponse(topic string, fimpMsg *FimpMessage, responseTopic string, responseService string, responseMsgType string, timeout int64, autoSubscribe bool) (*FimpMessage, error) {
	//log.Debug("Registering request uid = ",fimpMsg.UID)
	var conId int
	var conn *MqttTransport
	var inboundCh = make(MessageCh, 10)
	var responseChannel chan *FimpMessage
	var err error
	var chanName = uuid.NewV4().String()

	defer func() {
		if autoSubscribe && responseTopic != "" && conn != nil {
			if err := conn.Unsubscribe(responseTopic); err != nil {
				log.Error("<SyncClient> error unsubscribing from topic:", err)
			}
		}
		if conn != nil {
			conn.UnregisterChannel(chanName)
			close(inboundCh)
			if sc.isConnPoolEnabled {
				// force unset global prefix
				conn.SetGlobalTopicPrefix("")
				sc.mqttConnPool.ReturnConnection(conId)
			}
		}
	}()

	if sc.isConnPoolEnabled {
		conId, conn, err = sc.mqttConnPool.BorrowConnection()
		if err != nil {
			return nil, err
		}
	} else {
		conn = sc.mqttTransport
	}
	conn.RegisterChannel(chanName, inboundCh)

	responseChannel = sc.startResponseListener(fimpMsg, responseMsgType, responseService, responseTopic, inboundCh, timeout)

	// this if statement is currently dead code, as the autoSubscribe parameter is only called with false
	if autoSubscribe && responseTopic != "" {
		if err := conn.Subscribe(responseTopic); err != nil {
			log.Error("<SyncClient> error subscribing to topic:", err)
		}
	} else if responseTopic != "" {
		if err := conn.Subscribe(responseTopic); err != nil {
			log.Error("<SyncClient> error subscribing to topic:", err)
			return nil, errSubscribe
		}
	}

	// force the global prefix
	conn.SetGlobalTopicPrefix(sc.mqttTransport.getGlobalTopicPrefix())

	if err := conn.PublishToTopic(topic, fimpMsg); err != nil {
		log.Error("<SyncClient> error publishing to topic:", err)
		return nil, errPublish
	}

	select {
	case fimpResponse := <-responseChannel:
		return fimpResponse, nil
	case <-time.After(time.Second * time.Duration(timeout)):
		log.Info("<SyncClient> No response from queue for ", timeout)
		return nil, errTimeout
	}
}

// SendReqRespFimp sends msg to topic and expects to receive response on response topic . If autoSubscribe is set to true , the system will automatically subscribe and unsubscribe from response topic
func (sc *SyncClient) SendReqRespFimp(cmdTopic, responseTopic string, reqMsg *FimpMessage, timeout int64, autoSubscribe bool) (*FimpMessage, error) {
	return sc.sendFimpWithTopicResponse(cmdTopic, reqMsg, responseTopic, "", "", timeout, autoSubscribe)
}

// SendFimp sends message over mqtt and blocks until request is received or timeout is reached .
// messages are correlated using uid->corid
func (sc *SyncClient) SendFimp(topic string, fimpMsg *FimpMessage, timeout int64) (*FimpMessage, error) {
	return sc.SendFimpWithTopicResponse(topic, fimpMsg, "", "", "", timeout)
}

// SendFimpWithTopicResponse send message over mqtt and awaits response from responseTopic with responseService and responseMsgType (the method is for backward compatibility)
func (sc *SyncClient) SendFimpWithTopicResponse(topic string, fimpMsg *FimpMessage, responseTopic string, responseService string, responseMsgType string, timeout int64) (*FimpMessage, error) {
	return sc.sendFimpWithTopicResponse(topic, fimpMsg, responseTopic, responseService, responseMsgType, timeout, false)
}

func (sc *SyncClient) startResponseListener(requestMsg *FimpMessage, respMsgType, respService, respTopic string, inboundCh MessageCh, timeout int64) chan *FimpMessage {
	log.Debug("<SyncClient> Msg listener is started")
	respChan := make(chan *FimpMessage)
	go func() {
		for msg := range inboundCh {
			//log.Debug(msg.Payload.Value)
			if (respMsgType == msg.Payload.Type && respService == msg.Payload.Service && respTopic == msg.Topic) || requestMsg.UID == msg.Payload.CorrelationID {
				//log.Debug("Match")
				select {
				case respChan <- msg.Payload:
				case <-time.After(time.Second * time.Duration(timeout)):
				}
			}
		}
	}()
	return respChan
}
