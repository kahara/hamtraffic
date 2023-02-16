package hamtraffic

import "github.com/rs/zerolog/log"

var world *World

func Init() {
	world = NewWorld()

	for _, locale := range world.Locales {
		log.Info().Msgf("%+v", locale)
	}
}
