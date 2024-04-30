package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

func Create() (*Params, error) {
	p := &Params{}
	if err := initEnv(p); err != nil {
		return nil, err
	}
	initFlags(p)

	return p, nil
}

func initEnv(p *Params) error {
	err := env.Parse(p)
	if err != nil {
		return err
	}
	return nil
}

func initFlags(p *Params) {
	if p.A == "" {
		flag.StringVar(&p.A, "a", ":8080", "The address to listen on for HTTP requests.")
	}
	if p.B == "" {
		flag.StringVar(&p.B, "b", "", "Host for shorten url")
	}
	if p.FileStoragePath == "" {
		flag.StringVar(&p.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	}

	if p.DatabaseDsn == "" {
		flag.StringVar(&p.DatabaseDsn, "d", "postgres://postgres:pgpass@localhost:5432/shortener", "Database DSN")
	}

	flag.Parse()
}
