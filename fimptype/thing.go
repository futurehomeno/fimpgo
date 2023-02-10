package fimptype

// ThingInclusionReport is the object send as value in thing inclusion reports.
type ThingInclusionReport struct {
	Address  string    `json:"address"`  // Address of the thing, which is unique identifier within the adapter.
	Groups   []string  `json:"groups"`   // Groups are used to link multiple services in one logical group. Each group is effectively a separate device within a single thing, usually equal to channels in ZWave or endpoints in Zigbee.
	Services []Service `json:"services"` // Definitions of all services provided by the thing.

	ProductName    string `json:"product_name"`    // Optional initial human-readable name of the device as shown to the user.
	ProductHash    string `json:"product_hash"`    // Product hash is a unique identifier of the product consisting of adapter, manufacturer and product identifiers
	ProductId      string `json:"product_id"`      // Name of the model of the device.
	ManufacturerId string `json:"manufacturer_id"` // Name of the manufacturer of the device.
	DeviceId       string `json:"device_id"`       // Optional unique identifier or serial number of the device.

	HwVersion      string `json:"hw_ver"`          // Optional hardware version of the device.
	SwVersion      string `json:"sw_ver"`          // Optional software version of the device.
	CommTechnology string `json:"comm_tech"`       // Communication technology used by the adapter to communicate with the device. Possible values: "zw", "zigbee". TODO: work out conventions for cloud and local adapters.
	PowerSource    string `json:"power_source"`    // Power source of the device. Possible values: "dc", "ac", "battery". TODO: check if it should be "battery" or "bat".
	WakeUpInterval string `json:"wakeup_interval"` // Wakeup interval for battery powered devices. TODO: check if it should be integer or string and where it is used.
	Security       string `json:"security"`        // Level of communication security. Possible values: "insecure", "secure". // TODO: reevaluate usage and implementation.

	Alias         string   `json:"alias"`     // TODO: check if it is used and if it is needed.
	Tags          []string `json:"tags"`      // TODO: check if it is used and if it is needed.
	Type          string   `json:"type"`      // TODO: check if it is used and if it is needed.
	IntegrationId string   `json:"integr_id"` // TODO: check if it is used and if it is needed.

	TechSpecificProps map[string]string                 `json:"tech_specific_props"` // Custom properties of the thing specific to the technology adapter..
	PropSets          map[string]map[string]interface{} `json:"prop_set"`            // Custom property sets shared by multiple services.
}

// ThingExclusionReport is the object send as value in thing exclusion reports.
type ThingExclusionReport struct {
	Address string `json:"address"` // Address of the thing, which is unique identifier within the adapter.
}
