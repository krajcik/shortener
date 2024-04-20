package config

import "flag"
import "github.com/caarlos0/env/v11"

func Create() *Params {
	p := new(Params)
	err := env.Parse(p)
	if err != nil {
		panic(err)
	}
	if p.A == "" {
		flag.StringVar(&p.A, "a", ":8080", "The address to listen on for HTTP requests.")
	}
	if p.B == "" {
		flag.StringVar(&p.B, "b", "", "Host for shorten url")
	}

	flag.Parse()

	return p
}
