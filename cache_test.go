package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache_Get(t *testing.T) {
	mc := NewInMemoryCache()

	test, _ := mc.Get(CacheUserContext, "test")
	assert.Equal(t, nil, test)
}

func TestInMemoryCache_Set(t *testing.T) {
	mc := NewInMemoryCache()

	_ = mc.Put(CacheUserContext, "test", "test")
	test, _ := mc.Get(CacheUserContext, "test")
	assert.Equal(t, "test", test)
}

func TestInMemoryCache_Keys(t *testing.T) {
	mc := NewInMemoryCache()

	_ = mc.Put(CacheUserContext, "test", "test")

	keys := mc.Keys(CacheUserContext)
	assert.Equal(t, []string{"test"}, keys)
}
