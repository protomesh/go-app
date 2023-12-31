package app

import (
	"flag"
)

type FlagSet interface {
	Visit(fn func(*flag.Flag))
	VisitAll(fn func(*flag.Flag))
}

type flagSource struct {
	keyCase KeyCase
	flagSet FlagSet
	configs map[string]Config
	onlySet map[string]bool
}

func NewFlagSource(keyCase KeyCase, flagSet FlagSet) ConfigSource {
	return &flagSource{
		keyCase: keyCase,
		flagSet: flagSet,
		configs: make(map[string]Config),
		onlySet: make(map[string]bool),
	}
}

func (f *flagSource) Load() error {

	f.flagSet.Visit(func(fg *flag.Flag) {

		key := ConvertKeyCase(fg.Name, f.keyCase)

		f.onlySet[key] = true

	})

	f.flagSet.VisitAll(func(fg *flag.Flag) {

		key := ConvertKeyCase(fg.Name, f.keyCase)
		val := fg.Value.String()

		if len(val) == 0 {
			return
		}

		f.configs[key] = NewConfig(val)

	})

	return nil
}

func (f *flagSource) Get(k string) Config {

	if c, ok := f.configs[k]; ok {
		return c
	}

	return EmptyConfig()

}

func (f *flagSource) Has(k string) bool {

	_, ok := f.onlySet[k]

	return ok

}
