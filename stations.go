package hamtraffic

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

const (
	Prefix = "X0"
)

var nextSuffix = 0

type Band struct {
	Name            string
	CenterFrequency float64
	Bandwidth       float64
	Bias            float64
}

func NewBands() []*Band {
	var bands []*Band

	return bands
}

type Station struct {
	Callsing string
	Bands    []*Band
	Locale   *Locale
}

func NewStation(w *World) *Station {
	if nextSuffix > 9999 {
		log.Panic().Int("suffix", nextSuffix).Msg("too many suffixes")
	}

	callsign := fmt.Sprintf("%s%04d", Prefix, nextSuffix)

	nextSuffix += 1

	return &Station{
		Callsing: callsign,
		Bands:    NewBands(),
		Locale:   w.RandomLocale(),
	}
}
