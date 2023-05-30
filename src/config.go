package main

type Config struct {
	Path      string          `toml:"path"`
	Addresses AddressesCongif `toml:"addresses"`
}

type AddressesCongif struct {
	Server_addr string `toml:"server_addr"`
	Frontend    string `toml:"frontend"`
	Hil         string `toml:"hil"`
}
