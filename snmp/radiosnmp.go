package snmp

import (
	"github.com/gosnmp/gosnmp"
	"log"
	"os"
	"time"
)

type Radio struct {
	Status     *RadioStatus
	Ip         string
	StatusChan chan *RadioStatus
	killChan   chan bool
	Type       *RadioType
	Name       string
}

type RadioStatus struct {
	CapacityTx      int64
	CapacityRx      int64
	CapacityMatters int64
}

func (radio Radio) Close() {
	radio.killChan <- true
}

func (radio Radio) GetStatus() *RadioStatus {
	radio.StatusChan <- nil
	return <-radio.StatusChan
}

func CreateRadio(name string, ip string, radioType *RadioType, debug bool) (error, *Radio) {
	err, snmp := connectSNMP(ip, debug)
	if err != nil {
		return err, nil
	}
	radio := &Radio{
		Status:     nil,
		Ip:         ip,
		StatusChan: make(chan *RadioStatus),
		killChan:   make(chan bool),
		Type:       radioType,
		Name:       name,
	}
	go radioLoop(snmp, radio)
	return nil, radio
}

func radioLoop(snmp *gosnmp.GoSNMP, radio *Radio) {
	defer snmp.Conn.Close()
	for {
		select {
		case <-radio.killChan:
			return
		case <-radio.StatusChan:
			if err, status := (*radio.Type).GetCapacity(snmp); err == nil {
				radio.Status = status
				radio.StatusChan <- status
			} else {
				log.Println(err)
				radio.StatusChan <- nil
			}
		}
	}
}

func connectSNMP(ip string, debug bool) (error, *gosnmp.GoSNMP) {
	var logger gosnmp.Logger
	if debug {
		logger = gosnmp.NewLogger(log.New(os.Stdout, "<SNMP> ", log.Flags()))
	}
	snmp := gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Community: "public",
		Version:   gosnmp.Version1,
		Timeout:   5 * time.Second,
		Logger:    logger,
	}
	if err := snmp.Connect(); err != nil {
		return err, nil
	}
	return nil, &snmp
}

