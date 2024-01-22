package config

import "github.com/caarlos0/env"

type Config struct {
	File     string   `env:"FILE" envDefault:"/app/files/fixtures.json"`
	Words    []string `env:"WORDS" envDefault:"Potato,Mushroom,Veggie"`
	Postcode string   `env:"POSTCODE" envDefault:"10120"`
	FromTime string   `env:"FROM" envDefault:"10AM"`
	ToTime   string   `env:"TO" envDefault:"3PM"`
}

func ReadConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}

func (c Config) WithFile(file string) Config {
	c.File = file
	return c
}

func (c Config) WithWords(words []string) Config {
	c.Words = words
	return c
}

func (c Config) WithPostcode(postcode string) Config {
	c.Postcode = postcode
	return c
}

func (c Config) WithFromTime(fromTime string) Config {
	c.FromTime = fromTime
	return c
}

func (c Config) WithToTime(toTime string) Config {
	c.ToTime = toTime
	return c
}
