package layout

import (
	"strconv"
	"time"

	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
)

// Config represents typed map interface related to the "config" section in layout.
type Config struct {
	v *viper.Viper
}

// Unmarshal parses the whole config into the out value. It's useful when you want to
// describe and to pre-define the fields in your custom configuration struct.
func (c *Config) Unmarshal(v interface{}) error {
	return c.v.Unmarshal(v)
}

// UnmarshalKey parses the specific key in the config into the out value.
func (c *Config) UnmarshalKey(k string, v interface{}) error {
	return c.v.UnmarshalKey(k, v)
}

// Get returns a child map field wrapped into Config.
// If the field isn't a map, returns nil.
func (c *Config) Get(k string) *Config {
	v := c.v.Sub(k)
	if v == nil {
		return nil
	}
	return &Config{v: v}
}

// Slice returns a child slice of objects wrapped into Config.
// If the field isn't a slice, returns nil.
func (c *Config) Slice(k string) (slice []*Config) {
	a, ok := c.v.Get(k).([]interface{})
	if !ok {
		return nil
	}

	for i := range a {
		m, ok := a[i].(map[string]interface{})
		if !ok {
			return nil
		}

		v := viper.New()
		v.MergeConfigMap(m)
		slice = append(slice, &Config{v: v})
	}

	return
}

// String returns a field casted to the string.
func (c *Config) String(k string) string {
	return c.v.GetString(k)
}

// Int returns a field casted to the int.
func (c *Config) Int(k string) int {
	return c.v.GetInt(k)
}

// Int64 returns a field casted to the int64.
func (c *Config) Int64(k string) int64 {
	return c.v.GetInt64(k)
}

// Float returns a field casted to the float64.
func (c *Config) Float(k string) float64 {
	return c.v.GetFloat64(k)
}

// Bool returns a field casted to the bool.
func (c *Config) Bool(k string) bool {
	return c.v.GetBool(k)
}

// Duration returns a field casted to the time.Duration.
// Accepts number-represented duration or a string in 0nsuµmh format.
func (c *Config) Duration(k string) time.Duration {
	return c.v.GetDuration(k)
}

// ChatID returns a field casted to the ChatID.
// The value must be an integer.
func (c *Config) ChatID(k string) tele.ChatID {
	return tele.ChatID(c.Int64(k))
}

// Strings returns a field casted to the string slice.
func (c *Config) Strings(k string) []string {
	return c.v.GetStringSlice(k)
}

// Ints returns a field casted to the int slice.
func (c *Config) Ints(k string) []int {
	return c.v.GetIntSlice(k)
}

// Int64s returns a field casted to the int64 slice.
func (c *Config) Int64s(k string) (ints []int64) {
	for _, s := range c.Strings(k) {
		i, _ := strconv.ParseInt(s, 10, 64)
		ints = append(ints, i)
	}
	return ints
}

// Floats returns a field casted to the float slice.
func (c *Config) Floats(k string) (floats []float64) {
	for _, s := range c.Strings(k) {
		i, _ := strconv.ParseFloat(s, 64)
		floats = append(floats, i)
	}
	return floats
}
