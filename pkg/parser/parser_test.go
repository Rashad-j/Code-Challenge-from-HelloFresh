package parser

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rashad-j/jsonreader/pkg/config"
)

func TestJsonParser_Parse(t *testing.T) {
	var errInvalidTimeFormat error = errors.New("invalid time format")
	tests := []struct {
		name        string
		fileContent string
		expected    []Entry
	}{
		{
			name:        "Valid JSON Content",
			fileContent: `[{"Postcode": "12345", "Delivery": "Monday 9AM - 5PM", "Recipe": "RecipeA"}]`,
			expected: []Entry{
				{Recipe: Recipe{Postcode: "12345", Delivery: "Monday 9AM - 5PM", Recipe: "RecipeA"}},
			},
		},
		{
			name:        "Valid JSON Content with multiple entries",
			fileContent: `[{"Postcode": "12345", "Delivery": "Monday 9AM - 5PM", "Recipe": "RecipeA"}, {"Postcode": "12345", "Delivery": "Monday 10AM - 6PM", "Recipe": "RecipeB"}]`,
			expected: []Entry{
				{Recipe: Recipe{Postcode: "12345", Delivery: "Monday 9AM - 5PM", Recipe: "RecipeA"}},
				{Recipe: Recipe{Postcode: "12345", Delivery: "Monday 10AM - 6PM", Recipe: "RecipeB"}},
			},
		},
		{
			name:        "Invalid JSON Content",
			fileContent: `[{"Postcode": "12345", "Delivery": "InvalidTimeFormat", "Recipe": "RecipeA"}]`,
			expected:    []Entry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating a temporary JSON file and write content
			file, err := createTempJSONFile(t, tt.fileContent)
			if err != nil {
				t.Fatalf("Error creating temporary file: %v", err)
			}
			defer os.Remove(file.Name())

			cfg := config.Config{File: file.Name()}
			parser := NewJsonParser(cfg)

			// Start JSON stream parsing
			go parser.Parse()

			// Creating a channel to receive entries
			entryCh := make(chan Entry, len(tt.expected))
			go func() {
				defer close(entryCh)
				for entry := range parser.Stream() {
					entryCh <- entry
				}
			}()

			time.Sleep(time.Second)

			// Collect entries from the channel
			var actual []Entry
			for entry := range entryCh {
				// Check if error is expected
				if entry.Error != nil && entry.Error.Error() == errInvalidTimeFormat.Error() {
					continue
				}
				actual = append(actual, entry)
			}

			// Compare the result with the expected output
			if !entriesEqual(actual, tt.expected) {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}

// Helper function to create a temporary JSON file with content
func createTempJSONFile(t *testing.T, content string) (*os.File, error) {
	file, err := os.CreateTemp("./", "temp-*.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Helper function to compare slices of Entry, ignoring errors
func entriesEqual(actual, expected []Entry) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i := range actual {
		if actual[i].Error != nil || expected[i].Error != nil {
			continue
		}

		if actual[i].Recipe != expected[i].Recipe {
			return false
		}
	}

	return true
}

func TestJsonParser_SanitizeRecipe(t *testing.T) {
	tests := []struct {
		name     string
		recipe   Recipe
		expected bool
	}{
		{
			name:     "Valid Recipe",
			recipe:   Recipe{Postcode: "12345", Delivery: "Monday 9AM - 5PM", Recipe: "RecipeA"},
			expected: true,
		},
		{
			name:     "Empty Recipe",
			recipe:   Recipe{},
			expected: false,
		},
		{
			name: "Long Recipe name",
			recipe: Recipe{
				Postcode: "12345",
				Delivery: "Monday 9AM - 5PM",
				Recipe:   "RecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipeRecipe",
			},
		},
		{
			name:     "Empty Postcode",
			recipe:   Recipe{Delivery: "Monday 9AM - 5PM", Recipe: "RecipeA"},
			expected: false,
		},
		{
			name:     "Empty Delivery",
			recipe:   Recipe{Postcode: "12345", Recipe: "RecipeA"},
			expected: false,
		},
		{
			name:     "Invalid Delivery",
			recipe:   Recipe{Postcode: "12345", Delivery: "InvalidTimeFormat", Recipe: "RecipeA"},
			expected: false,
		},
		{
			name:     "Postcode longer than 10 characters",
			recipe:   Recipe{Postcode: "12345678901", Delivery: "Monday 9AM - 5PM", Recipe: "RecipeA"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewJsonParser(config.Config{})
			result := parser.sanitizeRecipe(tt.recipe)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
