package hamtraffic

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/rs/zerolog/log"
	"os"
	"sort"
	"strconv"
)

const CityDataPath = "data/export.geojson"

type World struct {
	Locales []Locale
}

type Locale struct {
	Name       string
	Population int
	Geometry   orb.Point
}

func NewWorld() *World {
	return &World{
		Locales: loadCities(),
	}
}

func loadCities() []Locale {
	var (
		locales []Locale
	)

	data, err := os.ReadFile(CityDataPath)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	fc, _ := geojson.UnmarshalFeatureCollection(data)

	count := 0
	for _, feature := range fc.Features {
		// Name
		name := feature.Properties.MustString("name:en", "")
		if name == "" {
			name = feature.Properties.MustString("name", "")
			if name == "" {
				continue
			}
		}

		// Population
		p := feature.Properties.MustString("population", "")
		if p == "" {
			continue
		}
		pop, err := strconv.Atoi(p)
		if err != nil {
			continue
		}

		// Geometry
		point := feature.Geometry.(orb.Point)

		locales = append(locales, Locale{
			Name:       name,
			Population: pop,
			Geometry:   point,
		})

		log.Trace().
			Str("name", name).
			Int("population", pop).
			Any("point", point).
			Msg("")
		count += 1
	}

	sort.Slice(locales, func(i, j int) bool {
		return locales[i].Population < locales[j].Population
	})

	log.Info().Int("cities", count).Msg("City loading complete")

	return locales
}
