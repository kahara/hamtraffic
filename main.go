package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"time"
)

var (
	config *Config
	world  *World
)

func Init() {
	var (
		stations        []*Station
		start, deadline time.Time
	)

	go Metrics()

	config = NewConfig()

	// https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap/
	ballast := make([]byte, config.StationCount*1048576)
	log.Info().Int("size", len(ballast)).Msg("Ballast reserved")

	world = NewWorld()

	for i := 0; i < config.StationCount; i++ {
		stations = append(stations, NewStation(config, world, metrics["report_packets"]))
	}

	log.Debug().Msg("Computing neighbourhoods")
	for _, station := range stations {
		station.ComputeNeighbourhood(stations)
	}

	// Start running at the start of next minute
	start = time.Now().UTC()
	start = start.Truncate(time.Duration(time.Minute))
	start = start.Add(time.Minute)

	if config.Freerun {
		Run(&start, nil, stations)
	} else {
		deadline = start.Add(*config.Runtime)
		Run(&start, &deadline, stations)
	}
}
