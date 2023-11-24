package plugins

type Config struct {
	Plugins map[string]pluginConfig
}

type pluginConfig struct {
	Enabled bool
}
