
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
)

func TestConfig(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{value1}",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.NoError(t, err) {
			return
		}
	})

	t.Run("invalid", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%value1}",
			"field":     "message",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("with tokenizer missing", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("with empty tokenizer", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{
			"tokenizer": "",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})

	t.Run("tokenizer with no field defined", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{
			"tokenizer": "hello world",
		})
		if !assert.NoError(t, err) {
			return
		}

		cfg := config{}
		err = c.Unpack(&cfg)
		if !assert.Error(t, err) {
			return
		}
	})
}
