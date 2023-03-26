package hamtraffic

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
)

// NOTE We're not attempting to model irl propagation of radio waves

const DistanceWiggle = 0.5

var bandmodepairToDistance map[string]float64

func propagate(transmission *Transmission) {
	if bandmodepairToDistance == nil {
		bandmodepairToDistance = func() map[string]float64 {
			var (
				basePropagation = NeighbourBins[1]
				distances       = make(map[string]float64)
			)

			for _, bandmodepair := range config.BandModePairs {
				distances[fmt.Sprintf("%s@%s", bandmodepair.Band, bandmodepair.Mode)] = (basePropagation / 2) / math.Sqrt(bandmodepair.CenterFrequency/1000000)
			}

			return distances
		}()
	}

	log.Debug().Str("callsign", transmission.Station.Callsign).Str("band", transmission.Band).Str("mode", transmission.Mode).Time("timestamp", transmission.Time).Msg("Propagating")

	// Compute a base value of how far this transmission will travel, based solely on the hyper-simplification that
	// transmissions on lower frequencies can be heard farther away than transmissions on higher frequencies
	baseDistance := bandmodepairToDistance[fmt.Sprintf("%s@%s", transmission.Band, transmission.Mode)]

	// Randomly wiggle the base value
	baseWiggle := rand.Float64() * DistanceWiggle
	baseDistance = baseDistance + (-baseWiggle + (baseWiggle * 2.0))

	// Send transmission to each station that is close enough to hear it
	for d, bin := range transmission.Station.Neighbours {
		for _, station := range bin {
			// Wiggle further for each receiving station
			wiggle := rand.Float64() * DistanceWiggle
			distance := baseDistance + (-wiggle + (wiggle * 2.0))

			// Send; receiver decides if it was receiving at the time of transmission
			if distance >= d {
				station.Receive(transmission)
			}
		}
	}
}
