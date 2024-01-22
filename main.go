package main

import (
	"github.com/rs/zerolog/log"

	"github.com/rashad-j/jsonreader/cmd/stats"
)

func main() {
	if err := stats.ExecuteStatsCMD(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute stats command")
	}
}
