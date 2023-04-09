package hamtraffic

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

type Runlog struct {
	f *os.File
}

func NewRunlog(path string) *Runlog {
	var (
		err    error
		f      *os.File
		runlog Runlog
	)

	if f, err = os.OpenFile(config.RunlogPath, os.O_WRONLY|os.O_CREATE, 0664); err != nil {
		log.Fatal().Err(err).Str("path", config.RunlogPath).Msg("")
	}

	runlog.f = f

	return &runlog
}

func (r *Runlog) Log(item interface{}) {
	log.Trace().Str("item", fmt.Sprintf("%+v", item)).Msg("")
}
