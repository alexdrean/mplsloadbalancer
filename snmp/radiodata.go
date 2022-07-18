package snmp

import "github.com/gosnmp/gosnmp"

type Direction int

const (
	DirectionTXMatters = 1
	DirectionRXMatters = 0
)

type RadioType interface {
	GetCapacity(snmp *gosnmp.GoSNMP) (error, *RadioStatus)
}
