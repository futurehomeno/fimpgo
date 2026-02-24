package fimpgo

import (
	"fmt"
	"strings"
)

const (
	DefaultPayload        = "j1"
	CompressedJsonPayload = "j1c1"
	MsgTypeCmd            = "cmd"
	MsgTypeEvt            = "evt"
	MsgTypeRsp            = "rsp"
	ResourceTypeDevice    = "dev"
	ResourceTypeApp       = "app"
	ResourceTypeAdapter   = "ad"
	ResourceTypeCloud     = "cloud"
	ResourceTypeDiscovery = "discovery"
	ResourceTypeLocation  = "loc"
)

type Address struct {
	GlobalPrefix    string
	PayloadType     string
	MsgType         string
	ResourceType    string
	ResourceName    string
	ResourceAddress string
	ServiceName     string
	ServiceAddress  string
}

func (addr *Address) Serialize() string {
	if addr.PayloadType == "" {
		addr.PayloadType = DefaultPayload
	}
	result := ""

	switch addr.ResourceType {
	case ResourceTypeAdapter, ResourceTypeApp, ResourceTypeCloud:
		result = fmt.Sprintf("%s/%s/%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType),
			addr.prepComp("rt", addr.ResourceType),
			addr.prepComp("rn", addr.ResourceName),
			addr.prepComp("ad", addr.ResourceAddress))
	case ResourceTypeDevice:
		result = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType),
			addr.prepComp("rt", addr.ResourceType),
			addr.prepComp("rn", addr.ResourceName),
			addr.prepComp("ad", addr.ResourceAddress),
			addr.prepComp("sv", addr.ServiceName),
			addr.prepComp("ad", addr.ServiceAddress))
	case ResourceTypeDiscovery:
		result = fmt.Sprintf("%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType),
			addr.prepComp("rt", addr.ResourceType))
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
				addr.MsgType = keyVal[1]
			case "rt":
				addr.ResourceType = keyVal[1]
			case "rn":
				addr.ResourceName = keyVal[1]
			case "ad":
				if addr.ServiceName == "" {
					addr.ResourceAddress = keyVal[1]
				} else {
					addr.ServiceAddress = keyVal[1]
				}

			case "sv":
				addr.ServiceName = keyVal[1]
			}
		default:
			return nil, fmt.Errorf("invalid address format key=%v", keyVal)
		}
	}

	return &addr, nil
}
