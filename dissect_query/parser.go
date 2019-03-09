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
	"sort"
)

type parser struct {
	delimiters      []delimiter
	fields          []field
	referenceFields []field
}

var isIndirectField = func(field field) bool {
	switch field.(type) {
	case indirectField:
		return true
	default:
		return false
	}
}

func newParser(tokenizer string) (*parser, error) {
	matches := delimiterRE.FindAllStringSubmatchIndex(tokenizer, -1)
	if len(matches) == 0 {
		return nil, errInvalidTokenizer
	}

	var delimiters []delimiter
	var fields []field

	pos := 0
	for id, m := range matches {
		d := newDelimiter(tokenizer[m[2]:m[3]])
		key := tokenizer[m[4]:m[5]]
		field, err := newField(id, key, d)
		if err != nil {
			return nil, err
		}
		if field.IsGreedy() {
			d.MarkGreedy()
		}
		fields = append(fields, field)
		delimiters = append(delimiters, d)
		pos = m[5] + 1
	}

	if pos < len(tokenizer) {
		d := newDelimiter(tokenizer[pos:])
		delimiters = append(delimiters, d)
	}


	for i := 0; i < len(delimiters); i++ {
		if i+1 < len(delimiters) {
			delimiters[i].SetNext(delimiters[i+1])
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Ordinal() < fields[j].Ordinal()
	})

	var referenceFields []field
	for _, f := range fields {
		if !f.IsSaveable() {
			referenceFields = append(referenceFields, f)
		}
	}

	return &parser{
		delimiters:      delimiters,
		fields:          fields,
		referenceFields: referenceFields,
	}, nil
}

func filterFieldsWith(fields []field, predicate func(field) bool) []field {
	var filtered []field
	for _, field := range fields {
		if predicate(field) {
			filtered = append(filtered, field)
		}
	}
	return filtered
}
