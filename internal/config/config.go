package config

import "os"

type Config struct{
	Sinker 		*Sinker
	SinkerAPI *SinkerAPI
}

type Sinker struct{
	BasePath string
}

type SinkerAPI struct {
	StoreEventPath string
}

func Load() *Config {
	return &Config{
		Sinker: &Sinker{
			BasePath: os.Getenv("SINKER_BASE_PATH"),
		},
		SinkerAPI: &SinkerAPI{
			StoreEventPath: os.Getenv("SINKER_API_STORE_EVENT_PATH"),
		},
	}
}
