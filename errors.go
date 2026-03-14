package astolform

import (
	"errors"
	"strings"
)

func GenerateErrorMessage(parts ...string) error {
	var b strings.Builder
	total := 0
	for _, s := range parts {
		total += len(s)
	}
	b.Grow(total)
	for _, s := range parts {
		b.WriteString(s)
	}
	return errors.New(b.String())
}
