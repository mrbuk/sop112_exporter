package main

// Powersocket describes a SOP112 power socket by Name (serial number) and IP
type Powersocket struct {
	Name string
	IP   string
}

// PowerConsumption describes the data strucutre returned by each SOP112 power
// socket for the 511 measurement command
type PowerConsumption struct {
	Response int                  `json:response`
	Code     int                  `json:code`
	Data     PowerConsumptionData `json:data`
}

// PowerConsumptionData describes the actual consumption data (wattage, amps)
// as part of a measurement command. Additionally it holds the switch state (on/off)
// All fields are
// 		1) kept as strings, as in the original data structure
//			(e.g. SwitchState 0 = off, 1 = on)
// 		2) SwitchState =
type PowerConsumptionData struct {
	Watt        []string `json:watt`
	Amp         []string `json:amp`
	SwitchState []int    `json:switch`
}
