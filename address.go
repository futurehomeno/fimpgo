package fimpgo

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	MsgTypeCmd          = "cmd"
	MsgTypeEvt          = "evt"
	ResourceTypeDevice  = "dev"
	ResourceTypeApp     = "app"
	ResourceTypeAdapter = "ad"
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

func (adr *Address) Serialize() string {
	if adr.PayloadType == "" {
		adr.PayloadType = "j1"
	}
	result := ""

	switch adr.ResourceType {

	case ResourceTypeAdapter, ResourceTypeApp:
		result = fmt.Sprintf("%s/%s/%s/%s/%s",
			adr.prepComp("pt", adr.PayloadType),
			adr.prepComp("mt", adr.MsgType),
			adr.prepComp("rt", adr.ResourceType),
			adr.prepComp("rn", adr.ResourceName),
			adr.prepComp("pt", adr.ResourceAddress))
	case ResourceTypeDevice:
		result = fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			adr.prepComp("pt", adr.PayloadType),
			adr.prepComp("mt", adr.MsgType),
			adr.prepComp("rt", adr.ResourceType),
			adr.prepComp("rn", adr.ResourceName),
			adr.prepComp("ad", adr.ResourceAddress),
			adr.prepComp("sv", adr.ServiceName),
			adr.prepComp("ad", adr.ServiceAddress))

	}
	if adr.GlobalPrefix != "" {
		result = adr.GlobalPrefix + "/" + result
	}
	return result
}

func (adr *Address) prepComp(prefix string, comp string) string {
	if comp == "+" || comp == "#" {
		return fmt.Sprintf("%s", comp)
	} else {
		return fmt.Sprintf("%s:%s", prefix, comp)
	}
}

func NewAddressFromString(address string) (*Address, error) {
	adr := Address{}
	tokens := strings.Split(address, "/")
	var err error
	for index, _ := range tokens {
		keyVal := strings.Split(tokens[index], ":")
		// detecting global prefix
		if len(keyVal) == 1 && index == 0 {
			adr.GlobalPrefix = keyVal[0]
		} else if len(keyVal) == 2 {
			switch keyVal[0] {
			case "pt":
				adr.PayloadType = keyVal[1]
			case "mt":
				adr.MsgType = keyVal[1]
			case "rt":
				adr.ResourceType = keyVal[1]
			case "rn":
				adr.ResourceName = keyVal[1]
			case "ad":
				if adr.ServiceName == "" {
					adr.ResourceAddress = keyVal[1]
				} else {
					adr.ServiceAddress = keyVal[1]
				}

			case "sv":
				adr.ServiceName = keyVal[1]

			}
		} else {
			return nil, errors.New("Incorrectly formatted address")
		}
	}

	return &adr, err

}
