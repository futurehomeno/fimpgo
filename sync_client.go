package fimpgo

import (
	"sync"
	"time"
	log "github.com/Sirupsen/logrus"
	"errors"
)

type transaction struct {
	respTopic string
	respService string
	respMsgType string
	respUid string
	isActive bool
	respChannel chan *FimpMessage
}

// SyncClient allows sync interaction over async channel.
type SyncClient struct {
	 mqttTransport       *MqttTransport
	 transactions        []transaction
	 mux                 sync.Mutex
	 transactionPoolSize int
	 inboundMsgChannel MessageCh
	 stopSignalCh chan bool
}

func NewSyncClient(mqttTransport *MqttTransport) *SyncClient {
	sc := SyncClient{mqttTransport:mqttTransport}
	sc.init()
	return &sc
}

func (sc *SyncClient) init() {
	sc.stopSignalCh = make (chan bool)
	sc.inboundMsgChannel = make (MessageCh)
	sc.transactions = []transaction{}
	go sc.MessageListener()
}

// Connect establishes connection to mqtt broker and initializes mqtt
// Should be used if MqttTransport instance is not provided in constructor .
func (sc *SyncClient) Connect(serverURI string, clientID string, username string, password string, cleanSession bool, subQos byte, pubQos byte) error {
	sc.mqttTransport = NewMqttTransport(serverURI,clientID,username,password,cleanSession,subQos,pubQos)
	err := sc.mqttTransport.Start()
	if err != nil {
		log.Error("Error connecting to broker ",err)
		return err
	}
	return nil
}

func (sc *SyncClient) Stop() {
	sc.stopSignalCh <- true;

}

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

func (sc *SyncClient) SendFimpWithTopicResponse(topic string, fimpMsg *FimpMessage,responseTopic string,responseService string,responseMsgType string,timeout int64) (*FimpMessage,error) {

	msgB , err := fimpMsg.SerializeToJson()
	if err != nil {
		return nil,err
	}
	responseChannel := sc.registerRequest(fimpMsg.UID,responseTopic,responseService,responseMsgType)
	sc.mqttTransport.PublishRaw(topic,msgB)
	select {
	case fimpResponse := <- responseChannel :
		return fimpResponse,nil

	case <- time.After(time.Second* time.Duration(timeout)):
		log.Info("No response from queue for 10 seconds")


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
			sc.transactions[i].respUid = responseUid
			sc.transactions[i].respService = responseService
			sc.transactions[i].respMsgType = responseMsgType
			sc.transactions[i].respTopic = responseTopic
			return sc.transactions[i].respChannel
		}
	}
	// no active transactions , let's create one
	respChan := make(chan *FimpMessage)
	runReq := transaction{respTopic:responseTopic,respMsgType:responseMsgType,respChannel:respChan,respUid:responseUid}
	sc.transactions = append(sc.transactions,runReq)
	return respChan
}


func (sc *SyncClient) unregisterRequest(responseUid string,responseTopic string , responseService string, responseMsgType string) {
	sc.mux.Lock()
	defer sc.mux.Unlock()
	var result []transaction
	for i := range sc.transactions {
		if (sc.transactions[i].respTopic == responseTopic &&
		   sc.transactions[i].respMsgType == responseMsgType &&
		   sc.transactions[i].respService == responseService) || sc.transactions[i].respUid == responseUid {

		   if len(sc.transactions)> sc.transactionPoolSize {
			   result = append(sc.transactions[:i],sc.transactions[i+1:]...)
			   break
		   }
			sc.transactions[i].respTopic = ""
			sc.transactions[i].respMsgType = ""
			sc.transactions[i].respService = ""
			sc.transactions[i].respUid = ""
			sc.transactions[i].isActive = false
		}
	}
	if result != nil {
		sc.transactions = result
	}

	//log.Info("Unregister .Number of in flight transactions = ",len(pr.RunningRequest))
}


// OnMessage is invoked by an adapter on every new message
// The code is executed in callers goroutine
func (sc *SyncClient) MessageListener() {
	var msg *Message
	select {
		case msg =<- sc.inboundMsgChannel:
			log.Debugf("New message from topic %s",msg.Topic)
			for i := range sc.transactions {
				if (sc.transactions[i].respMsgType == msg.Payload.Type && sc.transactions[i].respService == msg.Payload.Service && sc.transactions[i].respTopic == msg.Topic) ||
					sc.transactions[i].respUid == msg.Payload.CorrelationID{
					//sending message to coresponding channel
					sc.transactions[i].respChannel <- msg.Payload
				}
			}
		case _ =<-sc.stopSignalCh:
			log.Info("Stopping Message listener")
	}


}


