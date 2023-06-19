package app

import (
	"os"
	"strings"
)

type envSource struct {
	keyCase KeyCase
	configs map[string]Config
}

func NewEnvSource(keyCase KeyCase) ConfigSource {
	return &envSource{
		keyCase: keyCase,
		configs: make(map[string]Config),
	}
}

func (e *envSource) Load() error {

	envs := os.Environ()

	for _, env := range envs {

		sep := strings.Index(env, "=")

		key := ConvertKeyCase(env[:sep], e.keyCase)
		val := env[sep+1:]

		e.configs[key] = NewConfig(val)

	}

	return nil

}

func (e *envSource) Get(k string) Config {

	if c, ok := e.configs[k]; ok {
		return c
	}

	return EmptyConfig()

}

func (e *envSource) Has(k string) bool {

	_, ok := e.configs[k]

	return ok

}
