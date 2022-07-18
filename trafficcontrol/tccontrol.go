package trafficcontrol

import (
	"fmt"
	"github.com/alexdrean/go-tc"
	"math"
	"mplsloadbalancer/buckets"
	"net"
	"os"
)

type TC struct {
	handle *tc.Tc
	devID  *net.Interface
	kill   chan bool
	ttl    uint8
}

func Open(iface string, offset int, ttl uint8) (*TC, error) {
	result := TC{
		kill: make(chan bool),
	}
	result.waitForClose()
	devID, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	result.devID = devID
	handle, err := tc.Open(&tc.Config{})
	if err != nil {
		return nil, err
	}
	result.handle = handle
	return &result, nil
}

func (o TC) Close() {
	o.kill <- true
}

func (o TC) MapBuckets(buckets []buckets.Bucket, offset int) {
	actualOffset, mask := computeOffsetAndMask(offset, len(buckets))
	for i, labels := range buckets {
		mplsActions := make([]*tc.Action, len(labels))
		for j, label := range labels {
			mplsActions[j] = &tc.Action{
				Kind: "mpls",
				MPLS: &tc.MPLS{
					Parms: &tc.MPLSParam{
						Index:   uint32(j),
						MAction: tc.MPLSActPush,
					},
					Label: &label,
					TTL:   &o.ttl,
				},
			}
		}
		filter := tc.Object{
			Msg: tc.Msg{
				Ifindex: uint32(o.devID.Index),
				Handle:  actualOffset + uint32(i),
			},
			Attribute: tc.Attribute{
				Kind: "fw",
				Fw: &tc.Fw{
					Mask:    &mask,
					Actions: &mplsActions,
				},
			},
		}
		o.handle.Filter().Replace(&filter)
	}
}

func computeOffsetAndMask(offset int, count int) (uint32, uint32) {
	actualOffset := 1 << (offset - 1)
	mask := math.Pow(2, math.Floor(math.Log(float64(count))+1)) + actualOffset - 1
	return actualOffset, mask
}

func (o TC) waitForClose() {
	defer func() {
		if o.handle != nil {
			if err := o.handle.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
			}
		}
	}()
	for !<-o.kill {
	}
}
