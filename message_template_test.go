package astolform

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCompileTemplate(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  []part
	}{
		{
			name:  "text then placeholder",
			input: "foo {value}",
			want: []part{
				{typ: textPart, text: "foo "},
				{typ: valuePart},
			},
		},
		{
			name:  "placeholder then text",
			input: "{field} bar",
			want: []part{
				{typ: fieldPart},
				{typ: textPart, text: " bar"},
			},
		},
		{
			name:  "multiple placeholders",
			input: "foo {value} bar {field}",
			want: []part{
				{typ: textPart, text: "foo "},
				{typ: valuePart},
				{typ: textPart, text: " bar "},
				{typ: fieldPart},
			},
		},
		{
			name:  "only text",
			input: "hello world",
			want: []part{
				{typ: textPart, text: "hello world"},
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			result := CompileTemplate(tt.input)

			if !reflect.DeepEqual(result.parts, tt.want) {
				t.Errorf(
					"CompileTemplate(%q)\nGot: %#v\nWant: %#v",
					tt.input,
					result.parts,
					tt.want,
				)
			}

		})

	}
}

func TestRenderTemplate(t *testing.T) {

	tpl := CompileTemplate(
		"{field} must be greater than {params}, got {value}",
	)

	fmt.Println(tpl)

	got := tpl.Render("age", "5", "10")

	want := "age must be greater than 10, got 5"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
