package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rashad-j/jsonreader/pkg/config"
	"github.com/rashad-j/jsonreader/pkg/parser"
	"github.com/rashad-j/jsonreader/pkg/stats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	fileName string
	fromTime string
	toTime   string
	postcode string
	words    string
	helpFlag bool
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"})
}

func ExecuteStatsCMD() error {
	var statsCmd = &cobra.Command{
		Use:     "stats",
		Short:   "Generate Recipes statistics based on specified parameters",
		Run:     runStats,
		Example: `./parser stats --file ./files/test.json --postcode 10120 --words Potato,Mushroom,Veggie --fromTime 10AM --toTime 3PM`,
	}

	// get default values from config
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	statsCmd.Flags().StringVarP(&fileName, "file", "f", cfg.File, "File to use (optional)")
	statsCmd.Flags().StringVarP(&fromTime, "fromTime", "s", cfg.FromTime, "From time (optional)")
	statsCmd.Flags().StringVarP(&toTime, "toTime", "e", cfg.ToTime, "To time (optional)")
	statsCmd.Flags().StringVarP(&postcode, "postcode", "p", cfg.Postcode, "Postcode (required)")
	statsCmd.Flags().StringVarP(&words, "words", "w", strings.Join(cfg.Words, ","), "List of comma-separated words (optional)")
	statsCmd.Flags().BoolVarP(&helpFlag, "help", "h", false, "Show help information")

	if err := statsCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func runStats(cmd *cobra.Command, args []string) {
	log.Info().Msg("Calculating stats...")
	if helpFlag {
		cmd.Help()
		return
	}

	// sanitize parameters
	// postcode is less than 10 chars
	if len(postcode) > 10 || len(postcode) == 0 {
		fmt.Println("postcode must be less than 10 characters")
		return
	}
	// trim whitespace
	tirmmedWords := strings.TrimSpace(words)
	wordSlice := strings.Split(tirmmedWords, ",")
	// remove empty strings
	var words []string
	for _, word := range wordSlice {
		if word != "" {
			words = append(words, word)
		}
	}

	// Additional logic can be added to process the parameters as needed
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	// NOTE: config uses the builder pattern
	// Check if optional parameters are set, if yes, override the default values
	if fileName != cfg.File {
		cfg = cfg.WithFile(fileName)
	}
	if fromTime != cfg.FromTime {
		cfg = cfg.WithFromTime(fromTime)
	}
	if toTime != cfg.ToTime {
		cfg = cfg.WithToTime(toTime)
	}
	if postcode != cfg.Postcode {
		cfg = cfg.WithPostcode(postcode)
	}
	if len(words) > 0 {
		cfg = cfg.WithWords(words)
	}

	// Create JsonParser object - it implements the Parser interface
	p := parser.NewJsonParser(cfg)
	go p.Parse()
	// Create stats object
	s := stats.NewJsonStats(p, cfg)
	// generate stats
	data, err := s.Generate()
	if err != nil {
		fmt.Println("Error generating stats:", err)
		return
	}

	// Marshal the ResponseData to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// print the result
	fmt.Println(string(jsonData)) // this is piped to stdout
	log.Info().Msg("Done!")       // this is piped to stderr
}
