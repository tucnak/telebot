package layout

import (
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cast"
	tele "gopkg.in/tucnak/telebot.v3"
)

// Config represents typed map interface related to the "config" section in layout.
type Config struct {
	v map[string]interface{}
}

// Unmarshal parses the config into the out value. It's useful when you want to
// describe and to pre-define the fields in your custom configuration struct.
func (c *Config) Unmarshal(v interface{}) error {
	data, err := yaml.Marshal(c.v)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

// Get returns child map field wrapped into Config.
// If the field isn't map, returns nil.
func (c *Config) Get(k string) *Config {
	v, ok := c.v[k].(map[string]interface{})
	if !ok {
		return nil
	}
	return &Config{v: v}
}

// String returns a field casted to the string.
func (c *Config) String(k string) string {
	return cast.ToString(c.v[k])
}

// Int returns a field casted to the int.
func (c *Config) Int(k string) int {
	return cast.ToInt(c.v[k])
}

// Int64 returns a field casted to the int64.
func (c *Config) Int64(k string) int64 {
	return cast.ToInt64(c.v[k])
}

// Float returns a field casted to the float64.
func (c *Config) Float(k string) float64 {
	return cast.ToFloat64(c.v[k])
}

// Float returns a field casted to the bool.
func (c *Config) Bool(k string) bool {
	return cast.ToBool(c.v[k])
}

// Duration returns a field casted to the time.Duration.
// Accepts number-represented duration or a string in 0nsuµmh format.
func (c *Config) Duration(k string) time.Duration {
	return cast.ToDuration(c.v[k])
}

// ChatID returns a field casted to the ChatID.
// The value must be an integer.
func (c *Config) ChatID(k string) tele.ChatID {
	return tele.ChatID(c.Int64(k))
}
