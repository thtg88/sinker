package config

import "os"

type Config struct{
	Sinker *Sinker
}

type Sinker struct{
	BasePath string
}

func Load() *Config {
	return &Config{
		Sinker: &Sinker{
			BasePath: os.Getenv("SINKER_BASE_PATH"),
		},
	}
}
