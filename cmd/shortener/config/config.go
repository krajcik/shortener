package config

import "flag"

func Create() *Params {
	p := new(Params)
	flag.StringVar(&p.A, "a", ":8080", "The address to listen on for HTTP requests.")
	flag.StringVar(&p.B, "b", "", "Host for shorten url")

	flag.Parse()

	return p
}
