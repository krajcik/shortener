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
	if p.A == "" && flag.Lookup("a") != nil {
		flag.StringVar(&p.A, "a", ":8080", "The address to listen on for HTTP requests.")
	}
	if p.B == "" && flag.Lookup("b") != nil {
		flag.StringVar(&p.B, "b", "", "Host for shorten url")
	}
	if p.FileStoragePath == "" && flag.Lookup("f") != nil {
		flag.StringVar(&p.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	}
	flag.Parse()
}
