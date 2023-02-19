package hamtraffic

import (
	"github.com/logocomune/maidenhead"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

const CityDataPath = "data/cities.geojson"

type Locale struct {
	Name       string
	Population int
	Geometry   orb.Point
	Locators   []string
}

type World struct {
	Locales []Locale
}

func NewWorld() *World {
	return &World{
		Locales: loadCities(),
	}
}

func (w *World) RandomLocale() *Locale {
	return &w.Locales[rand.Intn(len(w.Locales))]
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

		// Locator
		var (
			locator  string
			locators []string
		)

		locator, _ = maidenhead.Locator(point.Lat(), point.Lon(), maidenhead.FieldPrecision)
		locators = append(locators, locator)

		locator, _ = maidenhead.Locator(point.Lat(), point.Lon(), maidenhead.SquarePrecision)
		locators = append(locators, locator)

		locator, _ = maidenhead.Locator(point.Lat(), point.Lon(), maidenhead.SubSquarePrecision)
		locators = append(locators, locator)

		locator, _ = maidenhead.Locator(point.Lat(), point.Lon(), maidenhead.ExtendedSquarePrecision)
		locators = append(locators, locator)

		locator, _ = maidenhead.Locator(point.Lat(), point.Lon(), maidenhead.SubExtendedSquarePrecision)
		locators = append(locators, locator)

		locales = append(locales, Locale{
			Name:       name,
			Population: pop,
			Geometry:   point,
			Locators:   locators,
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
