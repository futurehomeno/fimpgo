package fimptype

type Service struct {
	Name          string                    `json:"name" storm:"index"`
	Alias         string                    `json:"alias"`
	Address       string                    `json:"address"`
	Enabled       bool                      `json:"enabled"`
	Groups        []string                  `json:"groups"`
	Props         map[string]interface{}    `json:"props"`
	Tags          []string                  `json:"tags"`
	PropSetReference string 				`json:"prop_set_ref"`
	Interfaces    []Interface               `json:"interfaces"`
}

type Interface struct {
	Type      string `json:"intf_t"`
	MsgType   string `json:"msg_t"`
	ValueType string `json:"val_t"`
	Version   string `json:"ver"`
}


