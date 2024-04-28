package config

type Params struct {
	A               string `env:"SERVER_ADDRESS"`
	B               string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}
