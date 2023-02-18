package hamtraffic

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

const (
	Prefix = "X0"
)

var nextSuffix = -1

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
	nextSuffix += 1
	if nextSuffix > MaxStationCount-1 {
		log.Panic().Int("suffix", nextSuffix).Msg("too many suffixes")
	}

	suffix := fmt.Sprintf("%04X", nextSuffix)
	callsign := Prefix
	for _, c := range suffix {
		callsign += func(c rune) string {
			lookup := map[rune]rune{
				'0': 'A',
				'1': 'B',
				'2': 'C',
				'3': 'D',
				'4': 'E',
				'5': 'F',
				'6': 'G',
				'7': 'H',
				'8': 'I',
				'9': 'J',
				'A': 'K',
				'B': 'L',
				'C': 'M',
				'D': 'N',
				'E': 'O',
				'F': 'P',
			}
			return string(lookup[c])
		}(c)
	}

	return &Station{
		Callsing: callsign,
		Bands:    NewBands(),
		Locale:   w.RandomLocale(),
	}
}
