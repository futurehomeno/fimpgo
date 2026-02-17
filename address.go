package fimpgo

import (
	"fmt"
	"strings"

	"github.com/futurehomeno/fimpgo/fimptype"
)

type ResourceTypeT string
type MsgTypeT string

const (
	DefaultPayload        = "j1"
	CompressedJsonPayload = "j1c1"

	MsgTypeCmd     MsgTypeT = "cmd"
	MsgTypeEvt     MsgTypeT = "evt"
	MsgTypeRsp     MsgTypeT = "rsp"
	MsgTypeUnknown MsgTypeT = ""

	ResourceTypeDevice    ResourceTypeT = "dev"
	ResourceTypeApp       ResourceTypeT = "app"
	ResourceTypeAdapter   ResourceTypeT = "ad"
	ResourceTypeCloud     ResourceTypeT = "cloud"
	ResourceTypeDiscovery ResourceTypeT = "discovery"
	ResourceTypeLocation  ResourceTypeT = "loc"
	ResourceTypeUnknown   ResourceTypeT = ""
)

func (rn ResourceTypeT) Str() string {
	return string(rn)
}

func (rn MsgTypeT) Str() string {
	return string(rn)
}

func InterfaceMsgType(iface string) MsgTypeT {
	iface = strings.TrimSpace(iface)

	switch {
	case strings.HasPrefix(iface, MsgTypeCmd.Str()):
		return MsgTypeCmd
	case strings.HasPrefix(iface, MsgTypeEvt.Str()):
		return MsgTypeEvt
	case strings.HasPrefix(iface, MsgTypeCmd.Str()):
		return MsgTypeRsp
	}

	return MsgTypeUnknown
}

type Address struct {
	GlobalPrefix    string
	PayloadType     string
	MsgType         MsgTypeT
	ResourceType    ResourceTypeT
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
	case ResourceTypeAdapter, ResourceTypeApp, ResourceTypeCloud:
		result = fmt.Sprintf("%s/%s/%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType.Str()),
			addr.prepComp("rt", addr.ResourceType.Str()),
			addr.prepComp("rn", addr.ResourceName),
			addr.prepComp("ad", addr.ResourceAddress))
	case ResourceTypeDevice:
		result = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			addr.prepComp("pt", addr.PayloadType),
			addr.prepComp("mt", addr.MsgType.Str()),
			addr.prepComp("rt", addr.ResourceType.Str()),
			addr.prepComp("rn", addr.ResourceName),
			addr.prepComp("ad", addr.ResourceAddress),
			addr.prepComp("sv", addr.ServiceName.Str()),
			addr.prepComp("ad", addr.ServiceAddress))
	case ResourceTypeDiscovery:
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
				addr.MsgType = MsgTypeT(keyVal[1])
			case "rt":
				addr.ResourceType = ResourceTypeT(keyVal[1])
			case "rn":
				addr.ResourceName = keyVal[1]
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
