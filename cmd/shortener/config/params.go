package config

type Params struct {
	A string `env:"SERVER_ADDRESS"`
	B string `env:"BASE_URL"`
}
