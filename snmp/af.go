package snmp

import "github.com/gosnmp/gosnmp"

type AF struct {
	Direction
}

func (device AF) GetCapacity(snmp *gosnmp.GoSNMP) (error, *RadioStatus) {
	result, err := snmp.GetNext([]string{
		".1.3.6.1.4.1.41112.1.3.2.1.5", // RX
		".1.3.6.1.4.1.41112.1.3.2.1.6", // TX
	})
	if err != nil {
		return err, nil
	}
	status := RadioStatus{}
	status.CapacityRx = int64(result.Variables[0].Value.(int))
	status.CapacityTx = int64(result.Variables[1].Value.(int))
	if device.Direction == DirectionRXMatters {
		status.CapacityMatters = status.CapacityRx
	} else if device.Direction == DirectionTXMatters {
		status.CapacityMatters = status.CapacityTx
	}
	return nil, &status
}
