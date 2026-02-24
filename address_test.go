package fimpgo

import (
	"testing"

	"github.com/futurehomeno/fimpgo/fimptype"
)

func TestNewAddressFromStringDevice(t *testing.T) {
	addrString := "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:sensor_presence/ad:16"
	adr, err := NewAddressFromString(addrString)
	if err != nil {
		t.Error(err)
	}
	if adr.PayloadType != "j1" {
		t.Error("Wrong payload type")
	}
	if adr.MsgType != fimptype.MsgTypeEvt {
		t.Error("Wrong message type")
	}
	if adr.ResourceType != fimptype.ResourceTypeDevice {
		t.Error("Wrong resource type")
	}
	if adr.ResourceName != "zw" {
		t.Error("Wrong resource name type")
	}
	if adr.ResourceAddress != "1" {
		t.Error("Wrong resource address")
	}
	if adr.ServiceName != "sensor_presence" {
		t.Error("Wrong payload type")
	}
	if adr.ServiceAddress != "16" {
		t.Error("Wrong service address")
	}
}

func TestNewAddressFromStringAdapter(t *testing.T) {
	addrString := "pt:j1/mt:evt/rt:ad/rn:zw/ad:1"
	adr, err := NewAddressFromString(addrString)
	if err != nil {
		t.Error(err)
	}
	if adr.PayloadType != "j1" {
		t.Error("Wrong payload type")
	}
	if adr.MsgType != fimptype.MsgTypeEvt {
		t.Error("Wrong message type")
	}
	if adr.ResourceType != fimptype.ResourceTypeAdapter {
		t.Error("Wrong resource type")
	}
	if adr.ResourceName != "zw" {
		t.Error("Wrong resource name type")
	}
	if adr.ResourceAddress != "1" {
		t.Error("Wrong resource address")
	}
}

func TestNewAddressFromStringAdapterGlobalPrefix(t *testing.T) {
	addrString := "BDNF123/pt:j1/mt:evt/rt:ad/rn:zw/ad:1"
	adr, err := NewAddressFromString(addrString)
	if err != nil {
		t.Error(err)
	}
	if adr.GlobalPrefix != "BDNF123" {
		t.Error("Wrong global prefix type")
	}
	if adr.PayloadType != "j1" {
		t.Error("Wrong payload type")
	}
	if adr.MsgType != fimptype.MsgTypeEvt {
		t.Error("Wrong message type")
	}
	if adr.ResourceType != fimptype.ResourceTypeAdapter {
		t.Error("Wrong resource type")
	}
	if adr.ResourceName != "zw" {
		t.Error("Wrong resource name type")
	}
	if adr.ResourceAddress != "1" {
		t.Error("Wrong resource address")
	}
}

func TestAddress_Serialize(t *testing.T) {
	adr := Address{MsgType: fimptype.MsgTypeEvt, ResourceType: fimptype.ResourceTypeDevice, ResourceName: fimptype.ZwaveRn, ResourceAddress: "1", ServiceName: "sensor_presence", ServiceAddress: "16"}
	adrStr := adr.Serialize()
	if adrStr != "pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:sensor_presence/ad:16" {
		t.Error("Serialization is incorrect . Result is -  ", adrStr)
	}
}

func TestAddress_SerializeWithGlobalPrefix(t *testing.T) {
	adr := Address{MsgType: fimptype.MsgTypeEvt, ResourceType: fimptype.ResourceTypeDevice, ResourceName: fimptype.ZwaveRn, ResourceAddress: "1", ServiceName: "sensor_presence", ServiceAddress: "16", GlobalPrefix: "DNF123"}
	adrStr := adr.Serialize()
	if adrStr != "DNF123/pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:sensor_presence/ad:16" {
		t.Error("Serialization is incorrect . Result is -  ", adrStr)
	}
}
