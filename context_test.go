package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Context = (*nativeContext)(nil)

func TestContext(t *testing.T) {
	t.Run("Get,Set", func(t *testing.T) {
		var c Context
		c = new(nativeContext)
		c.Set("name", "Jon Snow")
		assert.Equal(t, "Jon Snow", c.Get("name"))
	})
}
