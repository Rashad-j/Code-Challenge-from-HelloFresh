package stats

import (
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rashad-j/jsonreader/pkg/config"
	"github.com/rashad-j/jsonreader/pkg/parser"
	"github.com/rs/zerolog/log"
)

type Stats interface {
	Generate() (ResponseData, error)
}

type JsonStats struct {
	parser parser.Parser
	cfg    config.Config
}

func NewJsonStats(p parser.Parser, cfg config.Config) *JsonStats {
	return &JsonStats{
		parser: p,
		cfg:    cfg,
	}
}

func (s *JsonStats) Generate() (ResponseData, error) {
	// Initialize maps and slices for data analysis
	recipeCounts := make(map[string]int, 2000)
	postCodeCounts := make(map[string]int, 1000_000)
	postCodeMaxDeliveries := ""
	specificPostCodeDeliveries := 0
	recipesContainingWords := make(map[string]int)

	// Convert words from slice to map for faster lookup
	wordsMap := make(map[string]bool, len(s.cfg.Words))
	for _, word := range s.cfg.Words {
		wordsMap[strings.ToLower(word)] = true
	}

	// Read json content over stream
	for entry := range s.parser.Stream() {
		if entry.Error != nil {
			log.Error().Err(entry.Error).Msg("failed to process entry")
			continue
		}

		// This is to count the number of unique recipes, and the total number of recipes
		recipeCounts[entry.Recipe.Recipe]++
		postCodeCounts[entry.Recipe.Postcode]++

		// Find postcode with most delivered recipes
		if postCodeCounts[entry.Recipe.Postcode] > postCodeCounts[postCodeMaxDeliveries] {
			postCodeMaxDeliveries = entry.Recipe.Postcode
		}

		// Find recipes containing words
		if s.containsWords(entry.Recipe.Recipe, wordsMap) {
			recipesContainingWords[entry.Recipe.Recipe]++
		}
		// Number of deliveries for postcode and time range
		if entry.Recipe.Postcode == s.cfg.Postcode {
			inRange, err := s.isDeliveryTimeInRange(entry.Recipe.Delivery, s.cfg.FromTime, s.cfg.ToTime)
			if err != nil {
				return ResponseData{}, errors.Wrapf(err, "failed to check if delivery time is in range: %s", entry.Recipe.Delivery)
			}
			if inRange {
				specificPostCodeDeliveries++
			}
		}
	}

	// sort recipe names alphabetically
	sortedKeys := sortKeys(recipeCounts)
	countPerRecipe := s.uniqueRecipeCount(sortedKeys, recipeCounts)

	// sort recipes containing words alphabetically
	matchByName := sortKeys(recipesContainingWords)

	responseData := ResponseData{
		UniqueRecipeCount: len(recipeCounts),
		CountPerRecipe:    countPerRecipe,
		BusiestPostcode: BusiestPostcode{
			Postcode:      postCodeMaxDeliveries,
			DeliveryCount: postCodeCounts[postCodeMaxDeliveries],
		},
		CountPerPostcodeAndTime: CountPerPostcodeAndTime{
			Postcode:      s.cfg.Postcode,
			From:          s.cfg.FromTime,
			To:            s.cfg.ToTime,
			DeliveryCount: specificPostCodeDeliveries,
		},
		MatchByName: matchByName,
	}

	return responseData, nil
}

// containsWords checks if the recipe contains any of the words
func (s *JsonStats) containsWords(recipe string, words map[string]bool) bool {
	recipeWords := strings.Fields(recipe)
	// check if recipe contains any of the words
	for _, recipeWord := range recipeWords {
		if words[strings.ToLower(recipeWord)] {
			return true
		}
	}

	return false
}

// isDeliveryTimeInRange checks if the delivery time is within the specified range.
// from startHour, but not including endHour.
func (s *JsonStats) isDeliveryTimeInRange(delivery string, startHour, endHour string) (bool, error) {
	// delivery format: Monday 9AM - 5PM
	// startHour is in the format 10AM, 10PM, etc. The same endHour is in the format 3PM, 1AM, etc.
	startHourInt, err := parseHour(startHour)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse start hour: %s", startHour)
	}
	endHourInt, err := parseHour(endHour)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse end hour: %s", endHour)
	}

	// split delivery string into day and time
	deliverySplit := strings.SplitN(delivery, " ", 2)
	if len(deliverySplit) != 2 {
		return false, errors.New("failed to split delivery string")
	}
	// split time into start and end time
	deliveryTimeSplit := strings.SplitN(deliverySplit[1], " - ", 2)
	if len(deliveryTimeSplit) != 2 {
		return false, errors.New("failed to split delivery time")
	}
	// parse split start and end time into int
	deliveryStartHourInt, err := parseHour(deliveryTimeSplit[0])
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse delivery start hour: %s", deliveryTimeSplit[0])
	}
	deliveryEndHourInt, err := parseHour(deliveryTimeSplit[1])
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse delivery end hour: %s", deliveryTimeSplit[1])
	}

	// check if delivery time is within range
	if deliveryStartHourInt <= startHourInt && endHourInt < deliveryEndHourInt {
		return true, nil
	}

	return false, nil
}

// uniqueRecipeCount returns a slice of RecipeCount objects alphabetically sorted by recipe name
func (s *JsonStats) uniqueRecipeCount(sortedKeys []string, recipeCounts map[string]int) []RecipeCount {
	// create RecipeCount objects
	recipeCountSlice := make([]RecipeCount, 0)
	for _, k := range sortedKeys {
		recipeCountSlice = append(recipeCountSlice, RecipeCount{
			Recipe: k,
			Count:  recipeCounts[k],
		})
	}

	return recipeCountSlice
}

// parseHour a helper function to parse the hour from the delivery time
func parseHour(hourString string) (int, error) {
	// remove AM/PM
	hourTrimmed := strings.TrimSuffix(hourString, "AM")
	hourTrimmed = strings.TrimSuffix(hourTrimmed, "PM")

	// convert to int
	hourInt, err := strconv.Atoi(hourTrimmed)
	if err != nil {
		return 0, err
	}

	// check the range between 1 and 12
	if hourInt < 1 || hourInt > 12 {
		return 0, errors.New("hour is not in range 1-12")
	}

	// 12AM is 0
	if strings.Contains(hourString, "AM") && hourInt == 12 {
		hourInt = 0
	}

	// if PM, add 12 hours
	if strings.Contains(hourString, "PM") {
		hourInt += 12
	}

	return hourInt, nil
}

// uniqueRecipeCount a helper function to sort the keys alphabetically
func sortKeys(m map[string]int) []string {
	// extract recipe names and sort alphabetically
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
