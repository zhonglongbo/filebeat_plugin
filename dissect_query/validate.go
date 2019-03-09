

package main

import (
	"fmt"
)

func validate(p *parser) error {
	indirectFields := filterFieldsWith(p.fields, isIndirectField)

	for _, field := range indirectFields {
		found := false
		for _, reference := range p.referenceFields {
			if reference.Key() == field.Key() {
				found = true
				break
			}
		}

		if found == false {
			return fmt.Errorf("missing reference for key '%s'", field.Key())
		}
	}

	return nil
}
