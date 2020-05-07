package integration

import (
	"encoding/json"
	"fmt"
	"github.com/futurehomeno/fimpgo"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"time"
)

const KeyTypePublic = "pub"

type KeyRecord struct {
	UserId    string
	DeviceId  string
	Algorithm string
	KeyType   string
	Key       string
	nonce     string // last nonce sent to remote client
	AddedAt   string
}

// The message is sent to remote client as request for public key
type AppKeyRequest struct {
	Nonce   string `json:"nonce"`
	KeyType string `json:"key_type"`
	Algorithm  string `json:"algo"`
}

// Remote client mus generate key-pair and respond with the message
type AppKeyResponse struct {
	SignedNonce string `json:"signed_nonce"`
	Key         string `json:"key"`
	KeyType     string `json:"key_type"`
	Algorithm   string `json:"algo"`
	UserId      string `json:"user_id"`
	DeviceId    string `json:"device_id"`
}

type KeyStore struct {
	keyStore         []KeyRecord
	keyStoreFilePath string
}

// full username is a string which identifies a user on given device
func (cs *KeyStore) CheckIfUserHasKey(userId, deviceId string) bool {
	for i := range cs.keyStore {
		if cs.keyStore[i].DeviceId == deviceId && cs.keyStore[i].UserId == userId {
			return true
		}
	}
	return false
}
// Requesting crypto key from remote client (mobile app).
func (cs *KeyStore) RequestKeyFromRemoteClient(mqtt *fimpgo.MqttTransport) error {
	// send request message with nonce
	// save response to internal store
	// Edge app ---(evt.auth.get_client_key)--> Mobile APP
	//    /|\                                    |
	//     |_____(cmd.auth.set_client_key)_______|
	//
	edgeAppTopic := "edge_app_topic"
	requestTopic := "request_topic"
	request := AppKeyRequest{
		Nonce: cs.getNonce(),
		KeyType: KeyTypePublic,
	}
	reqMsg := fimpgo.NewObjectMessage("evt.auth.get_client_key","mobile-app",request,nil,nil,nil)
	syncClient := fimpgo.NewSyncClient(mqtt)
	syncClient.AddSubscription(edgeAppTopic)
	response,err := syncClient.SendFimp(requestTopic,reqMsg,5)
	syncClient.RemoveSubscription(edgeAppTopic)
	if err != nil {
		log.Error("<key-man> Key request timed out . Err:",err.Error())
		return err
	}
	resp := AppKeyResponse{}
	err = response.GetObjectValue(&resp)
	if err != nil {
		log.Error("<key-man> Key response can't be mapped to response object . Err:",err.Error())
		return err
	}
	return cs.AddKey(resp.UserId,resp.DeviceId,resp.Key,resp.KeyType,resp.Algorithm)
}

func (cs *KeyStore) getNonce() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return fmt.Sprint(r1.Int31())
}

func (cs *KeyStore) AddKey(user, device, key,keyType, cipher string) error {
	cs.keyStore = append(cs.keyStore, KeyRecord{
		UserId:    user,
		DeviceId:  device,
		Key:       key,
		KeyType:   keyType,
		Algorithm: cipher,
		AddedAt:   time.Now().Format(time.RFC3339),
	})
	return cs.SaveToDisk()
}

func (cs *KeyStore) GetKey(userId, deviceId, keyType string) *KeyRecord {
	for i := range cs.keyStore {
		if cs.keyStore[i].DeviceId == deviceId && cs.keyStore[i].UserId == userId && cs.keyStore[i].KeyType == keyType {
			return &cs.keyStore[i]
		}
	}
	return nil
}

//
func (cs *KeyStore) GetAllUserKeys(userId, keyType string) []KeyRecord {
	var result []KeyRecord
	for i := range cs.keyStore {
		if cs.keyStore[i].UserId == userId {
			result = append(result, cs.keyStore[i])
		}
	}
	return result
}

func (cs *KeyStore) SaveToDisk() error {
	bpayload, err := json.Marshal(cs.keyStore)
	err = ioutil.WriteFile(cs.keyStoreFilePath, bpayload, 0664)
	if err != nil {
		return err
	}
	return err
	return nil
}

func (cs *KeyStore) LoadFromDisk() error {
	configFileBody, err := ioutil.ReadFile(cs.keyStoreFilePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configFileBody, &cs.keyStore)
	if err != nil {
		return err
	}
	return nil
}

// An app should call the method to authenticate message
//func (cs *KeyStore) IsMessageAuthenticated(msg *fimpgo.FimpMessage) (bool) {
//	//1. Extract username and signature from user message
//	//2. Query public key from local key store
//	//3. Validate signature using public key
//
//}

