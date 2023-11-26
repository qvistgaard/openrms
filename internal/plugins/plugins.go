package plugins

import "sort"

type Plugins struct {
	car    []Car
	race   []Race
	config map[string]pluginConfig
}

func New(config *Config) (*Plugins, error) {
	return &Plugins{
		car:    make([]Car, 0),
		race:   make([]Race, 0),
		config: config.Plugins,
	}, nil
}

func (p *Plugins) Car() []Car {
	return p.car
}

func (p *Plugins) Race() []Race {
	return p.race
}

func (p *Plugins) Enabled(plugin Plugin) bool {
	if plug, ok := p.config[plugin.Name()]; ok {
		return plug.Enabled
	}
	return false
}

func (p *Plugins) Append(plugin Plugin) Plugin {
	if c, ok := p.config[plugin.Name()]; ok && c.Enabled {
		if rule, ok := plugin.(Car); ok {
			p.car = append(p.car, rule)
		}
		if rule, ok := plugin.(Race); ok {
			p.race = append(p.race, rule)
		}
		sort.Slice(p.car, func(i, j int) bool {
			return p.car[i].(Plugin).Priority() < p.car[j].(Plugin).Priority()
		})
		sort.Slice(p.race, func(i, j int) bool {
			return p.race[i].(Plugin).Priority() < p.race[j].(Plugin).Priority()
		})
	}
	return plugin
}
