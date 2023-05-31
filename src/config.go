package main

type Config struct {
	Path      string          `toml:"path"`
	Addresses AddressesConfig `toml:"addresses"`
}

type AddressesConfig struct {
	Server_addr string `toml:"server_addr"`
	Frontend    string `toml:"frontend"`
	Hil         string `toml:"hil"`
}
