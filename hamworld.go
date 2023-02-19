package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"time"
)

func Run(start, deadline *time.Time, stations []*Station) {
	var dones, acks []chan bool

	if deadline != nil {
		log.Info().Time("start", *start).Time("deadline", *deadline).Str("delay", time.Until(*start).String()).Msg("Starting run")
	} else {
		log.Info().Time("start", *start).Str("delay", time.Until(*start).String()).Msg("Starting run")
	}

	// Start up the stations
	for _, station := range stations {
		done := make(chan bool, 1)
		ack := make(chan bool, 1)
		go station.Run(done, ack)
		dones = append(dones, done)
		acks = append(acks, ack)
	}

	// Aim near "start"
	time.Sleep(time.Until(*start))
	ticker := time.NewTicker(time.Second)
Loop:
	for {
		now := time.Now().UTC()
		if deadline != nil && now.After(*deadline) {
			log.Info().Time("deadline", *deadline).Msg("Deadline reached, ending run")
			break Loop
		}
		log.Info().Msg("Running")

		// Run the world

		select {
		case <-ticker.C:
		}
	}

	// Tell stations to shut down
	for _, done := range dones {
		done <- true
	}

	// Wait for station shutdown confirmations
	for _, ack := range acks {
		<-ack
	}
}

// TODO log results for integration testing purposes
