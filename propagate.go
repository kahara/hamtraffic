package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"time"
)

// NOTE We're not attempting to model irl propagation of radio waves

type Timeline struct {
}

func propagate(backlog <-chan Transmission) {
	for {
		select {
		case transmission := <-backlog:
			log.Debug().Str("callsign", transmission.Station.Callsign).Float64("frequency", transmission.Frequency).Str("mode", transmission.Mode).Time("timestamp", transmission.Time).Msg("Propagator received a transmission")
		}

		time.Sleep(time.Duration(1 * time.Millisecond))
	}
}
