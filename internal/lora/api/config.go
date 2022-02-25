package api

type Config struct {
	URL      string `koanf:"url"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`

	DeviceProfileID string `koanf:"device_profile_id"`
	ApplicationID   int64  `koanf:"application_id"`
}
