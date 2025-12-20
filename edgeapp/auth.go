package edgeapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/utils"
	log "github.com/sirupsen/logrus"
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

func (oac *FhOAuth2Client) SetHubToken(hubToken string) {
	oac.hubToken = hubToken
}

func (oac *FhOAuth2Client) AuthCodeApiUrl() string {
	return oac.authCodeApiUrl
}

func (oac *FhOAuth2Client) SetAuthCodeApiUrl(authCodeApiUrl string) {
	oac.authCodeApiUrl = authCodeApiUrl
}

func (oac *FhOAuth2Client) RefreshTokenApiUrl() string {
	return oac.refreshTokenApiUrl
}

func (oac *FhOAuth2Client) SetRefreshTokenApiUrl(refreshTokenApiUrl string) {
	oac.refreshTokenApiUrl = refreshTokenApiUrl
}

// NewFhOAuth2Client implements OAuth client which communicates to 3rd party API over FH Auth proxy.
func NewFhOAuth2Client(partnerName string, appName string, env string) *FhOAuth2Client {
	client := &FhOAuth2Client{partnerName: partnerName, mqttServerURI: "tcp://127.0.0.1:1883", mqttClientID: "auth_client_" + appName}
	if env == utils.EnvBeta {
		client.refreshTokenApiUrl = "https://partners-beta.futurehome.io/api/control/edge/proxy/refresh"
		client.authCodeApiUrl = "https://partners-beta.futurehome.io/api/control/edge/proxy/auth-code"
	} else {
		client.refreshTokenApiUrl = "https://partners.futurehome.io/api/control/edge/proxy/refresh"
		client.authCodeApiUrl = "https://partners.futurehome.io/api/control/edge/proxy/auth-code"
	}
	client.retryDelay = 60
	client.refreshRetry = 5
	client.cbRetryDelay = 30
	client.cbRetry = 7
	client.appName = appName
	return client
}

// Init has to be invoked before requesting access token
func (oac *FhOAuth2Client) Init() error {
	return oac.LoadHubTokenFromCB()
}

// SetParameters can be used to change default configuration parameter parameters. Parameters which are set to null values will be ignored
func (oac *FhOAuth2Client) SetParameters(mqttServerUri, authCodeApiUrl, refreshTokenApiUrl string, retryDelay time.Duration, refreshRetry int, cbRetry int, cbRetryDelay time.Duration) {
	if mqttServerUri != "" {
		oac.mqttServerURI = mqttServerUri
	}
	if authCodeApiUrl != "" {
		oac.authCodeApiUrl = authCodeApiUrl
	}
	if refreshTokenApiUrl != "" {
		oac.refreshTokenApiUrl = refreshTokenApiUrl
	}
	if retryDelay != 0 {
		oac.retryDelay = retryDelay
	}
	if refreshRetry != 0 {
		oac.refreshRetry = refreshRetry
	}
	if cbRetry != 0 {
		oac.cbRetry = cbRetry
	}
	if cbRetryDelay != 0 {
		oac.cbRetryDelay = cbRetryDelay
	}
}

// ConfigureFimpSyncClient configures fimp sync client , which is used to obtain Hub token from cloud bridge.
func (oac *FhOAuth2Client) ConfigureFimpSyncClient() error {
	if oac.mqt == nil {
		oac.mqt = fimpgo.NewMqttTransport(oac.mqttServerURI, oac.mqttClientID, "", "", true, 1, 1)
		err := oac.mqt.Start()
		if err != nil {
			log.Error("Error connecting to broker ", err)
			return err
		}
		log.Debug("Auth mqtt client connected")
		oac.syncClient = fimpgo.NewSyncClient(oac.mqt)
	} else {
		log.Error("Mqtt client is not configured")
	}
	return nil
}

// LoadHubTokenFromCB - requests hub token from CloudBridge
func (oac *FhOAuth2Client) LoadHubTokenFromCB() error {
	if oac.mqt == nil || oac.syncClient == nil {
		if err := oac.ConfigureFimpSyncClient(); err != nil {
			log.Error(err)
		}
	}
	responseTopic := fmt.Sprintf("pt:j1/mt:rsp/rt:app/rn:%s/ad:1", oac.appName)
	if err := oac.syncClient.AddSubscription(responseTopic); err != nil {
		return err
	}

	reqMsg := fimpgo.NewStringMessage("cmd.clbridge.get_auth_token", "clbridge", "", nil, nil, nil)
	reqMsg.ResponseToTopic = responseTopic
	var err error
	var response *fimpgo.FimpMessage
	for range oac.cbRetry {
		response, err = oac.syncClient.SendFimp("pt:j1/mt:cmd/rt:app/rn:clbridge/ad:1", reqMsg, 5)
		if err == nil {
			break
		}
		log.Error("[edgeapp] CB is not responding.Retrying")
		time.Sleep(time.Second * oac.cbRetryDelay)
	}

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

// ExchangeCodeForTokens - exchanging code for access token
func (oac *FhOAuth2Client) ExchangeCodeForTokens(code string) (*OAuth2TokenResponse, error) {
	req := OAuth2AuthCodeProxyRequest{AuthCode: code, PartnerCode: oac.partnerName}
	return oac.postMsg(req, oac.refreshTokenApiUrl)
}

// ExchangeRefreshToken - exchange refresh token for new access
func (oac *FhOAuth2Client) ExchangeRefreshToken(refreshToken string) (*OAuth2TokenResponse, error) {
	req := OAuth2RefreshProxyRequest{RefreshToken: refreshToken, PartnerCode: oac.partnerName}
	return oac.postMsg(req, oac.refreshTokenApiUrl)
}

func (oac *FhOAuth2Client) postMsg(req interface{}, url string) (*OAuth2TokenResponse, error) {
	if oac.hubToken == "" {
		err := oac.LoadHubTokenFromCB()
		if err != nil {
			return nil, errors.New("empty hub token.operation aborted")
		}
	}
	reqB, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: time.Second * 60}
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqB))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+oac.hubToken)

	var resp *http.Response
	for range oac.refreshRetry {
		resp, err = client.Do(r)
		if err == nil && resp.StatusCode < 400 {
			break
		}
		log.Error("[fimpgo] Error response from auth endpoint.Retrying...")
		time.Sleep(time.Second * oac.retryDelay)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error %s response from server", resp.Status)
	}

	bData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tResp := &OAuth2TokenResponse{}
	err = json.Unmarshal(bData, tResp)
	if err != nil {
		return nil, err
	}
	return tResp, nil
}
