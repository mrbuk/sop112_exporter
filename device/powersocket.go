package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// defaultTimeout for the http.Client used to connect
// to the Powersocket
var defaultTimeout = 10

// Powersocket describes a SOP112 power socket by Name (serial number) and IP
type Powersocket struct {
	Name       string
	IP         string
	httpClient *http.Client
}

// PowerConsumption describes the data strucutre returned by each SOP112 power
// socket for the 511 measurement command
type PowerConsumption struct {
	Response int                  `json:"response"`
	Code     int                  `json:"code"`
	Data     PowerConsumptionData `json:"data"`
}

// PowerConsumptionData describes the actual consumption data (wattage, amps)
// as part of a measurement command. Additionally it holds the switch state (on/off)
// All fields are
// 		1) kept as strings, as in the original data structure
//			(e.g. SwitchState 0 = off, 1 = on)
// 		2) SwitchState =
type PowerConsumptionData struct {
	Watt        []string `json:"watt"`
	Amp         []string `json:"amp"`
	SwitchState []int    `json:"switch"`
}

// Measureable describes how a measuring can be executed
//go:generate counterfeiter . Measureable
type Measureable interface {
	Get() (float64, error)
}

// NewPowersocket creates a new instance of a Powersocket with
// provided values and a default http.Client (using defaultTimeout)
func NewPowersocket(name string, ip string) *Powersocket {
	return &Powersocket{
		Name: name,
		IP:   ip,
		httpClient: &http.Client{
			Timeout: time.Duration(defaultTimeout) * time.Second,
		},
	}
}

// SetTimeout changes the timeout for the http.Client associated
// with a Powersocket
func (ps *Powersocket) SetTimeout(t time.Duration) {
	ps.httpClient.Timeout = t
}

// Get fetches the wattage from the Powersocket
func (ps *Powersocket) Get() (float64, error) {

	var consumption PowerConsumption

	log.Debugf("Checking '%s' for consumption data", ps.IP)
	response, err := ps.httpClient.Get(fmt.Sprintf("http://%s/?cmd=511", ps.IP))
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != err {
		return 0.0, err
	}

	err = json.Unmarshal(body, &consumption)
	if err != nil {
		return 0.0, err
	}

	rawWattage := consumption.Data.Watt
	if len(rawWattage) < 1 {
		return 0.0, errors.New("No measures")
	}

	wattage, err := strconv.ParseFloat(rawWattage[0], 64)
	if err != nil {
		return 0.0, err
	}

	return wattage, nil
}

func (ps *Powersocket) String() {
	fmt.Sprintf("%s[%s]", ps.Name, ps.IP)
}
