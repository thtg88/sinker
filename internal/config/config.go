package config

import "os"

const WatcherIntervalSeconds = 5

type Config struct{
	Sinker 		*Sinker
	SinkerAPI *SinkerAPI
}

type Sinker struct{
	BasePath								string
	S3BucketName						string
	WatcherIntervalSeconds	int64
}

type SinkerAPI struct {
	APIKey					string
	BaseURL					string
	HeaderNames			*SinkerAPIHeaderNames
	StoreDevicePath string
	StoreEventPath	string
	UserID					string
}

type SinkerAPIHeaderNames struct {
	APIKey		string
	DeviceID	string
	UserID		string
}

func Load() *Config {
	return &Config{
		Sinker: &Sinker{
			BasePath:			os.Getenv("SINKER_BASE_PATH"),
			S3BucketName:	os.Getenv("AWS_BUCKET"),
			WatcherIntervalSeconds: WatcherIntervalSeconds,
		},
		SinkerAPI: &SinkerAPI{
			APIKey:						os.Getenv("SINKER_API_KEY_HEADER_VALUE"),
			BaseURL:					os.Getenv("SINKER_API_BASE_URL"),
			StoreDevicePath:	os.Getenv("SINKER_API_STORE_DEVICE_PATH"),
			StoreEventPath:		os.Getenv("SINKER_API_STORE_EVENT_PATH"),
			UserID:						os.Getenv("SINKER_API_USER_ID_HEADER_VALUE"),
			HeaderNames: &SinkerAPIHeaderNames{
				APIKey:   os.Getenv("SINKER_API_KEY_HEADER_NAME"),
				DeviceID: os.Getenv("SINKER_API_DEVICE_ID_HEADER_NAME"),
				UserID:   os.Getenv("SINKER_API_USER_ID_HEADER_NAME"),
			},
		},
	}
}
