package parser

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/rashad-j/jsonreader/pkg/config"
	"github.com/rs/zerolog/log"
)

type Parser interface {
	Parse()
	Stream() <-chan Entry
}

type JsonParser struct {
	cfg    config.Config
	stream chan Entry
}

func NewJsonParser(cfg config.Config) *JsonParser {
	return &JsonParser{
		cfg:    cfg,
		stream: make(chan Entry),
	}
}

func (r *JsonParser) Stream() <-chan Entry {
	return r.stream
}

// Parse reads the JSON file and streams Recipe objects over the channel
func (r *JsonParser) Parse() {
	defer close(r.stream)

	file, err := os.Open(r.cfg.File)
	if err != nil {
		r.stream <- Entry{Error: errors.Wrap(err, "failed to open file")}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	// read opening delimiter `[`
	if _, err := decoder.Token(); err != nil {
		r.stream <- Entry{Error: errors.Wrap(err, "failed to read opening delimiter")}
	}

	for decoder.More() {
		var recipe Recipe
		// decode an array value (Recipe)
		if err := decoder.Decode(&recipe); err != nil {
			r.stream <- Entry{Error: errors.Wrap(err, "failed to decode recipe")}
		}
		if !r.sanitizeRecipe(recipe) {
			continue
		}
		r.stream <- Entry{Recipe: recipe}
	}

	// read closing delimiter `]`
	if _, err := decoder.Token(); err != nil {
		r.stream <- Entry{Error: errors.Wrap(err, "failed to read closing delimiter")}
	}
}

func (r *JsonParser) sanitizeRecipe(recipe Recipe) bool {
	// postcode is not empty
	if recipe.Postcode == "" {
		log.Error().Msg("postcode is empty")
		return false
	}
	// postcode is less than 10 characters
	if len(recipe.Postcode) > 10 {
		log.Error().Msg("postcode is longer than 10 characters")
		return false
	}

	// delivery is not empty
	if recipe.Delivery == "" {
		log.Error().Msg("delivery is empty")
		return false
	}

	pattern := `^(\w+)\s+([1-9]|1[0-2])\s*(AM)\s*-\s*([1-9]|1[0-2])\s*(PM)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(recipe.Delivery)
	// Check if the format matches
	if len(matches) != 6 {
		log.Error().Str("delivery", recipe.Delivery).Msg("delivery format does not match")
		return false
	}

	// check if recipe is not empty
	if recipe.Recipe == "" {
		log.Error().Msg("recipe is empty")
		return false
	}
	// check that recipe is less than 100 characters
	if len(recipe.Recipe) > 100 {
		log.Error().Msg("recipe is longer than 100 characters")
		return false
	}

	return true
}
