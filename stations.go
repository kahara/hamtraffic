package hamtraffic

import (
	"fmt"
	"github.com/paulmach/orb/geo"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

const (
	Prefix                 = "X0"
	WeightWiggle           = 0.1
	SmallestCommonDuration = time.Duration(500 * time.Millisecond)
)

var (
	nextSuffix    = -1
	NeighbourBins = []float64{40960000, 20480000, 10240000, 5210000, 2560000, 1280000, 640000, 320000, 160000, 80000, 40000, 20000, 10000, 0}
)

type BandModePair struct {
	Name             string
	Weight           float64
	Band             string
	Mode             string
	CenterFrequency  float64
	BandWidth        float64
	ChannelSpacing   float64
	TransmitDuration time.Duration
	TransmitPeriods  []time.Duration
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
				pair.TransmitDuration = time.Duration(15 * time.Second)
				pair.TransmitPeriods = []time.Duration{
					time.Duration(0 * time.Second),
					time.Duration(15 * time.Second),
					time.Duration(30 * time.Second),
					time.Duration(45 * time.Second),
				}
			case "FT4":
				pair.BandWidth = 3000
				pair.ChannelSpacing = 83.3
				pair.TransmitDuration = time.Duration(7500 * time.Millisecond)
				pair.TransmitPeriods = []time.Duration{
					time.Duration(0 * time.Millisecond),
					time.Duration(7500 * time.Millisecond),
					time.Duration(15000 * time.Millisecond),
					time.Duration(22500 * time.Millisecond),
					time.Duration(30000 * time.Millisecond),
					time.Duration(37500 * time.Millisecond),
					time.Duration(45000 * time.Millisecond),
					time.Duration(52500 * time.Millisecond),
				}
			case "CW":
				pair.ChannelSpacing = 100
				pair.TransmitDuration = time.Duration(30 * time.Second)
				pair.TransmitPeriods = []time.Duration{
					time.Duration(0 * time.Second),
					time.Duration(30 * time.Second),
				}
			}

			pairs = append(pairs, pair)
		}
	}

	return pairs
}

type Transmission struct {
	Station   *Station
	Time      time.Time
	Duration  time.Duration
	Frequency float64
	Band      string
	Mode      string
	Power     float64
}

type Station struct {
	Callsign            string
	Antenna             string
	BandModePairs       []BandModePair
	CurrentBandModePair *BandModePair
	Frequency           float64
	TransmitEven        bool
	TransmitPeriods     []time.Duration
	Locale              *Locale
	Neighbours          map[float64][]*Station
	Receiver            chan Transmission
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
			TransmitPeriods: pair.TransmitPeriods,
		})
	}

	station := Station{
		Callsign:      callsign,
		Antenna:       antenna,
		BandModePairs: bandModePairs,
		Locale:        w.RandomLocale(),
		Receiver:      make(chan Transmission, 100),
	}

	if rand.Float64() < 0.5 {
		station.TransmitEven = true
	} else {
		station.TransmitEven = false
	}

	station.PickBandModePair()

	return &station
}

func (s *Station) ComputeNeighbourhood(neighbours []*Station) {
	var g = s.Locale.Geometry

	s.Neighbours = make(map[float64][]*Station)
	for _, bin := range NeighbourBins {
		s.Neighbours[bin] = []*Station{}
	}

	for _, neighbour := range neighbours {
		// Skip self
		if neighbour == s {
			continue
		}

		// Put each remote station into its distance bin, wasting space while conserving time during lookups
		distance := geo.DistanceHaversine(g, neighbour.Locale.Geometry)
		for index, bin := range NeighbourBins {
			// Skip the last (zero) bin
			if bin == 0 {
				break
			}
			if distance < bin && distance > NeighbourBins[index+1] {
				s.Neighbours[bin] = append(s.Neighbours[bin], neighbour)
				break
			}
		}
	}

	log.Debug().Str("callsign", s.Callsign).Any("bins", func() map[string]int {
		x := make(map[string]int)
		for bin, n := range s.Neighbours {
			x[fmt.Sprintf("%.0fkm", bin/1000)] = len(n)
		}
		return x
	}()).Msg("Computed distances to neighbours")
}

func (s *Station) PickBandModePair() {
	var (
		highscore float64 = 0
		selected  BandModePair
	)

	// Pick a (weightedly) random band-mode-pair
	for _, pair := range s.BandModePairs {
		score := rand.Float64() * pair.Weight
		if score > highscore {
			highscore = score
			selected = pair
		}
	}
	s.CurrentBandModePair = &selected

	// Set the frequency
	channels := int(s.CurrentBandModePair.BandWidth / s.CurrentBandModePair.ChannelSpacing)
	floor := s.CurrentBandModePair.CenterFrequency - (s.CurrentBandModePair.BandWidth / 2)
	s.Frequency = floor + (s.CurrentBandModePair.ChannelSpacing * float64(rand.Intn(channels)))

	// Resolve the even, odd periods
	s.TransmitPeriods = []time.Duration{}
	for index, period := range s.CurrentBandModePair.TransmitPeriods {
		if ((index%2) == 0 && s.TransmitEven) || ((index%2) == 1 && !s.TransmitEven) {
			s.TransmitPeriods = append(s.TransmitPeriods, period)
		}
	}
}

func (s *Station) Receive(transmission Transmission) {
	//log.Debug().Str("sender", transmission.Station.Callsign).Str("receiver", s.Callsign).Str("band", transmission.Band).Str("mode", transmission.Mode).Msg("Transmission received")

	// TODO check if station was listening

	// TODO check if concurrent increments are an actual problem
	metrics["receptions"].WithLabelValues(transmission.Band, transmission.Mode, s.Callsign).Inc()
}

func (s *Station) Run(start *time.Time, xmit chan<- Transmission, done <-chan bool, ack chan<- bool) {
	var t time.Time

	time.Sleep(time.Until(*start))
	ticker := time.NewTicker(SmallestCommonDuration)

Loop:
	for {
		select {
		case <-done:
			break Loop
		case transmission := <-s.Receiver:
			s.Receive(transmission)
		default:
			t = <-ticker.C
		}

		// Maybe decide to change band and mode
		if rand.Float64() > config.Stickiness {
			s.PickBandModePair()
			log.Debug().Str("callsign", s.Callsign).Str("bandmodepair", s.CurrentBandModePair.Name).Msg("Station changed band and mode, and is waiting until beginning of next minute before proceeding transmission")
			time.Sleep(time.Until(t.Truncate(time.Minute).Add(time.Duration(time.Minute))))
		}

		for _, period := range s.TransmitPeriods {
			if t.Sub(t.Truncate(time.Duration(time.Minute)).Add(period)).Abs() < (SmallestCommonDuration / 2) {
				// Decide to transmit or not
				if rand.Float64() > config.TransmissionProbability {
					log.Debug().Str("callsign", s.Callsign).Msg("Skipping transmission")
					break
				}
				xmit <- Transmission{
					Station:   s,
					Time:      time.Now().UTC(),
					Duration:  s.CurrentBandModePair.TransmitDuration,
					Frequency: s.Frequency,
					Band:      s.CurrentBandModePair.Band,
					Mode:      s.CurrentBandModePair.Mode,
					Power:     0,
				}
				log.Debug().Time("time", t).Str("callsign", s.Callsign).Str("bandmodepair", s.CurrentBandModePair.Name).Bool("even", s.TransmitEven).Dur("period", period).Msg("Transmitting")
			}
		}
	}

	log.Debug().Str("callsign", s.Callsign).Msg("Station shutting down")
	ack <- true
}
