package edgeapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type OAuth2TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	RefreshToken string      `json:"refresh_token"`
	Scope        interface{} `json:"scope"`
}

type OAuth2RefreshProxyRequest struct {
	RefreshToken string `json:"refreshToken"`
	PartnerCode  string `json:"partnerCode"`
}

type OAuth2AuthCodeProxyRequest struct {
	AuthCode    string `json:"code"`
	PartnerCode string `json:"partnerCode"`
}

type OAuth2PasswordProxyRequest struct {
	PartnerCode string `json:"partnerCode"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type FhOAuth2Client struct {
	hubToken           string
	syncClient         *fimpgo.SyncClient
	appName            string
	partnerName        string
	mqt                *fimpgo.MqttTransport
	mqttServerURI      string
	mqttClientID       string
	refreshTokenApiUrl string
	authCodeApiUrl     string
	refreshRetry       int
	retryDelay         time.Duration // delay in seconds
	cbRetry            int
	cbRetryDelay       time.Duration
}
//NewFhOAuth2Client implements OAuth client which communicates to 3rd party API over FH Auth proxy.
func NewFhOAuth2Client(partnerName string, appName string) (*FhOAuth2Client,error) {
	client := &FhOAuth2Client{partnerName: partnerName, mqttServerURI: "tcp://localhost:1883", mqttClientID: "fhouthclient"}
	client.retryDelay = 60
	client.refreshRetry = 5
	client.cbRetryDelay = 30
	client.cbRetry = 7
	err := client.ConfigureFimpSyncClient()
	if err != nil {
		return nil, err
	}
	err = client.LoadHubTokenFromCB()
	return client,err
}

func (oac *FhOAuth2Client) ConfigureFimpSyncClient() error {
	if oac.mqt == nil {
		oac.mqt = fimpgo.NewMqttTransport(oac.mqttServerURI, oac.mqttClientID, "", "", true, 1, 1)
		err := oac.mqt.Start()
		log.Debug("Auth mqtt client connected")
		if err != nil {
			log.Error("Error connecting to broker ", err)
			return err
		}
		oac.syncClient = fimpgo.NewSyncClient(oac.mqt)
	} else {
		log.Error("Mqtt client is not configured")
	}
	return nil
}

func (oac *FhOAuth2Client) LoadHubTokenFromCB() error {
	if oac.mqt == nil || oac.syncClient == nil {
		oac.ConfigureFimpSyncClient()
	}
	responseTopic := fmt.Sprintf("pt:j1/mt:rsp/rt:app/rn:%s/ad:1", oac.appName)
	oac.syncClient.AddSubscription(responseTopic)
	reqMsg := fimpgo.NewStringMessage("cmd.clbridge.get_auth_token", "clbridge", "", nil, nil, nil)
	reqMsg.ResponseToTopic = responseTopic
	var err error
	var response *fimpgo.FimpMessage
	for i:=0;i<oac.cbRetry;i++ {
		response, err = oac.syncClient.SendFimp("pt:j1/mt:cmd/rt:app/rn:clbridge/ad:1", reqMsg, 5)
		if err == nil {
			break
		}
		log.Error("CB is not responding.Retrying")
		time.Sleep(time.Second*oac.cbRetryDelay)
	}

	// TODO:retry
	oac.syncClient.Stop()
	oac.mqt.Stop()
	if err != nil {
		return err
	}
	if response.Type != "evt.clbridge.auth_token_report" {
		return errors.New("wrong response msg type")
	}
	oac.hubToken, err = response.GetStringValue()
	return err
}

func (oac *FhOAuth2Client) ExchangeCodeForTokens(code string) (*OAuth2TokenResponse,error) {
	req := OAuth2AuthCodeProxyRequest{AuthCode: code,PartnerCode: oac.partnerName}
	return oac.postMsg(req,oac.refreshTokenApiUrl)
}

func (oac *FhOAuth2Client) ExchangeRefreshToken(refreshToken string) (*OAuth2TokenResponse,error) {
	req := OAuth2RefreshProxyRequest{RefreshToken: refreshToken,PartnerCode: oac.partnerName}
	return oac.postMsg(req,oac.refreshTokenApiUrl)
}

func (oac *FhOAuth2Client) postMsg(req interface{},url string) (*OAuth2TokenResponse,error) {
	if oac.hubToken == "" {
		log.Info("Empty token.Re-requesting new token")
		err := oac.LoadHubTokenFromCB()
		if err != nil {
			return nil,errors.New("empty hub token.operation aborted")
		}
	}
	reqB,err  := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: time.Second * 60}
	r, _ := http.NewRequest("POST", url,bytes.NewBuffer(reqB))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+oac.hubToken)
	//log.Info("Sending using token :",oac.hubToken
	var resp *http.Response
	for i:=0;i<oac.refreshRetry;i++ {
		resp, err = client.Do(r)
		if err == nil &&  resp.StatusCode < 400 {
			break
		}
		log.Error("Error response from auth endpoint.Retrying...")
		time.Sleep(time.Second*oac.retryDelay)
	}
	if err != nil {
		return nil,err
	}
	if resp.StatusCode >= 400 {
		return nil,fmt.Errorf("error %s response from server",resp.Status)
	}

	bData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	tResp := &OAuth2TokenResponse{}
	err = json.Unmarshal(bData,tResp)
	if err != nil {
		return nil,err
	}
	return tResp,nil
}

