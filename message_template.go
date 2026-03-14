package astolform

import "strings"

type partType uint8

const (
	textPart partType = iota
	valuePart
	fieldPart
	paramsPart
)

type part struct {
	typ  partType
	text string
}

type Template struct {
	parts []part
}

func CompileTemplate(tpl string) *Template {

	parts := make([]part, 0, 8)
	start := 0

	for i := 0; i < len(tpl); i++ {

		if tpl[i] == '{' {

			if start < i {
				parts = append(parts, part{
					typ:  textPart,
					text: tpl[start:i],
				})
			}

			j := i + 1
			for j < len(tpl) && tpl[j] != '}' {
				j++
			}

			key := tpl[i+1 : j]

			switch key {

			case "field":
				parts = append(parts, part{typ: fieldPart})

			case "value":
				parts = append(parts, part{typ: valuePart})

			case "params":
				parts = append(parts, part{typ: paramsPart})

			default:
				parts = append(parts, part{
					typ:  textPart,
					text: tpl[i : j+1],
				})
			}

			i = j
			start = j + 1
		}
	}

	if start < len(tpl) {
		parts = append(parts, part{
			typ:  textPart,
			text: tpl[start:],
		})
	}

	return &Template{parts: parts}
}

func (t *Template) Render(field, value, param string) string {

	var b strings.Builder

	b.Grow(len(field) + len(value) + len(param) + 32)

	for _, p := range t.parts {

		switch p.typ {

		case textPart:
			b.WriteString(p.text)

		case fieldPart:
			b.WriteString(field)

		case valuePart:
			b.WriteString(value)

		case paramsPart:
			b.WriteString(param)

		}
	}

	return b.String()
}
