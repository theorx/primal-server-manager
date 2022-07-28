package Config

import "errors"

type Config struct {
	params map[string]Param
}

func (c *Config) Add(p Param) error {
	if _, ok := c.params[p.GetKey()]; ok {
		return errors.New("config parameter with given key already exists")
	}

	c.params[p.GetKey()] = p

	return nil
}

func (c *Config) Get(key string) (Param, error) {

	if _, ok := c.params[key]; ok {
		return c.params[key], nil
	}

	return Param{}, errors.New("")
}

func (c *Config) List() map[string]string {

	m := make(map[string]string)
	for _, p := range c.params {
		m[p.GetKey()] = p.GetValue()
	}

	return m
}

func (c *Config) ListTyped(t ParamType) map[string]string {

	m := make(map[string]string)
	for _, p := range c.params {
		if p.GetType() == t {
			m[p.GetKey()] = p.GetValue()
		}
	}

	return m
}
