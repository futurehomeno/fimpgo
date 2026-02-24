package fimptype

import "strings"

type MsgTypeT string

const (
	MsgTypeCmd     MsgTypeT = "cmd"
	MsgTypeEvt     MsgTypeT = "evt"
	MsgTypeRsp     MsgTypeT = "rsp"
	MsgTypeUnknown MsgTypeT = ""
)

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
	case strings.HasPrefix(iface, MsgTypeRsp.Str()):
		return MsgTypeRsp
	}

	return MsgTypeUnknown
}
