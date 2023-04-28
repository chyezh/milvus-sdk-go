package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := Config{}
	enable, err := c.parseEnableTLSAuth()
	assert.Nil(t, err)
	assert.False(t, enable)

	c = Config{
		EnableTLSAuth: true,
	}
	enable, err = c.parseEnableTLSAuth()
	assert.Nil(t, err)
	assert.True(t, enable)

	c = Config{
		Address: "https://xxxx-xxxx-xxxx",
	}
	enable, err = c.parseEnableTLSAuth()
	assert.Nil(t, err)
	assert.True(t, enable)
	addr, err := c.parseRemoteAddr()
	assert.Nil(t, err)
	assert.Equal(t, addr, "xxxx-xxxx-xxxx")

	c = Config{
		Address: "http://xxxx-xxxx-xxxx",
	}
	enable, err = c.parseEnableTLSAuth()
	assert.Nil(t, err)
	assert.False(t, enable)
	assert.Nil(t, err)
	assert.Equal(t, addr, "xxxx-xxxx-xxxx")
}
