package lora

type Config struct {
	MAC     string `koanf:"mac"`
	Keys    `koanf:"keys"`
	Devices []Device `koanf:"devices"`
}

type Device struct {
	Addr   string `koanf:"addr" yaml:"addr" json:"addr"`
	DevEUI string `koanf:"dev_eui" yaml:"dev_eui" json:"dev_eui"` // nolint: tagliatelle
}

type Keys struct {
	NetworkSKey     string `koanf:"network_skey"`
	ApplicationSKey string `koanf:"application_skey"`
}
