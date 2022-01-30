package lora

type Config struct {
	MAC string
	Keys
	Device
}

type Device struct {
	Addr string
}

type Keys struct {
	NetworkSKey     string
	ApplicationSKey string
}
