package stats

// Unit tests for the stats package

import (
	"reflect"
	"testing"

	"github.com/rashad-j/jsonreader/pkg/config"
	"github.com/rashad-j/jsonreader/pkg/parser"
)

func TestNewJsonStats(t *testing.T) {
	// arrange
	parser := parser.NewJsonParser(config.Config{})
	config := config.Config{}

	// act
	jsonStats := NewJsonStats(parser, config)

	// assert
	if got := NewJsonStats(parser, config); !reflect.DeepEqual(got, jsonStats) {
		t.Errorf("NewJsonStats() = %v, want %v", got, jsonStats)
	}
}

func TestJsonStats_ContainsWords(t *testing.T) {
	// arrange
	parser := parser.NewJsonParser(config.Config{})
	config := config.Config{}
	jsonStats := NewJsonStats(parser, config)
	tests := []struct {
		name   string
		recipe string
		words  map[string]bool
		want   bool
	}{
		{
			name:   "recipe contains words",
			recipe: "Hot Honey Barbecue Mushroom Chicken Legs",
			words: map[string]bool{
				"potato":   true,
				"mushroom": true,
				"veggie":   true,
			},
			want: true,
		},
		{
			name:   "recipe does not contain words",
			recipe: "Grilled Cheese Jumble",
			words: map[string]bool{
				"potato":   true,
				"mushroom": true,
				"veggie":   false,
			},
			want: false,
		},
		{
			name:   "recipe contains words with different case",
			recipe: "Korean-Style Chicken Potato Thighs",
			words: map[string]bool{
				"potato":   true,
				"mushroom": true,
				"veggie":   true,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		// avoid closure
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// act
			if got := jsonStats.containsWords(tt.recipe, tt.words); got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestJsonStats_IsDeliveryTimeInRange(t *testing.T) {
	// arrange
	parser := parser.NewJsonParser(config.Config{})
	config := config.Config{}
	jsonStats := NewJsonStats(parser, config)
	tests := []struct {
		name      string
		delivery  string
		startHour string
		endHour   string
		want      bool
	}{
		{
			name:      "delivery time is within range",
			delivery:  "Monday 9AM - 5PM",
			startHour: "10AM",
			endHour:   "3PM",
			want:      true,
		},
		{
			name:      "delivery time is not within range",
			delivery:  "Monday 11AM - 5PM",
			startHour: "10AM",
			endHour:   "3PM",
			want:      false,
		},
		{
			name:      "delivery time is not correct format",
			delivery:  "Monday 99AM - 5PM",
			startHour: "10PM",
			endHour:   "3AM",
			want:      false,
		},
		{
			name:      "start delivery time is 12AM",
			delivery:  "Monday 12AM - 11PM",
			startHour: "8AM",
			endHour:   "9AM",
			want:      true,
		},
	}

	for _, tt := range tests {
		// avoid closure
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			// act
			if got, _ := jsonStats.isDeliveryTimeInRange(testCase.delivery, testCase.startHour, testCase.endHour); got != testCase.want {
				t.Errorf("%s = %v, want %v", testCase.name, got, testCase.want)
			}
		})
	}
}

func TestJsonStats_UniqueRecipeCount(t *testing.T) {
	tests := []struct {
		name          string
		sortedKeys    []string
		recipeCounts  map[string]int
		expectedCount []RecipeCount
	}{
		{
			name:       "Normal Case",
			sortedKeys: []string{"RecipeA", "RecipeB", "RecipeC"},
			recipeCounts: map[string]int{
				"RecipeC": 8,
				"RecipeB": 5,
				"RecipeA": 10,
			},
			expectedCount: []RecipeCount{
				{Recipe: "RecipeA", Count: 10},
				{Recipe: "RecipeB", Count: 5},
				{Recipe: "RecipeC", Count: 8},
			},
		},
		{
			name:       "Empty Case",
			sortedKeys: []string{},
			recipeCounts: map[string]int{
				"RecipeC": 8,
				"RecipeB": 5,
				"RecipeA": 10,
			},
			expectedCount: []RecipeCount{},
		},
		{
			name:       "Missing Recipe Case",
			sortedKeys: []string{"RecipeA", "RecipeB", "RecipeC"},
			recipeCounts: map[string]int{
				"RecipeC": 8,
				"RecipeB": 5,
			},
			expectedCount: []RecipeCount{
				{Recipe: "RecipeA", Count: 0},
				{Recipe: "RecipeB", Count: 5},
				{Recipe: "RecipeC", Count: 8},
			},
		},
	}

	for _, tt := range tests {
		// avoid closure
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stats := &JsonStats{}
			result := stats.uniqueRecipeCount(tt.sortedKeys, tt.recipeCounts)
			if !reflect.DeepEqual(result, tt.expectedCount) {
				t.Errorf("Expected %v, but got %v", tt.expectedCount, result)
			}
		})
	}
}

func TestJsonStats_ParseHour(t *testing.T) {
	tests := []struct {
		name     string
		hour     string
		expected int
	}{
		{
			name:     "Normal Case AM",
			hour:     "9AM",
			expected: 9,
		},
		{
			name:     "Normal Case PM",
			hour:     "9PM",
			expected: 21,
		},
		{
			name:     "Normal Case 12PM",
			hour:     "12PM",
			expected: 24,
		},
		{
			name:     "Normal Case 12AM",
			hour:     "12AM",
			expected: 0,
		},
		{
			name:     "Normal Case 1AM",
			hour:     "1AM",
			expected: 1,
		},
		{
			name:     "Normal Case 1PM",
			hour:     "1PM",
			expected: 13,
		},
		{
			name:     "Normal Case",
			hour:     "11AM",
			expected: 11,
		},
	}

	for _, tt := range tests {
		// avoid closure
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHour(tt.hour)
			if err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
