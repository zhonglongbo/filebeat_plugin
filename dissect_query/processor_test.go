// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)


func TestDissector_Dissect(t *testing.T) {

	config := map[string]interface{}{
		"tokenizer":   "%{logtime} [%{level}] [%{pid}:%{}] %{filename} zsearch query process %{query_state}[name:%{name}  %{timerstring}request:%{query}log_id:%{log_id} %{other}]",
		"field": "message",
		"target_prefix": "message_content",
	}
	testConfig, err := common.NewConfigFrom(config)
	assert.NoError(t, err)
	p, err := newProcessor(testConfig)
	require.NoError(t, err)

	for i:= 0; i<1;i++  {
		// 生成 event
		event := &beat.Event{
			Fields:    common.MapStr{},
			Timestamp: time.Now(),
		}
		_, _ = event.Fields.Put("message",
			"2019-02-11 00:03:52.464735 [info] " +
				"[144409:0x7f75bfbb6700] ckit_rpc_server.cc:376 " +
				"zsearch query process finish![name:select,  " +
				"cost_time:0, recv_time:0, wait_time:0, pt:0, cf:0, " +
				"qt:0, st:0, rt:0, ft:0, ot:0, ct:0, encodet:0, sendt:0," +
				"request:fl=tradeItemId,itemMarks&q=shopId:(1302438807) AND " +
				"categoryIds:(4121)&sort=sale desc&start=0&rows=50&filter=itemMarks:" +
				"{800}&debug=false&cache=true&appName=algo_prism, log_id:8airprhGUfT2u," +
				" c:true, ret:true, o:m=2538, r:]")

		newEvent, err := p.Run(event)
		assert.NoError(t, err)


		fmt.Println(newEvent)
		//fmt.Printf("%T\n",newEvent)
	}
}


func TestProcessor(t *testing.T) {
	tests := []struct {
		name   string
		c      map[string]interface{}
		fields common.MapStr
		values map[string]string
	}{
		{
			name:   "default field/default target",
			c:      map[string]interface{}{"tokenizer": "hello %{key}"},
			fields: common.MapStr{"message": "hello world"},
			values: map[string]string{"dissect.key": "world"},
		},
		{
			name:   "default field/target root",
			c:      map[string]interface{}{"tokenizer": "hello %{key}", "target_prefix": ""},
			fields: common.MapStr{"message": "hello world"},
			values: map[string]string{"key": "world"},
		},
		{
			name: "specific field/target root",
			c: map[string]interface{}{
				"tokenizer":     "hello %{key}",
				"target_prefix": "",
				"field":         "new_field",
			},
			fields: common.MapStr{"new_field": "hello world"},
			values: map[string]string{"key": "world"},
		},
		{
			name: "specific field/specific target",
			c: map[string]interface{}{
				"tokenizer":     "hello %{key}",
				"target_prefix": "new_target",
				"field":         "new_field",
			},
			fields: common.MapStr{"new_field": "hello world"},
			values: map[string]string{"new_target.key": "world"},
		},
		{
			name: "extract to already existing namespace not conflicting",
			c: map[string]interface{}{
				"tokenizer":     "hello %{key} %{key2}",
				"target_prefix": "extracted",
				"field":         "message",
			},
			fields: common.MapStr{"message": "hello world super", "extracted": common.MapStr{"not": "hello"}},
			values: map[string]string{"extracted.key": "world", "extracted.key2": "super", "extracted.not": "hello"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := common.NewConfigFrom(test.c)
			if !assert.NoError(t, err) {
				return
			}

			processor, err := newProcessor(c)
			if !assert.NoError(t, err) {
				return
			}

			e := beat.Event{Fields: test.fields}
			newEvent, err := processor.Run(&e)
			if !assert.NoError(t, err) {
				return
			}

			for field, value := range test.values {
				v, err := newEvent.GetValue(field)
				if !assert.NoError(t, err) {
					return
				}

				assert.Equal(t, value, v)
			}
		})
	}
}

func TestFieldDoesntExist(t *testing.T) {
	c, err := common.NewConfigFrom(map[string]interface{}{"tokenizer": "hello %{key}"})
	if !assert.NoError(t, err) {
		return
	}

	processor, err := newProcessor(c)
	if !assert.NoError(t, err) {
		return
	}

	e := beat.Event{Fields: common.MapStr{"hello": "world"}}
	_, err = processor.Run(&e)
	if !assert.Error(t, err) {
		return
	}
}

func TestFieldAlreadyExist(t *testing.T) {
	tests := []struct {
		name      string
		tokenizer string
		prefix    string
		fields    common.MapStr
	}{
		{
			name:      "no prefix",
			tokenizer: "hello %{key}",
			prefix:    "",
			fields:    common.MapStr{"message": "hello world", "key": "exists"},
		},
		{
			name:      "with prefix",
			tokenizer: "hello %{key}",
			prefix:    "extracted",
			fields:    common.MapStr{"message": "hello world", "extracted": "exists"},
		},
		{
			name:      "with conflicting key in prefix",
			tokenizer: "hello %{key}",
			prefix:    "extracted",
			fields:    common.MapStr{"message": "hello world", "extracted": common.MapStr{"key": "exists"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := common.NewConfigFrom(map[string]interface{}{
				"tokenizer":     test.tokenizer,
				"target_prefix": test.prefix,
			})

			if !assert.NoError(t, err) {
				return
			}

			processor, err := newProcessor(c)
			if !assert.NoError(t, err) {
				return
			}

			e := beat.Event{Fields: test.fields}
			_, err = processor.Run(&e)
			if !assert.Error(t, err) {
				return
			}
		})
	}
}

func TestErrorFlagging(t *testing.T) {
	t.Run("when the parsing fails add a flag", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{ok} - %{notvalid}",
		})

		if !assert.NoError(t, err) {
			return
		}

		processor, err := newProcessor(c)
		if !assert.NoError(t, err) {
			return
		}

		e := beat.Event{Fields: common.MapStr{"message": "hello world"}}
		event, err := processor.Run(&e)

		if !assert.Error(t, err) {
			return
		}

		flags, err := event.GetValue(beat.FlagField)
		if !assert.NoError(t, err) {
			return
		}

		assert.Contains(t, flags, flagParsingError)
	})

	t.Run("when the parsing is succesful do not add a flag", func(t *testing.T) {
		c, err := common.NewConfigFrom(map[string]interface{}{
			"tokenizer": "%{ok} %{valid}",
		})

		if !assert.NoError(t, err) {
			return
		}

		processor, err := newProcessor(c)
		if !assert.NoError(t, err) {
			return
		}

		e := beat.Event{Fields: common.MapStr{"message": "hello world"}}
		event, err := processor.Run(&e)

		if !assert.NoError(t, err) {
			return
		}

		_, err = event.GetValue(beat.FlagField)
		assert.Error(t, err)
	})
}
