package state

import (
	"errors"
	"fmt"
	"mplsloadbalancer/config"
	"mplsloadbalancer/snmp"
)

type link struct {
	label uint32
	radio *snmp.Radio
}

var InvalidPolarityError = errors.New("invalid radio polarity")
var InvalidRadioType = errors.New("invalid radio type")


func (compiler Compiler) loadConfig(cfg config.Config) error {
	compiler.paths = make([][]link, len(cfg.Paths))
	for i, cpath := range cfg.Paths {
		compiler.paths[i] = make([]link, len(cpath.Links))
		for j, clink := range cpath.Links {
			cradio := clink.Radio

			radioName := fmt.Sprintf("%s %s-%s (%s)", cradio.Type, cpath.Close, cpath.Far, cradio.Polarity)

			var direction snmp.Direction
			if cradio.Polarity == config.PolarityClose {
				direction = snmp.DirectionTXMatters
			} else if cradio.Polarity == config.PolarityFar {
				direction = snmp.DirectionRXMatters
			} else {
				return InvalidPolarityError
			}

			var radioType *snmp.RadioType = nil
			if cradio.Type == config.TypeAF {
				af := snmp.RadioType(snmp.AF{Direction: direction})
				radioType = &af
			} else if cradio.Type == config.TypeAFLTU {
				afltu := snmp.RadioType(snmp.AFLTU{Direction: direction})
				radioType = &afltu
			} else {
				return InvalidRadioType
			}

			err, radio := snmp.CreateRadio(radioName, cradio.Ip, radioType, false)
			if err != nil {
				return err
			}

			compiler.paths[i][j] = link{
				label: clink.Label,
				radio: radio,
			}
		}
	}
	return nil
}
