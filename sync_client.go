package fimpgo

import (
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
	"errors"
	"github.com/satori/go.uuid"
)

type transaction struct {
	respTopic   string
	respService string
	respMsgType string
	requestUid  string
	isActive    bool
	respChannel chan *FimpMessage
}

// SyncClient allows sync interaction over async channel.
type SyncClient struct {
	 mqttTransport       *MqttTransport
	 transactions        []transaction
	 mux                 sync.Mutex
	 transactionPoolSize int // Max transaction pool size
	 inboundBufferSize   int // Inbound message channel buffer size
	 inboundMsgChannel MessageCh
	 inboundChannelName string
	 stopSignalCh chan bool
	 isStartedUsingConnect bool
}

func (sc *SyncClient) SetTransactionPoolSize(transactionPoolSize int) {
	sc.transactionPoolSize = transactionPoolSize
}

func NewSyncClient(mqttTransport *MqttTransport) *SyncClient {
	sc := SyncClient{mqttTransport:mqttTransport}
	sc.transactionPoolSize = 20
	sc.inboundBufferSize = 10
	sc.init()
	return &sc
}

func NewSyncClientV2(mqttTransport *MqttTransport,transactionPoolSize int , inboundBuffSize int) *SyncClient {
	sc := SyncClient{mqttTransport:mqttTransport}
	sc.transactionPoolSize = transactionPoolSize
	sc.inboundBufferSize = inboundBuffSize
	sc.init()
	return &sc
}

func (sc *SyncClient) init() {
	sc.stopSignalCh = make (chan bool)
	sc.inboundMsgChannel = make (MessageCh,sc.inboundBufferSize)
	sc.inboundChannelName = uuid.NewV4().String()
	sc.transactions = []transaction{}
	if sc.mqttTransport != nil {
		sc.mqttTransport.RegisterChannel(sc.inboundChannelName,sc.inboundMsgChannel)
	}
	go sc.messageListener()
}

// Connect establishes internal connection to mqtt broker and initializes mqtt
// Should be used if MqttTransport instance is not provided in constructor .
func (sc *SyncClient) Connect(serverURI string, clientID string, username string, password string, cleanSession bool, subQos byte, pubQos byte) error {
	if sc.mqttTransport == nil {
		log.Info("<SyncClient> Connecting to mqtt broker")
		sc.mqttTransport = NewMqttTransport(serverURI,clientID,username,password,cleanSession,subQos,pubQos)
		err := sc.mqttTransport.Start()
		if err != nil {
			log.Error("<SyncClient> Error connecting to broker :",err)
			return err
		}
		sc.isStartedUsingConnect = true
		sc.mqttTransport.RegisterChannel(sc.inboundChannelName,sc.inboundMsgChannel)
	}else {
		log.Info("<SyncClient> Already connected")
	}

	return nil
}

// Stop has to be invoked to stop message listener
func (sc *SyncClient) Stop() {
	sc.mqttTransport.UnregisterChannel(sc.inboundChannelName)
	close(sc.inboundMsgChannel)
	if sc.isStartedUsingConnect {
		sc.mqttTransport.Stop()
	}

}
// AddSubscription has to be invoked before Send methods
func (sc *SyncClient) AddSubscription(topic string) {
	sc.mqttTransport.Subscribe(topic)
}

func (sc *SyncClient) RemoveSubscription(topic string) {
	sc.mqttTransport.Unsubscribe(topic)
}

// SendFimp sends message over mqtt and blocks until request is received or timeout is reached .
// meesages are corelated using uid->corid
func (sc *SyncClient) SendFimp(topic string, fimpMsg *FimpMessage,timeout int64) (*FimpMessage,error) {
	return sc.SendFimpWithTopicResponse(topic,fimpMsg,"","","",timeout)
}

// SendFimpWithTopicResponse send message over mqtt and awaits response from responseTopic with responseService and responseMsgType
func (sc *SyncClient) SendFimpWithTopicResponse(topic string, fimpMsg *FimpMessage,responseTopic string,responseService string,responseMsgType string,timeout int64) (*FimpMessage,error) {

	msgB , err := fimpMsg.SerializeToJson()
	if err != nil {
		return nil,err
	}
	responseChannel := sc.registerRequest(fimpMsg.UID,responseTopic,responseService,responseMsgType)
	sc.mqttTransport.PublishRaw(topic,msgB)
	select {
	case fimpResponse := <- responseChannel :
		sc.unregisterRequest(fimpMsg.UID,responseTopic,responseService,responseMsgType)
		return fimpResponse,nil

	case <- time.After(time.Second* time.Duration(timeout)):
		log.Info("<SyncClient> No response from queue for ",timeout)


	}
	sc.unregisterRequest(fimpMsg.UID,responseTopic,responseService,responseMsgType)
	return nil, errors.New("request timed out")
}


func (sc *SyncClient) registerRequest(responseUid string,responseTopic string , responseService string, responseMsgType string) chan *FimpMessage {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	// searching for non-active transaction in pool
	for i := range sc.transactions {
		if sc.transactions[i].isActive == false {
			sc.transactions[i].isActive = true
			sc.transactions[i].requestUid = responseUid
			sc.transactions[i].respService = responseService
			sc.transactions[i].respMsgType = responseMsgType
			sc.transactions[i].respTopic = responseTopic
			//log.Info("<SyncClient> Transaction from pool")
			return sc.transactions[i].respChannel
		}
	}
	// no active transactions , let's create one
	respChan := make(chan *FimpMessage)
	runReq := transaction{respTopic:responseTopic,respService:responseService,respMsgType:responseMsgType,respChannel:respChan, requestUid:responseUid}
	sc.transactions = append(sc.transactions,runReq)
	//log.Info("<SyncClient> Transaction was added")
	return respChan
}


func (sc *SyncClient) unregisterRequest(responseUid string,responseTopic string , responseService string, responseMsgType string) {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	var result []transaction
	for i := range sc.transactions {
		if (sc.transactions[i].respTopic == responseTopic &&
		   sc.transactions[i].respMsgType == responseMsgType &&
		   sc.transactions[i].respService == responseService) || sc.transactions[i].requestUid == responseUid {

		   if len(sc.transactions)> sc.transactionPoolSize {
			   log.Debugf("<SyncClient> Removing transaction from pool")
			   result = append(sc.transactions[:i],sc.transactions[i+1:]...)
			   break
		   }
			sc.transactions[i].respTopic = ""
			sc.transactions[i].respMsgType = ""
			sc.transactions[i].respService = ""
			sc.transactions[i].requestUid = ""
			sc.transactions[i].isActive = false
		}else {
			log.Debug("<SyncClient> Nothing to unregister")
		}
	}
	if result != nil {
		sc.transactions = result
	}

	//log.Info("Unregister .Number of in flight transactions = ",len(pr.RunningRequest))
}


// OnMessage is invoked by an adapter on every new message
// The code is executed in callers goroutine
func (sc *SyncClient) messageListener() {
	log.Debug("<SyncClient> Msg listener is started")
	for msg := range sc.inboundMsgChannel{
			sc.mux.Lock()
			for i := range sc.transactions {
				if (sc.transactions[i].respMsgType == msg.Payload.Type && sc.transactions[i].respService == msg.Payload.Service && sc.transactions[i].respTopic == msg.Topic) ||
					sc.transactions[i].requestUid == msg.Payload.CorrelationID{
					//log.Debug("<SyncClient> Transaction match , transaction size = ",len(sc.transactions))
					//sending message to coresponding channel
					select {
					case sc.transactions[i].respChannel <- msg.Payload:
					default:
						log.Error("<SyncClient> No channel to send the message.")

					}
				}
			}
			sc.mux.Unlock()
		}
	log.Debug("Stopping Message listener")

}


