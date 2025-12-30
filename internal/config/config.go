package config

import "os"

const (
	AWSRegionEUWest1 = "eu-west-1"

	WatcherIntervalSeconds = 5
)

type Config struct{
	AWS				*AWS
	Sinker 		*Sinker
	SinkerAPI *SinkerAPI
}

type AWS struct {
	Region		string
	S3Bucket	string
}

type Sinker struct{
	BasePath								string
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
		AWS: &AWS{
			Region:   AWSRegionEUWest1,
			S3Bucket: os.Getenv("AWS_BUCKET"),
		},
		Sinker: &Sinker{
			BasePath:								os.Getenv("SINKER_BASE_PATH"),
			WatcherIntervalSeconds:	WatcherIntervalSeconds,
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
