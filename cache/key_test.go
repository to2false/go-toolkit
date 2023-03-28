package cache

import (
	"testing"

	"gotest.tools/assert"
)

func Test_CacheKey(t *testing.T) {
	var c Key = "test"
	assert.Equal(t, c.String(), "test")

	SetCacheNamespace("ns")
	assert.Equal(t, c.String(), "ns:test")

	var f Key = "test:%s"
	assert.Equal(t, f.Format("v").String(), "ns:test:v")
}
