package device

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	message = "00dv=all,2016-01-31,19:27:45,13;"
	timeout = 5 * time.Second
)

type SOP112Response struct {
	Response int                `json:"response"`
	Data     SOP112ResponseData `json:"data"`
}

type SOP112ResponseData struct {
	Serialnumber string `json:"sn"`
}

func Detect(broadcast string) ([]*Powersocket, error) {

	dst, err := net.ResolveUDPAddr("udp", broadcast+":8888")
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	detected := []*Powersocket{}

	// write message as broadcast, expect that powersockets will reply
	conn.WriteTo([]byte(message), dst)

	// process messages
	for {
		b := make([]byte, 4096)

		// wait for answers only for a specific duration
		conn.SetReadDeadline(time.Now().Add(timeout))

		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			// a time is expected and handled as a 'exit' event
			if isTimeout(err) {
				return detected, nil
			}
			// other errors are not expected and returned
			return nil, err
		}

		response := SOP112Response{}
		err = json.Unmarshal(b[:n], &response)
		if err != nil {
			return nil, err
		}

		if (response.Data == SOP112ResponseData{} && response.Data.Serialnumber == "") {
			return nil, fmt.Errorf("couldn't extract name from received data: %s", string(b))
		}

		ps := NewPowersocket(response.Data.Serialnumber, raddr.IP.String())
		detected = append(detected, ps)
	}
}

func isTimeout(err error) bool {
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return true
	}
	return false
}
