package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(start, deadline *time.Time, stations []*Station) {
	// FIXME this currently blocks until the main loop is entered
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if deadline != nil {
		log.Info().Time("start", *start).Time("deadline", *deadline).Str("delay", time.Until(*start).String()).Msg("Starting run")
	} else {
		log.Info().Time("start", *start).Str("delay", time.Until(*start).String()).Msg("Starting run")
	}

	// Aim at "start", which is the beginning of minute
	time.Sleep(time.Until(*start))
	ticker := time.NewTicker(SmallestCommonDuration)

Loop:
	for {
		now := time.Now().UTC()
		if deadline != nil && now.After(*deadline) {
			log.Info().Time("deadline", *deadline).Msg("Deadline reached, ending run")
			break Loop
		}
		log.Info().Msg("Running")

		// Run each station
		var transmissions []*Transmission
		for _, station := range stations {
			if transmission := station.Run(now); transmission != nil {
				metrics["transmissions"].WithLabelValues(transmission.Band, transmission.Mode, transmission.Station.Callsign).Inc()
				transmissions = append(transmissions, transmission)
			}
		}

		// Propagate each transmission
		for _, transmission := range transmissions {
			propagate(transmission)
		}

		// Adjust each station, maybe
		for _, station := range stations {
			station.Adjust()
		}

		select {
		case sig := <-sigs:
			log.Info().Str("signal", sig.String()).Msg("Signal caught, preparing to exit")
			for _, station := range stations {
				station.Spotter.Close()
			}
			break Loop
		case <-ticker.C:
		}
	}
}

// TODO log results for integration testing purposes
