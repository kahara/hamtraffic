package hamtraffic

import (
	"github.com/rs/zerolog/log"
)

// NOTE We're not attempting to model irl propagation of radio waves

const ()

func propagate(backlog <-chan Transmission) {
	var transmission Transmission

	for {
		select {
		case transmission = <-backlog:
			log.Debug().Str("callsign", transmission.Station.Callsign).Str("band", transmission.Band).Str("mode", transmission.Mode).Time("timestamp", transmission.Time).Msg("Propagating")
		default:
			continue
		}

		// Compute a base value of how far this transmission will travel, based solely on the hyper-simplification that
		// transmissions on lower frequencies can be heard farther away than transmissions on higher frequencies

		// Randomly wiggle the base value

		// Send transmission to each station that is close enough to hear it

		// Wiggle further
	}
}
