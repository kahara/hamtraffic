package main

import (
	"github.com/kahara/hamtraffic"
	"github.com/rs/zerolog"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	hamtraffic.Init()
}
