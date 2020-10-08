package layout

import (
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cast"
	tele "gopkg.in/tucnak/telebot.v3"
)

type Config struct {
	v map[string]interface{}
}

func (c *Config) Unmarshal(v interface{}) error {
	data, err := yaml.Marshal(c.v)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c.v)
}

func (c *Config) Get(k string) *Config {
	v, ok := c.v[k].(map[string]interface{})
	if !ok {
		return nil
	}
	return &Config{v: v}
}

func (c *Config) String(k string) string {
	return cast.ToString(c.v[k])
}

func (c *Config) Int(k string) int {
	return cast.ToInt(c.v[k])
}

func (c *Config) Int64(k string) int64 {
	return cast.ToInt64(c.v[k])
}

func (c *Config) Float(k string) float64 {
	return cast.ToFloat64(c.v[k])
}

func (c *Config) Bool(k string) bool {
	return cast.ToBool(c.v[k])
}

func (c *Config) Duration(k string) time.Duration {
	return cast.ToDuration(c.v[k])
}

func (c *Config) ChatID(k string) tele.ChatID {
	return tele.ChatID(c.Int64(k))
}
