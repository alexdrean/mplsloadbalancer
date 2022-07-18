package main

import (
	"fmt"
	"log"
	"mplsloadbalancer/config"
	"mplsloadbalancer/snmp"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg.Paths)
	repeat := 1

	dumpRadio("AFLTU Camelback", "100.64.10.200", snmp.AFLTU{Direction: snmp.DirectionRXMatters}, repeat, false)
	dumpRadio("AFLTU Oatman - Ajo", "100.64.10.180", snmp.AFLTU{Direction: snmp.DirectionRXMatters}, repeat, false)
	dumpRadio("AF11FX Oatman - Gila Bend", "100.64.1.4", snmp.AF{Direction: snmp.DirectionRXMatters}, repeat, false)
}


func dumpRadio(name string, ip string, radioType snmp.RadioType, repeat int, debug bool) {
	err, radio := snmp.CreateRadio(name, ip, &radioType, debug)
	defer radio.Close()
	if err != nil {
		log.Panicln(err)
	}
	for i := 0; i < repeat; i++ {
		status := radio.GetStatus()
		if status != nil {
			println(status.CapacityMatters, radio.Name)
		}
	}
}
