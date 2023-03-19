package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(start, deadline *time.Time, stations []*Station) {
	var (
		xmits       []chan Transmission
		dones, acks []chan bool
		backlog     = make(chan Transmission, 10000)
	)

	// FIXME this currently blocks until the main loop is entered
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if deadline != nil {
		log.Info().Time("start", *start).Time("deadline", *deadline).Str("delay", time.Until(*start).String()).Msg("Starting run")
	} else {
		log.Info().Time("start", *start).Str("delay", time.Until(*start).String()).Msg("Starting run")
	}

	// Start up the stations
	for _, station := range stations {
		xmit := make(chan Transmission, 1)
		done := make(chan bool, 1)
		ack := make(chan bool, 1)
		go station.Run(start, xmit, done, ack)
		xmits = append(xmits, xmit)
		dones = append(dones, done)
		acks = append(acks, ack)
	}

	// Propagator resolves which station heard each transmission
	go propagate(backlog)

	// Aim at "start", which is the beginning of minute
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
		for _, xmit := range xmits {
			select {
			case transmission := <-xmit:
				metrics["transmissions"].WithLabelValues(transmission.Band, transmission.Mode, transmission.Station.Callsign).Inc()
				backlog <- transmission // For propagator's consumption
			default:
			}
		}

		select {
		case sig := <-sigs:
			log.Info().Str("signal", sig.String()).Msg("Signal caught, preparing to exit")
			break Loop
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
