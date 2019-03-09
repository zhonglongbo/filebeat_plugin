package main

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"
	"github.com/pkg/errors"
	"strings"
	"sync/atomic"
)

type SampleData struct {
	config Config
}

var (
	count int64 = 0
)

const (
	processorName = "sample_data"
)

func newSampleProcessor(cfg *common.Config) (processors.Processor, error) {
	config := defaultConfig()

	if err := cfg.Unpack(&config); err != nil {
		return nil, errors.Wrapf(err, "fail to unpack the %v configuration", processorName)
	}

	sampledata := &SampleData{
		config: config,
	}
	return sampledata, nil
}

func (sampledata SampleData) Run(event *beat.Event) (*beat.Event, error) {
	//根据采样配置进行采样
	if strings.Contains(event.Fields.String(), sampledata.config.QueryType) &&
		strings.Contains(event.Fields.String(), sampledata.config.LogType) &&
		strings.Contains(event.Fields.String(), "request:") {
		atomic.AddInt64(&count, 1)
		if int64(sampledata.config.Sample*float64(count)) >= 1 {
			logp.Debug(processorName, "sample successful :%s", event)
			count = 0
			return event, nil
		} else {
			logp.Debug(processorName, "sample data abnormal :%s", event)
			return nil, nil
		}

	} else {
		return event, nil
	}
}

func (sampledata SampleData) String() string {
	return "sampledata"
}
