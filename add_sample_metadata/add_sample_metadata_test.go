package main

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)


func TestSampleDataRun(t *testing.T) {

	config := map[string]interface{}{
		"sample":                 0.00001,
		"log_type":               "warn",
		"query_type":             "finish",
	}

	testConfig, err := common.NewConfigFrom(config)
	assert.NoError(t, err)
	p, err := newSampleProcessor(testConfig)
	require.NoError(t, err)

	for i := 1; ; i++ {
		// 生成 event
		event := &beat.Event{
			Fields:    common.MapStr{},
			Timestamp: time.Now(),
		}
		_, _ = event.Fields.Put("message",
			"2019-02-11 00:03:52.464735 [warn] " +
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
		if i % 10 == 0 {
			assert.NotNil(t, newEvent, "")
			//fmt.Println(newEvent.Fields.GetValue("message"))
		} else {
			assert.Nil(t, newEvent)
		}
	}


}


