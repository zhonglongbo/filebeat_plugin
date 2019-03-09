
package main

import (
	"fmt"
	"github.com/elastic/beats/libbeat/logp"
	"strings"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

const flagParsingError = "dissect_parsing_error"

type processor struct {
	config config
}

const (
	processorName = "dissect_query"
)

func newProcessor(c *common.Config) (processors.Processor, error) {
	config := defaultConfig
	err := c.Unpack(&config)
	if err != nil {
		return nil, err
	}
	p := &processor{config: config}

	return p, nil
}

func ParseQuery(event *beat.Event) (*beat.Event, error) {

	query,err := event.Fields.GetValue("message_content.query")

	//fmt.Println("message:",message_val)
	if err != nil {
		logp.Info("Get message_content.query fail reason of [%s]",err)
		return event,err
	}
	var m Map
	m = map[string]string{}

	//解析query
	spiltRet := strings.Split(query.(string),"&")
	for k,_ := range spiltRet {
		var ret = strings.Split(spiltRet[k],"=")
		if len(ret)==2 {
			m[ret[0]] = ret[1]
		}
	}
	err = event.Fields.Delete("message_content.query")
	if err!=nil {
		logp.Info("Delete message_content.query fail reason of [%s]",err)

	}
	event,err = mapperToContent(event,mapToMapStr(m))
	//fmt.Println("mapper_to_content",event)
	return event,nil
}

func (p *processor) Run(event *beat.Event) (*beat.Event, error) {
	v, err := event.GetValue(p.config.Field)
	if err != nil {
		return event, err
	}

	s, ok := v.(string)
	if !ok {
		return event, fmt.Errorf("field is not a string, value: `%v`, field: `%s`", v, p.config.Field)
	}

	m, err := p.config.Tokenizer.Dissect(s)

	if err != nil {
		if err := common.AddTagsWithKey(
			event.Fields,
			"log.flags",
			[]string{flagParsingError},
		); err != nil {
			return event, errors.Wrap(err, "cannot add new flag the event")
		}
		return event, err
	}

	event, err = p.mapper(event, mapToMapStr(m))

	event,err = ParseQuery(event)
	if err != nil {
		logp.Info("ParseQuery fail reason of [%s]",err)
		return event, err
	}

	return event, nil
}

func mapperToContent(event *beat.Event, m common.MapStr) (*beat.Event, error) {
	copy := event.Fields.Clone()
	prefix := ""
	prefix = "message_content" + "."
	var prefixKey string
	for k, v := range m {
		prefixKey = prefix + k
		if _, err := event.GetValue(prefixKey); err == common.ErrKeyNotFound {
			event.PutValue(prefixKey, v)
		} else {
			event.Fields = copy
			if err != nil {
				return event, errors.Wrapf(err, "cannot override existing key with `%s`", prefixKey)
			}
			return event, fmt.Errorf("cannot override existing key with `%s`", prefixKey)
		}
	}
	return event, nil
}

func (p *processor) mapper(event *beat.Event, m common.MapStr) (*beat.Event, error) {
	copy := event.Fields.Clone()

	prefix := ""
	if p.config.TargetPrefix != "" {
		prefix = p.config.TargetPrefix + "."
	}
	var prefixKey string
	for k, v := range m {
		prefixKey = prefix + k
		if _, err := event.GetValue(prefixKey); err == common.ErrKeyNotFound {
			event.PutValue(prefixKey, v)
		} else {
			event.Fields = copy
			if err != nil {
				return event, errors.Wrapf(err, "cannot override existing key with `%s`", prefixKey)
			}
			return event, fmt.Errorf("cannot override existing key with `%s`", prefixKey)
		}
	}

	return event, nil
}

func (p *processor) String() string {
	return "dissect=" + p.config.Tokenizer.Raw() +
		",field=" + p.config.Field +
		",target_prefix=" + p.config.TargetPrefix
}

func mapToMapStr(m Map) common.MapStr {
	newMap := make(common.MapStr, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
