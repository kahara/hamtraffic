package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"time"
)

func Run(start, deadline *time.Time, stations []*Station) {
	log.Info().Time("start", *start).Time("deadline", *deadline).Str("delay", time.Until(*start).String()).Msg("Starting run")

	// Aim near "start"
	time.Sleep(time.Until(*start))
	ticker := time.NewTicker(time.Second)
	for {
		now := time.Now().UTC()
		if now.After(*deadline) {
			log.Info().Time("deadline", *deadline).Msg("Deadline reached, ending run")
			return
		}
		log.Info().Msg("Running")

		// Advance world state

		select {
		case <-ticker.C:
		}
	}
}

// TODO log results for integration testing purposes
