package hamtraffic

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

const (
	Prefix       = "X0"
	WeightWiggle = 0.1
)

var (
	nextSuffix = -1
)

type BandModePair struct {
	Name            string
	Weight          float64
	Band            string
	Mode            string
	CenterFrequency float64
	BandWidth       float64
	ChannelSpacing  float64
	TimeSpacing     time.Duration
}

func NewBandModePairs(bandWeights string, modeWeights string) []BandModePair {

	var (
		pairs []BandModePair
		bands = KeyValueStringSplit(bandWeights)
		modes = KeyValueStringSplit(modeWeights)
	)

	// Have product of bands x modes
Loop:
	for _, b := range bands {
		for _, m := range modes {
			var pair BandModePair

			pair.Name = fmt.Sprintf("%s@%s", m.Key, b.Key)
			pair.Weight = b.Value * m.Value
			pair.Band = b.Key
			pair.Mode = m.Key

			// BANDS=160m:0.25,80m:0.40,40m:0.65,20m:1.0,10m:0.65,6m:0.40,2m:0.25
			// MODES=FT8:1.0,FT4:0.25,CW:0.15

			switch b.Key {
			case "160m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 1841500
				case "FT4":
					continue Loop // No FT4 on 160m
				case "CW":
					pair.CenterFrequency = 1852000
					pair.BandWidth = 28000
				}
			case "80m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 3574500
				case "FT4":
					pair.CenterFrequency = 3576500
				case "CW":
					pair.CenterFrequency = 3605000
					pair.BandWidth = 35000
				}

			case "40m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 7075500
				case "FT4":
					pair.CenterFrequency = 7049000
				case "CW":
					pair.CenterFrequency = 7020000
					pair.BandWidth = 40000
				}

			case "20m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 14075500
				case "FT4":
					pair.CenterFrequency = 14081500
				case "CW":
					pair.CenterFrequency = 14035000
					pair.BandWidth = 70000
				}

			case "10m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 28075500
				case "FT4":
					pair.CenterFrequency = 28181500
				case "CW":
					pair.CenterFrequency = 28035000
					pair.BandWidth = 70000
				}

			case "6m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 50314500
				case "FT4":
					pair.CenterFrequency = 50319500
				case "CW":
					pair.CenterFrequency = 50200000
					pair.BandWidth = 200000
				}

			case "2m":
				switch m.Key {
				case "FT8":
					pair.CenterFrequency = 144175500
				case "FT4":
					pair.CenterFrequency = 144171500
				case "CW":
					pair.CenterFrequency = 144212500
					pair.BandWidth = 375000
				}
			}

			switch m.Key {
			case "FT8":
				pair.BandWidth = 3000
				pair.ChannelSpacing = 50
				pair.TimeSpacing = time.Duration(15 * time.Second)
			case "FT4":
				pair.BandWidth = 3000
				pair.ChannelSpacing = 83.3
				pair.TimeSpacing = time.Duration(7500 * time.Millisecond)
			case "CW":
				pair.ChannelSpacing = 100
				pair.TimeSpacing = time.Duration(30 * time.Second)
			}

			pairs = append(pairs, pair)
		}
	}

	return pairs
}

type Station struct {
	Callsign            string
	Antenna             string
	BandModePairs       []BandModePair
	CurrentBandModePair *BandModePair
	Locale              *Locale
}

func NewStation(config *Config, w *World) *Station {
	nextSuffix += 1
	if nextSuffix > MaxStationCount-1 {
		log.Fatal().Int("suffix", nextSuffix).Msg("too many suffixes, passing out")
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

	antenna := [4]string{"Dipole", "Vertical", "Random wire", "Bedsprings"}[rand.Intn(4)]

	var bandModePairs []BandModePair
	for _, pair := range config.BandModePairs {
		wiggle := rand.Float64() * WeightWiggle
		bandModePairs = append(bandModePairs, BandModePair{
			Name:            pair.Name,
			Weight:          pair.Weight + (-wiggle + (wiggle * 2.0)),
			Band:            pair.Band,
			Mode:            pair.Mode,
			CenterFrequency: pair.CenterFrequency,
			BandWidth:       pair.BandWidth,
			ChannelSpacing:  pair.ChannelSpacing,
			TimeSpacing:     pair.TimeSpacing,
		})
	}

	station := Station{
		Callsign:      callsign,
		Antenna:       antenna,
		BandModePairs: bandModePairs,
		Locale:        w.RandomLocale(),
	}

	station.PickBandModePair()

	return &station
}

func (s *Station) Run(done <-chan bool, ack chan<- bool) {
Loop:
	for {
		select {
		case <-done:
			break Loop
		default:
			time.Sleep(time.Second)
		}
	}
	log.Debug().Str("callsign", s.Callsign).Msg("Station shutting down")
	ack <- true
}

func (s *Station) Tick() {
	if rand.Float64() > config.Stickiness {
		s.PickBandModePair()
		log.Debug().Str("callsign", s.Callsign).Str("bandmodepair", s.CurrentBandModePair.Name).Msg("Station changed band and mode")
	}
}

func (s *Station) PickBandModePair() {
	var (
		highscore float64 = 0
		selected  BandModePair
	)

	for _, pair := range s.BandModePairs {
		score := rand.Float64() * pair.Weight
		if score > highscore {
			highscore = score
			selected = pair
		}
	}

	s.CurrentBandModePair = &selected
}
