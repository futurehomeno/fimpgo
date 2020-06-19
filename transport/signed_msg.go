package transport

import (
	"crypto"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/security"
	"strings"
)

/*
https://tools.ietf.org/html/rfc7518#section-3
+--------------+-------------------------------+--------------------+
| "alg" Param  | Digital Signature or MAC      | Implementation     |
| Value        | Algorithm                     | Requirements       |
+--------------+-------------------------------+--------------------+
| HS256        | HMAC using SHA-256            | Required           |
| HS384        | HMAC using SHA-384            | Optional           |
| HS512        | HMAC using SHA-512            | Optional           |
| RS256        | RSASSA-PKCS1-v1_5 using       | Recommended        |
|              | SHA-256                       |                    |
| RS384        | RSASSA-PKCS1-v1_5 using       | Optional           |
|              | SHA-384                       |                    |
| RS512        | RSASSA-PKCS1-v1_5 using       | Optional           |
|              | SHA-512                       |                    |
| ES256        | ECDSA using P-256 and SHA-256 | Recommended+       |
| ES384        | ECDSA using P-384 and SHA-384 | Optional           |
| ES512        | ECDSA using P-521 and SHA-512 | Optional           |
| PS256        | RSASSA-PSS using SHA-256 and  | Optional           |
|              | MGF1 with SHA-256             |                    |
| PS384        | RSASSA-PSS using SHA-384 and  | Optional           |
|              | MGF1 with SHA-384             |                    |
| PS512        | RSASSA-PSS using SHA-512 and  | Optional           |
|              | MGF1 with SHA-512             |                    |
| none         | No digital signature or MAC   | Optional           |
|              | performed                     |                    |
+--------------+-------------------------------+--------------------+
*/

// Message encapsulation
// Message signing
// Plain fimp msg -> serialize into []binary -> generate signature -> base64 encode original message -> encapsulate into evt.transport.signed or cmd.transport.signed message -> add signature to props
// {
//  "type": "evt.transport.signed",
//  "serv": "sensor_presence",
//  "val_t": "bin",
//  "val": "ewogICJ0eXBlIjogImV2dC5wcmVzZW5jZS5yZXBvcnQiLAogICJzZXJ2IjogInNlbnNvcl9wcmVzZW5jZSIsCiAgInZhbF90IjogImJvb2wiLAogICJ2YWwiOiB0cnVlLAogICJ0YWdzIjogbnVsbCwKICAicHJvcHMiOiBudWxsLAogICJ2ZXIiOiAiMSIsCiAgImNvcmlkIjogIiIsCiAgImN0aW1lIjogIjIwMjAtMDUtMDZUMDk6Mjk6NTkuNTI3KzA1OjAwIiwKICAidWlkIjogIjczZjYxMDMwLTQzOTktNGQyMS1iYjk3LTRjYTdjMTYyM2FjMyIKfQ==",
//  "tags": null,
//  "props": {
//  	"user_id":"aleks@gmail.com",
//      "device_id":"dafdsafsf",
//      "sig":"<computed signature>",
//		"alg":"ES256"
//  },
//  "ver": "1",
//  "corid": "",
//  "ctime": "2020-05-06T09:29:59.527+05:00",
//  "uid": "73f61030-4399-4d21-bb97-4ca7c1623ac3"
//}


// SignMessageES256 encapsulate original message into special transport message with added signature.
func SignMessageES256(payload *fimpgo.FimpMessage,requestMsg *fimpgo.FimpMessage,userId string,keys *security.EcdsaKey,props *fimpgo.Props) (*fimpgo.FimpMessage,error) {
	serializedMsg,err := payload.SerializeToJson()
	if err != nil {
		return nil, err
	}
	if props == nil {
		props = &fimpgo.Props{"user_id":userId}
	}
	msgType := "evt.transport.signed"
	if strings.Contains(payload.Type,"cmd") {
		msgType = "cmd.transport.signed"
	}
	signedMsg := fimpgo.NewBinaryMessage(msgType,payload.Service,serializedMsg,*props,nil,requestMsg)

	signingMethodES256 := &jwt.SigningMethodECDSA{Name: "ES256", Hash: crypto.SHA256, KeySize: 32, CurveBits: 256}
	signature , err := signingMethodES256.Sign(signedMsg.Value.(string),keys.PrivateKey())
	if err != nil {
		return nil,err
	}
	signedMsg.Properties["sig"] = signature
	return signedMsg,nil
}
//
func GetVerifiedMessageES256(signedMsg *fimpgo.FimpMessage,key *security.EcdsaKey) (*fimpgo.FimpMessage,error) {

	if signedMsg.Type != "cmd.transport.signed" && signedMsg.Type != "evt.transport.signed"  {
		return nil,errors.New("incorrect message type")
	}
	origMsgBin , ok1 := signedMsg.Value.(string)
	if !ok1 {
		return nil,errors.New("incorrect encapsulated message format")
	}
	sig , ok2 := signedMsg.Properties["sig"]
	if !ok2 {
		return nil,errors.New("missing signature")
	}
	signingMethodES256 := &jwt.SigningMethodECDSA{Name: "ES256", Hash: crypto.SHA256, KeySize: 32, CurveBits: 256}
	err := signingMethodES256.Verify(origMsgBin,sig,key.PublicKey())
	if err == nil {
		decodedPayloadBin,err := base64.StdEncoding.DecodeString(origMsgBin)
		if err != nil {
			return nil,err
		}
		return fimpgo.NewMessageFromBytes(decodedPayloadBin)
	}else {
		return nil,err
	}
}

