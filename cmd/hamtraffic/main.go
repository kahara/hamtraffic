package main

import (
	"github.com/kahara/hamtraffic"
	"github.com/rs/zerolog"
	"math/rand"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	rand.Seed(time.Now().UnixNano())
	hamtraffic.Init()
}
