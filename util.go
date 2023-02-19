package hamtraffic

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type KeyValue struct {
	Key   string
	Value float64
}

func KeyValueStringSplit(s string) []KeyValue {
	var pairs []KeyValue

	for _, pair := range strings.Split(s, ",") {
		kv := strings.Split(pair, ":")
		value, err := strconv.ParseFloat(kv[1], 64)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		pairs = append(pairs, KeyValue{
			Key:   kv[0],
			Value: value,
		})
	}

	return pairs
}
