package teletest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	what interface{}
	opts []interface{}
}

func (r Response) Expect(t *testing.T, what interface{}, opts []interface{}) {
	assert.Equal(t, what, r.what)
	assert.Equal(t, opts, r.opts)
}
