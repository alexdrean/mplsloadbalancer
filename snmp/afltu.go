package snmp

import "github.com/gosnmp/gosnmp"

type AFLTU struct {
	Direction
}

func (device AFLTU) GetCapacity(snmp *gosnmp.GoSNMP) (error, *RadioStatus) {
	const (
		ROLE_AP  = 0
		ROLE_STA = 1
	)
	result, err := snmp.GetNext([]string{
		".1.3.6.1.4.1.41112.1.10.1.2.1",   // Role. AP = 0. STA = 1
		".1.3.6.1.4.1.41112.1.10.1.4.1.3", // STA Tx
		".1.3.6.1.4.1.41112.1.10.1.4.1.4", // STA Rx
	})
	if err != nil {
		return err, nil
	}
	deviceRole := result.Variables[0].Value.(int)
	staTX := int64(result.Variables[1].Value.(int)) * 1000
	staRX := int64(result.Variables[2].Value.(int)) * 1000

	status := RadioStatus{}
	if deviceRole == ROLE_STA {
		status.CapacityRx = staRX
		status.CapacityTx = staTX
	} else if deviceRole == ROLE_AP {
		status.CapacityRx = staTX
		status.CapacityTx = staRX
	}

	if device.Direction == DirectionRXMatters {
		status.CapacityMatters = status.CapacityRx
	} else if device.Direction == DirectionTXMatters {
		status.CapacityMatters = status.CapacityTx
	}
	return nil, &status
}
