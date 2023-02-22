package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

const (
	DefaultFreerun      = true
	DefaultRuntime      = "900s"
	DefaultStationCount = 10000
	MaxStationCount     = 65536
	DefaultBands        = "160m:0.25,80m:0.40,40m:0.65,20m:1.0,10m:0.65,6m:0.25,2m:0.15"
	DefaultModes        = "FT8:1.0,FT4:0.25,CW:0.15"
	DefaultStickiness   = 0.9999

	PrometheusAddrPort  = ":9108"
	PrometheusNamespace = "hamtraffic"
)

type Config struct {
	Freerun       bool
	Runtime       *time.Duration
	StationCount  int
	BandModePairs []BandModePair
	Stickiness    float64
}

func NewConfig() *Config {
	var (
		config Config
	)

	// Running mode
	fr := os.Getenv("FREERUN")
	if fr == "" {
		config.Freerun = DefaultFreerun
	} else {
		f, err := strconv.ParseBool(fr)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		config.Freerun = f
	}

	// Running time
	if !config.Freerun {
		rt := os.Getenv("RUNTIME")
		if rt == "" {
			r, _ := time.ParseDuration(DefaultRuntime)
			config.Runtime = &r
		} else {
			r, err := time.ParseDuration(rt)
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			config.Runtime = &r
		}
	} else {
		config.Runtime = nil
	}

	// Station count
	sc := os.Getenv("STATION_COUNT")
	if sc == "" {
		config.StationCount = DefaultStationCount
	} else {
		c, err := strconv.Atoi(sc)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		config.StationCount = c
	}

	// Combine bands and modes
	bands := os.Getenv("BANDS")
	if bands == "" {
		bands = DefaultBands
	}
	modes := os.Getenv("MODES")
	if modes == "" {
		modes = DefaultModes
	}
	config.BandModePairs = NewBandModePairs(bands, modes)

	// Stickiness
	stick := os.Getenv("STICKINESS")
	if stick == "" {
		config.Stickiness = DefaultStickiness
	} else {
		s, err := strconv.ParseFloat(stick, 64)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		config.Stickiness = s
	}

	log.Info().Any("config", config).Msg("Configuration complete")

	return &config
}
