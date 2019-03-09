

package main

import "fmt"

type Map = map[string]string

type positions []position

type position struct {
	start int
	end   int
}

type Dissector struct {
	raw    string
	parser *parser
}

func (d *Dissector) Dissect(s string) (Map, error) {
	if len(s) == 0 {
		return nil, errEmpty
	}

	positions, err := d.extract(s)
	if err != nil {
		return nil, err
	}

	if len(positions) == 0 {
		return nil, errParsingFailure
	}

	return d.resolve(s, positions), nil
}

func (d *Dissector) Raw() string {
	return d.raw
}

func (d *Dissector) extract(s string) (positions, error) {
	positions := make([]position, len(d.parser.fields))
	var i, start, lookahead, end int

	dl := d.parser.delimiters[0]
	offset := dl.IndexOf(s, 0)
	if offset == -1 || offset != 0 {
		return nil, fmt.Errorf(
			"could not find beginning delimiter: `%s` in remaining: `%s`, (offset: %d)",
			dl.Delimiter(), s, 0,
		)
	}
	offset += dl.Len()

	for dl.Next() != nil {
		start = offset
		end = dl.Next().IndexOf(s, offset)
		if end == -1 {
			return nil, fmt.Errorf(
				"could not find delimiter: `%s` in remaining: `%s`, (offset: %d)",
				dl.Delimiter(), s[offset:], offset,
			)
		}

		offset = end

		if dl.IsGreedy() {
			for {
				lookahead = dl.Next().IndexOf(s, offset+1)
				if lookahead != offset+1 {
					break
				} else {
					offset = lookahead
				}
			}
		}

		positions[i] = position{start: start, end: end}
		offset += dl.Next().Len()
		i++
		dl = dl.Next()
	}

	if offset < len(s) && i < len(d.parser.fields) {
		positions[i] = position{start: offset, end: len(s)}
	}
	return positions, nil
}

func (d *Dissector) resolve(s string, p positions) Map {
	m := make(Map, len(p))
	for _, f := range d.parser.fields {
		pos := p[f.ID()]
		f.Apply(s[pos.start:pos.end], m)
	}

	for _, f := range d.parser.referenceFields {
		delete(m, f.Key())
	}
	return m
}

func New(tokenizer string) (*Dissector, error) {
	p, err := newParser(tokenizer)
	if err != nil {
		return nil, err
	}

	if err := validate(p); err != nil {
		return nil, err
	}

	return &Dissector{parser: p, raw: tokenizer}, nil
}
