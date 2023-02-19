package main

import (
	"github.com/kahara/hamtraffic"
	"github.com/rs/zerolog"
	"math/rand"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	rand.Seed(time.Now().UnixNano())
	hamtraffic.Init()
}
