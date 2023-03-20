package hamtraffic

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
	"time"
)

// NOTE We're not attempting to model irl propagation of radio waves

const DistanceWiggle = 0.1

func propagate(backlog <-chan Transmission) {
	var (
		bandmodepairToDistance = func() map[string]float64 {
			var (
				basePropagation = NeighbourBins[0]
				distances       = make(map[string]float64)
			)

			for _, bandmodepair := range config.BandModePairs {
				distances[fmt.Sprintf("%s@%s", bandmodepair.Band, bandmodepair.Mode)] = (basePropagation / 2) / math.Sqrt(bandmodepair.CenterFrequency)
			}

			return distances
		}()
		transmission Transmission
	)

	for {
		select {
		case transmission = <-backlog:
			log.Debug().Str("callsign", transmission.Station.Callsign).Str("band", transmission.Band).Str("mode", transmission.Mode).Time("timestamp", transmission.Time).Msg("Propagating")
		default:
			time.Sleep(time.Duration(1 * time.Millisecond))
			continue
		}

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
					station.Receiver <- transmission
				}
			}
		}
	}
}
