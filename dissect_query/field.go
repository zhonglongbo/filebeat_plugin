

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type field interface {
	MarkGreedy()
	IsGreedy() bool
	Ordinal() int
	Key() string
	ID() int
	Apply(b string, m Map)
	String() string
	IsSaveable() bool
}

type baseField struct {
	id      int
	key     string
	ordinal int
	greedy  bool
}

func (f baseField) IsGreedy() bool {
	return f.greedy
}

func (f baseField) MarkGreedy() {
	f.greedy = true
}

func (f baseField) Ordinal() int {
	return f.ordinal
}

func (f baseField) Key() string {
	return f.key
}

func (f baseField) ID() int {
	return f.id
}

func (f baseField) IsSaveable() bool {
	return true
}

func (f baseField) String() string {
	return fmt.Sprintf("field: %s, ordinal: %d, greedy: %v", f.key, f.ordinal, f.IsGreedy())
}

type normalField struct {
	baseField
}

func (f normalField) Apply(b string, m Map) {
	m[f.Key()] = b
}

type skipField struct {
	baseField
}

func (f skipField) Apply(b string, m Map) {
}

func (f skipField) IsSaveable() bool {
	return false
}

type namedSkipField struct {
	baseField
}

func (f namedSkipField) Apply(b string, m Map) {
	m[f.Key()] = b
}

func (f namedSkipField) IsSaveable() bool {
	return false
}

type pointerField struct {
	baseField
}

func (f pointerField) Apply(b string, m Map) {
	m[f.Key()] = b
}

func (f pointerField) IsSaveable() bool {
	return false
}

type indirectField struct {
	baseField
}

func (f indirectField) Apply(b string, m Map) {
	v, ok := m[f.Key()]
	if ok {
		m[v] = b
		return
	}
}

// dissect: %{+key} %{+key}
// message: hello world
// result:
//	key: hello world
//
// dissect: %{+key/2} %{+key/1}
// message: hello world
// result:
//	key: world hello
type appendField struct {
	baseField
	previous delimiter
}

func (f appendField) Apply(b string, m Map) {
	v, ok := m[f.Key()]
	if ok {
		m[f.Key()] = v + f.JoinString() + b
		return
	}
	m[f.Key()] = b
}

func (f appendField) JoinString() string {
	if f.previous == nil || f.previous.Len() == 0 {
		return defaultJoinString
	}
	return f.previous.Delimiter()
}

func newField(id int, rawKey string, previous delimiter) (field, error) {
	if len(rawKey) == 0 {
		return newSkipField(id), nil
	}

	key, ordinal, greedy := extractKeyParts(rawKey)

	if strings.HasPrefix(key, appendIndirectPrefix) {
		return nil, errMixedPrefixIndirectAppend
	}

	if strings.HasPrefix(key, indirectAppendPrefix) {
		return nil, errMixedPrefixAppendIndirect
	}

	if strings.HasPrefix(key, skipFieldPrefix) {
		return newNamedSkipField(id, key[1:]), nil
	}

	if strings.HasPrefix(key, pointerFieldPrefix) {
		return newPointerField(id, key[1:]), nil
	}

	if strings.HasPrefix(key, appendFieldPrefix) {
		return newAppendField(id, key[1:], ordinal, greedy, previous), nil
	}

	if strings.HasPrefix(key, indirectFieldPrefix) {
		return newIndirectField(id, key[1:]), nil
	}

	return newNormalField(id, key, ordinal, greedy), nil
}

func newSkipField(id int) skipField {
	return skipField{baseField{id: id}}
}

func newNamedSkipField(id int, key string) namedSkipField {
	return namedSkipField{
		baseField{id: id, key: key},
	}
}

func newPointerField(id int, key string) pointerField {
	return pointerField{
		baseField{id: id, key: key},
	}
}

func newAppendField(id int, key string, ordinal int, greedy bool, previous delimiter) appendField {
	return appendField{
		baseField: baseField{
			id:      id,
			key:     key,
			ordinal: ordinal,
			greedy:  greedy,
		},
		previous: previous,
	}
}

func newIndirectField(id int, key string) indirectField {
	return indirectField{
		baseField{
			id:  id,
			key: key,
		},
	}
}

func newNormalField(id int, key string, ordinal int, greedy bool) normalField {
	return normalField{
		baseField{
			id:      id,
			key:     key,
			ordinal: ordinal,
			greedy:  greedy,
		},
	}
}

func extractKeyParts(rawKey string) (key string, ordinal int, greedy bool) {
	m := suffixRE.FindAllStringSubmatch(rawKey, -1)

	if m[0][3] != "" {
		ordinal, _ = strconv.Atoi(m[0][3])
	}

	if strings.EqualFold(greedySuffix, m[0][4]) {
		greedy = true
	}
	return m[0][1], ordinal, greedy
}
