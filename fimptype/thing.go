package fimptype

// ThingInclusionReport is the object send as value in thing inclusion reports.
type ThingInclusionReport struct {
	Address  string    `json:"address"`  // An arbitrary unique identifier of the thing within the adapter. Must consist only from alphanumeric characters. For example in Z-Wave it is equal to Node ID, while in Zigbee to UDID.
	Groups   []string  `json:"groups"`   // Groups are used to link multiple services into one logical group. Each group is effectively a separate device within a single thing, equal to channels in Z-Wave or endpoints in Zigbee.
	Services []Service `json:"services"` // An array of service definition objects for all services provided by the thing.

	ProductName    string `json:"product_name"`    // Optional initial human-readable name of the device as shown to the user. If empty falls back to product hash.
	ProductHash    string `json:"product_hash"`    // Product hash is a unique identifier of the product consisting of joined adapter, manufacturer and product identifiers.
	ProductId      string `json:"product_id"`      // Name of identification of the model of the device.
	ManufacturerId string `json:"manufacturer_id"` // Name or identification of the manufacturer of the device.
	DeviceId       string `json:"device_id"`       // Optional unique identifier or serial number of the device.

	HwVersion      string `json:"hw_ver"`          // Optional hardware version of the device.
	SwVersion      string `json:"sw_ver"`          // Optional software version of the device.
	CommTechnology string `json:"comm_tech"`       // Technology used by the adapter to communicate with the device, one of "zw", "zigbee", "local_network", "cloud values".
	PowerSource    string `json:"power_source"`    // Power source of the device. Possible values: "dc", "ac", "battery".
	WakeUpInterval string `json:"wakeup_interval"` // Wakeup interval for battery powered devices in seconds, value "-1" indicates that it is not applicable.
	Security       string `json:"security"`        // Level of communication security, either insecure or secure.

	TechSpecificProps map[string]string         `json:"tech_specific_props"` // Optional custom properties of the thing specific to the technology adapter.
	PropSets          map[string]map[string]any `json:"prop_set"`            // Optional map of custom property sets of services specific to the technology adapter. These sets can be referenced from service definition.
}

// ThingExclusionReport is the object send as value in thing exclusion reports.
type ThingExclusionReport struct {
	Address string `json:"address"` // Address of the thing, which is unique identifier within the adapter.
}
