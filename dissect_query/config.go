
package main

type config struct {
	Tokenizer    *tokenizer `config:"tokenizer" validate:"required"`
	Field        string     `config:"field"`
	TargetPrefix string     `config:"target_prefix"`
}

var defaultConfig = config{
	Field:        "message",
	TargetPrefix: "dissect",
}

type tokenizer = Dissector

func (t *tokenizer) Unpack(v string) error {
	d, err := New(v)
	if err != nil {
		return err
	}
	*t = *d
	return nil
}
