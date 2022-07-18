package buckets

import (
	"math"
	"mplsloadbalancer/state"
)

type Bucket []uint32

func GetBuckets(size int, state state.State) []Bucket {
	buckets := make([]Bucket, size)
	for n := 0; n < size; n++ {
		buckets[n] = make(Bucket, len(state.Paths))
	}
	for j, path := range state.Paths {
		reductionRatio := float64(size) / float64(path.Capacity)
		roundedLinks := make([]MixerItem, len(path.Links))
		for i, link := range path.Links {
			roundedLinks[i] = MixerItem{
				Value: link.Label,
				Count: int(roundAlternate(float64(link.Capacity) * reductionRatio, i)),
			}
		}
		for i, label := range mix(roundedLinks) {
			buckets[i][j] = label
		}
	}
	return buckets
}

type MixerItem struct {
	Value   uint32
	Count   int
	current float64
}

func mix(items []MixerItem) []uint32 {
	n := 0
	for _, item := range items {
		n += item.Count
	}
	result := make([]uint32, n)
	for i := 0; i < n; {
		lowest := MixerItem{current: math.MaxInt}
		for _, item := range items {
			if lowest.current > item.current {
				lowest = item
			}
		}
		result[i] = lowest.Value
		lowest.current += float64(n) / float64(lowest.Count)
	}
	return result
}

func roundAlternate(n float64, i int) float64 {
	if i % 2 == 0 {
		return math.Ceil(n)
	} else {
		return math.Floor(n)
	}
}

// https://www.spinics.net/lists/lartc/msg23460.html don't do this, do this instead:
// tc filter add dev enp40s0f0 protocol ip handle 6 fw action mpls push label 101 ttl 64 action mpls push label 102 ttl 64
// https://wiki.nftables.org/wiki-nftables/index.php/Math_operations
