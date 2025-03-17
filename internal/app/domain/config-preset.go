package domain

type ConfigPreset struct {
	Key      string
	Services []string
}

func NewConfigPreset(
	Key string,
	Services []string,
) ConfigPreset {
	return ConfigPreset{
		Key:      Key,
		Services: Services,
	}
}
