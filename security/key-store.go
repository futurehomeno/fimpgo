package security

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

const KeyTypePublic = "pub"
const KeyTypePrivate = "private"
const KeyTypeSymmetric = "sym"

const AlgEcdsa256 = "ES256"

type KeyRecord struct {
	UserId        string
	DeviceId      string // Client device ID (mobile phone id)
	Algorithm     string // ES256
	KeyType       string // public/private/symmetric
	SerializedKey string // serialized key
	EcdsaKey      *EcdsaKey
	AddedAt       string // SerializedKey added timestamp
}

type KeyStore struct {
	keyStore         []KeyRecord
	keyStoreFilePath string
	isPrivate        bool // private store should restrict other application from reading key store file.
}

func NewKeyStore(keyStoreFilePath string,isPrivate bool) *KeyStore {
	if keyStoreFilePath == "" {
		keyStoreFilePath = "/var/lib/futurehome/hub/pub_key_store.json"
	}
	return &KeyStore{keyStoreFilePath: keyStoreFilePath,isPrivate: isPrivate}
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

func (cs *KeyStore) AddSerializedKey(user, device, key, keyType, algo string) error {
	cs.keyStore = append(cs.keyStore, KeyRecord{
		UserId:        user,
		DeviceId:      device,
		SerializedKey: key,
		KeyType:       keyType,
		Algorithm:     algo,
		AddedAt:       time.Now().Format(time.RFC3339),
	})
	return cs.SaveToDisk()
}

func (cs *KeyStore) UpdateSerializedKey(userId, deviceId, key, keyType, algo string) (bool,error) {
	for i := range cs.keyStore {
		if cs.keyStore[i].DeviceId == deviceId && cs.keyStore[i].UserId == userId && cs.keyStore[i].KeyType == keyType && cs.keyStore[i].Algorithm == algo {
			cs.keyStore[i].SerializedKey = key
			cs.keyStore[i].AddedAt = time.Now().Format(time.RFC3339)
			return true,cs.SaveToDisk()

		}
	}
	return false,nil
}

func (cs *KeyStore) UpsertSerializedKey(userId, deviceId, key, keyType, algo string) error {
	keyExists , err := cs.UpdateSerializedKey(userId,deviceId,key,keyType,algo)
	if !keyExists {
		err = cs.AddSerializedKey(userId,deviceId,key,keyType,algo)
	}
	return err
}


func (cs *KeyStore) GetKey(userId, deviceId, keyType string) *KeyRecord {
	for i := range cs.keyStore {
		if cs.keyStore[i].DeviceId == deviceId && cs.keyStore[i].UserId == userId && cs.keyStore[i].KeyType == keyType {
			return &cs.keyStore[i]
		}
	}
	return nil
}

func (cs *KeyStore) GetEcdsaKey(userId, deviceId ,keyType string) (*EcdsaKey,error) {
	for i := range cs.keyStore {
		if cs.keyStore[i].DeviceId == deviceId && cs.keyStore[i].UserId == userId && cs.keyStore[i].Algorithm == AlgEcdsa256 && cs.keyStore[i].KeyType == keyType{
			if cs.keyStore[i].EcdsaKey == nil {
				if cs.keyStore[i].SerializedKey == "" {
					log.Warn("<kstore> Empty key string")
					return nil,fmt.Errorf("empty key string")
				}else {
					cs.keyStore[i].EcdsaKey = NewEcdsaKey()
					var err error
					if cs.keyStore[i].KeyType == KeyTypePrivate {
						err = cs.keyStore[i].EcdsaKey.ImportX509PrivateKey(cs.keyStore[i].SerializedKey)
					}else if cs.keyStore[i].KeyType == KeyTypePublic {
						err = cs.keyStore[i].EcdsaKey.ImportX509PublicKey(cs.keyStore[i].SerializedKey)
					}else {
						return nil,fmt.Errorf("unknown key type %s",keyType)
					}
					if err != nil {
						return nil, err
					}
					return cs.keyStore[i].EcdsaKey,nil
				}
			}else{
				return cs.keyStore[i].EcdsaKey,nil
			}
		}
	}
	return nil,fmt.Errorf("key not found")
}


//
func (cs *KeyStore) GetAllUserKeys(userId string) []KeyRecord {
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
	var mode os.FileMode
	if cs.isPrivate {
		mode = 0600
	}else {
		mode = 0664
	}
	err = ioutil.WriteFile(cs.keyStoreFilePath, bpayload, mode)
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
