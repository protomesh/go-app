package app

type compositeSource struct {
	s   []ConfigSource
	crs map[string]Config
}

func NewCompositeSource(s ...ConfigSource) ConfigSource {
	return &compositeSource{
		s:   s,
		crs: make(map[string]Config),
	}
}

func (c *compositeSource) Load() error {

	for _, cs := range c.s {

		err := cs.Load()
		if err != nil {
			return err
		}

	}

	return nil

}

func (c *compositeSource) Get(k string) Config {

	if cr, ok := c.crs[k]; ok {
		return cr
	}

	var defVal Config

	for _, s := range c.s {

		cr := s.Get(k)

		if cr != nil && cr.IsSet() {

			if s.Has(k) {
				c.crs[k] = cr
				return cr
			}

			defVal = cr

		}

	}

	if defVal != nil {
		return defVal
	}

	return EmptyConfig()
}

func (c *compositeSource) Has(k string) bool {

	cfg := c.Get(k)

	return cfg != nil && cfg.IsSet()

}
