package hamtraffic

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Runlog struct {
	f *os.File
}

type Logitem struct {
	Receiver struct {
		Callsign string `json:"callsign"`
		Locator  string `json:"locator"`
	} `json:"receiver"`

	Spot struct {
		Time      time.Time `json:"time"`
		Frequency float64   `json:"frequency"`
		Callsign  string    `json:"callsign"`
		Locator   string    `json:"locator"`
	} `json:"spot"`
}

func NewRunlog(path string) *Runlog {
	var (
		err    error
		f      *os.File
		runlog Runlog
	)

	if path != "" {
		if f, err = os.OpenFile(config.RunlogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664); err != nil {
			log.Fatal().Err(err).Str("path", config.RunlogPath).Msg("")
		}
		runlog.f = f
	}

	return &runlog
}

func (r *Runlog) Log(station *Station, transmission *Transmission) {
	if r.f == nil {
		return
	}

	var (
		logitem Logitem
	)

	logitem.Receiver.Callsign = station.Callsign
	logitem.Receiver.Locator = station.Locale.Locators[2]
	logitem.Spot.Time = transmission.Time
	logitem.Spot.Frequency = transmission.Frequency
	logitem.Spot.Callsign = transmission.Station.Callsign
	logitem.Spot.Locator = transmission.Station.Locale.Locators[1]

	if buf, err := json.Marshal(logitem); err != nil {
		log.Err(err).Msg("Runlog marshalling failed")
	} else {
		if _, err := r.f.Write(buf); err != nil {
			log.Err(err).Msg("Runlog writing failed")
		} else {
			_, _ = r.f.Write([]byte{'\n'})
		}
	}
}
