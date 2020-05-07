package edgeapp

import "github.com/futurehomeno/fimpgo"
import log "github.com/sirupsen/logrus"

type FhOAuth2Client struct {
	hubToken string
	syncClient *fimpgo.SyncClient
	appName string
	partnerName string
	mqt *fimpgo.MqttTransport
	mqttServerURI string
	mqttClientID  string
}

func NewFhOAuth2Client(partnerName string,appName string) *FhOAuth2Client {
	return &FhOAuth2Client{partnerName: partnerName,mqttServerURI: "tcp://localhost:1883",mqttClientID: "fhouthclient"}
}

func (oac *FhOAuth2Client) ConfigureFimpSyncClient() error {
	if oac.mqt == nil {
		oac.mqt = fimpgo.NewMqttTransport(oac.mqttServerURI,oac.mqttClientID,"","",true,1,1)
		err := oac.mqt.Start()
		log.Debug("Auth mqtt client connected")
		if err != nil {
			log.Error("Error connecting to broker ",err)
			return err
		}
	}
	oac.syncClient = fimpgo.NewSyncClient(oac.mqt)
	return nil
}

func (oac *FhOAuth2Client) LoadHubTokenFromCB()error  {
	oac.syncClient.AddSubscription("pt:j1/mt:rsp/rt:app/rn:clbridge/ad:1")

}

func (oac *FhOAuth2Client) ExchangeCodeForTokens(code string)  {

}
func (oac *FhOAuth2Client) ExchangeRefreshToken(refreshToken string) {

}
