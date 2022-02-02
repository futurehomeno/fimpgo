package fimptype

// ThingInclusionReport is the object send as value in thing inclusion reports.
type ThingInclusionReport struct {
	IntegrationId     string                            `json:"integr_id" storm:"index"`
	Address           string                            `json:"address" storm:"index"`
	Type              string                            `json:"type"`
	ProductHash       string                            `json:"product_hash"`
	Alias             string                            `json:"alias"`
	CommTechnology    string                            `json:"comm_tech" storm:"index"`
	ProductId         string                            `json:"product_id"`
	ProductName       string                            `json:"product_name"`
	ManufacturerId    string                            `json:"manufacturer_id"`
	DeviceId          string                            `json:"device_id"`
	HwVersion         string                            `json:"hw_ver"`
	SwVersion         string                            `json:"sw_ver"`
	PowerSource       string                            `json:"power_source"`
	WakeUpInterval    string                            `json:"wakeup_interval"`
	Security          string                            `json:"security"`
	Tags              []string                          `json:"tags"`
	Groups            []string                          `json:"groups"`
	PropSets          map[string]map[string]interface{} `json:"prop_set"`
	TechSpecificProps map[string]string                 `json:"tech_specific_props"`
	Services          []Service                         `json:"services"`
}

// ThingExclusionReport is the object send as value in thing exclusion reports.
type ThingExclusionReport struct {
	Address string `json:"address"`
}