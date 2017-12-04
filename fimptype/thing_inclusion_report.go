package fimptype


type ThingInclusionReport struct {
	IntegrationId  string    `json:"integr_id" storm:"index"`
	Address        string    `json:"address" storm:"index"`
	Type           string    `json:"type"`
	ProductHash    string    `json:"product_hash"`
	Alias          string    `json:"alias"`
	CommTechnology string    `json:"comm_tech" storm:"index"`
	ProductId      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ManufacturerId string    `json:"manufacturer_id"`
	DeviceId       string    `json:"device_id"`
	HwVersion      string    `json:"hw_ver"`
	SwVersion      string    `json:"sw_ver"`
	PowerSource    string    `json:"power_source"`
	WakeUpInterval string    `json:"wakeup_interval"`
	Security       string    `json:"security"`
	Tags           []string  `json:"tags"`
	PropSets                   map[string]map[string]interface{}  `json:"prop_set"`
	TechSpecificProps          map[string]string             `json:"tech_specific_props"`
	Services       []Service `json:"services"`

}

type Service struct {
	Name          string                    `json:"name" storm:"index"`
	Alias         string                    `json:"alias"`
	Address       string                    `json:"address"`
	Enabled       bool                    `json:"enabled"`
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

