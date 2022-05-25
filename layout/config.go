package layout

import (
	"strconv"
	"time"

	tele "github.com/TGeniusFamily/GOFSMtelebot"
	"github.com/spf13/cast"
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
// If the field isn't a map, returns nil.
func (c *Config) Get(k string) *Config {
	v, ok := c.v[k].(map[string]interface{})
	if !ok {
		return nil
	}
	return &Config{v: v}
}

// Index returns an i element from the array field, wrapped into Config.
// If the element isn't a map, returns nil. See also: Strings, Ints, Floats.
func (c *Config) Index(k string, i int) *Config {
	a, ok := c.v[k].([]interface{})
	if !ok {
		return nil
	}
	if i >= len(a) {
		return nil
	}
	v, ok := a[i].(map[string]interface{})
	if !ok {
		return nil
	}
	return &Config{v: v}
}

// Each iterates over the array field. Use it only with map elements.
func (c *Config) Each(k string, f func(int, *Config)) {
	a, ok := c.v[k].([]interface{})
	if !ok {
		return
	}

	for i, e := range a {
		v, ok := e.(map[string]interface{})
		if ok {
			f(i, &Config{v: v})
		}
	}
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

// Bool returns a field casted to the bool.
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

// Len returns the array field length.
func (c *Config) Len(k string) int {
	a := c.v[k].([]interface{})
	return len(a)
}

// Strings returns a field casted to the string slice.
func (c *Config) Strings(k string) []string {
	return cast.ToStringSlice(c.v[k])
}

// Ints returns a field casted to the int slice.
func (c *Config) Ints(k string) []int {
	return cast.ToIntSlice(c.v[k])
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
func (c *Config) Floats(k string) []float64 {
	slice := cast.ToSlice(c.v[k])

	fs := make([]float64, len(slice))
	for i, a := range slice {
		fs[i] = cast.ToFloat64(a)
	}

	return fs
}
