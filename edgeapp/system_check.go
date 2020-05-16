package edgeapp

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)


type SystemCheck struct {

}

func NewSystemCheck() *SystemCheck {
	return &SystemCheck{}
}

func (sc *SystemCheck) IsNetworkAvailable() bool {
	netIntfs,err := net.Interfaces()
	if err != nil {
		log.Error("<sys-check> Interface check error:",err.Error())
		return false
	}

	for i := range netIntfs {
		log.Tracef("Name %s , flags %s ",netIntfs[i].Name,netIntfs[i].Flags.String())
		if strings.Contains(netIntfs[i].Name,"tap0")|| strings.Contains(netIntfs[i].Name,"bridge") ||
			strings.Contains(netIntfs[i].Flags.String(),"loopback")  {
			// skipping zipgateway and local interfaces
			continue
		}else {
			if strings.Contains(netIntfs[i].Flags.String(),"up") && strings.Contains(netIntfs[i].Flags.String(),"broadcast") {
				addrs,err := netIntfs[i].Addrs()
				if err != nil {
					log.Trace("Address returned error :",err.Error())
					continue
				}
				for i2 := range addrs {
					log.Trace("Checking address :",addrs[i2].String())
					if len(addrs[i2].String())>=4 {
						return true
					}
				}
			}
		}
	}
	return false
}

func (sc *SystemCheck) IsInternetAvailable()bool {
	if !sc.IsNetworkAvailable() {
		return false
	}
	ips, err := net.LookupIP("google.com")
	if err != nil {
		return false
	}
	log.Trace("google.com resolved to ",ips)
	return true
}

func (sc *SystemCheck) WaitForNetwork(timeout time.Duration) error {
	var elapsedTime time.Duration
	for {
		startTime := time.Now()
		if sc.IsNetworkAvailable() {
			return nil
		}
		time.Sleep(time.Second*5)
		if timeout == 0 {
			continue
		}
		elapsedTime+= time.Since(startTime)
		if elapsedTime>timeout {
			return errors.New("timeout")
		}
	}
}

func (sc *SystemCheck) WaitForInternet(timeout time.Duration) error {
	var elapsedTime time.Duration
	for {
		startTime := time.Now()
		if sc.IsInternetAvailable() {
			return nil
		}
		time.Sleep(time.Second*5)
		if timeout == 0 {
			continue
		}
		elapsedTime+= time.Since(startTime)
		if elapsedTime>timeout {
			return errors.New("timeout")
		}
	}
}