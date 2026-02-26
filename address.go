package fimpgo

import (
	"fmt"
	"strings"

	"github.com/futurehomeno/fimpgo/fimptype"
)

const (
	DefaultPayload        = "j1"
	CompressedJsonPayload = "j1c1"
)

type Address struct {
	GlobalPrefix    string
	PayloadType     string
	MsgType         fimptype.MsgTypeT
	ResourceType    fimptype.ResourceTypeT
	ResourceName    fimptype.ResourceNameT
	ResourceAddress string
	ServiceName     fimptype.ServiceNameT
	ServiceAddress  string
}

func (addr *Address) Serialize() string {
	if addr.PayloadType == "" {
		addr.PayloadType = DefaultPayload
	}
	result := ""

	switch addr.ResourceType {
	case fimptype.ResourceTypeAdapter, fimptype.ResourceTypeApp, fimptype.ResourceTypeCloud:
		result = fmt.Sprintf("%s/%s/%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType.Str()),
			addr.prepComp("rt", addr.ResourceType.Str()),
			addr.prepComp("rn", addr.ResourceName.Str()),
			addr.prepComp("ad", addr.ResourceAddress))
	case fimptype.ResourceTypeDevice:
		result = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType.Str()),
			addr.prepComp("rt", addr.ResourceType.Str()),
			addr.prepComp("rn", addr.ResourceName.Str()),
			addr.prepComp("ad", addr.ResourceAddress),
			addr.prepComp("sv", addr.ServiceName.Str()),
			addr.prepComp("ad", addr.ServiceAddress))
	case fimptype.ResourceTypeDiscovery:
		result = fmt.Sprintf("%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType.Str()),
			addr.prepComp("rt", addr.ResourceType.Str()))
	}
	if addr.GlobalPrefix != "" {
		result = addr.GlobalPrefix + "/" + result
	}
	return result
}

func (addr *Address) prepComp(prefix string, comp string) string {
	if comp == "+" || comp == "#" {
		return comp
	} else {
		return fmt.Sprintf("%s:%s", prefix, comp)
	}
}

func NewAddressFromString(address string) (*Address, error) {
	addr := Address{}
	tokens := strings.Split(address, "/")

	for index, tok := range tokens {
		keyVal := strings.Split(tok, ":")
		// detecting global prefix
		switch {
		case len(keyVal) == 1 && index == 0:
			addr.GlobalPrefix = keyVal[0]
		case len(keyVal) == 2:
			switch keyVal[0] {
			case "pt":
				addr.PayloadType = keyVal[1]
			case "mt":
				addr.MsgType = fimptype.MsgTypeT(keyVal[1])
			case "rt":
				addr.ResourceType = fimptype.ResourceTypeT(keyVal[1])
			case "rn":
				addr.ResourceName = fimptype.ResourceNameT(keyVal[1])
			case "ad":
				if addr.ServiceName == "" {
					addr.ResourceAddress = keyVal[1]
				} else {
					addr.ServiceAddress = keyVal[1]
				}

			case "sv":
				addr.ServiceName = fimptype.ServiceNameT(keyVal[1])
			}
		default:
			return nil, fmt.Errorf("invalid address format key=%v", keyVal)
		}
	}

	return &addr, nil
}
