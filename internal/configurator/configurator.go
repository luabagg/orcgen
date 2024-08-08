package configurator

import "github.com/luabagg/orcgen/internal/generator"

// Configurator interface allows for a generator to be configured
type Configurator interface {
	// Configure sets up a generator to start working
	Configure(c generator.Config) generator.Generator
}
