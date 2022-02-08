package lora

type Config struct {
	MAC    string `koanf:"mac"`
	Keys   `koanf:"keys"`
	Device `koanf:"device"`
}

type Device struct {
	Addr string `koanf:"addr"`
}

type Keys struct {
	NetworkSKey     string `koanf:"network_skey"`
	ApplicationSKey string `koanf:"application_skey"`
}
