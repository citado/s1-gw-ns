package lora

type Config struct {
	MAC     string `json:"mac"`
	Keys    `json:"keys"`
	Devices []Device `json:"devices"`
}

type Device struct {
	Addr   string `json:"addr"`
	DevEUI string `json:"dev_eui"` // nolint: tagliatelle
}

type Keys struct {
	NetworkSKey     string `json:"network_skey"`     // nolint: tagliatelle
	ApplicationSKey string `json:"application_skey"` // nolint: tagliatelle
}
