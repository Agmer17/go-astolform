package astolform

import (
	"fmt"
	"reflect"
	"testing"
)

type Account struct {
	Name     string `should:"required"`
	Email    string `should:"required"`
	Password string `should:"required"`
}

var globalOk bool
var globalErrs map[string]string

func BenchmarkValidator(b *testing.B) {
	v := NewDefault()
	user := Account{
		Name:     "Rizky",
		Email:    "rizky@mail.com",
		Password: "secret",
	}

	b.ReportAllocs()
	for b.Loop() {
		globalOk, globalErrs = v.Validate(user) // sink ke package-level var
	}
	_ = globalOk
	_ = globalErrs
}

func BenchmarkValidatorError(b *testing.B) {

	v := NewDefault()

	user := Account{
		Name:     "",
		Email:    "",
		Password: "",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		globalOk, globalErrs = v.Validate(user) // sink ke package-level var
	}
	_ = globalOk
	_ = globalErrs
}

func BenchmarkParseStruct(b *testing.B) {

	t := reflect.TypeFor[Account]()

	for b.Loop() {
		parseStruct(t)
	}
}

func BenchmarkValidatorCold(b *testing.B) {
	user := Account{
		Name:     "Rizky",
		Email:    "rizky@mail.com",
		Password: "secret",
	}

	for b.Loop() {
		v := NewDefault()
		globalOk, globalErrs = v.Validate(user)
	}
}

func TestValueToString(t *testing.T) {

	vInt := reflect.ValueOf(42)
	vString := reflect.ValueOf("foo")
	vBool := reflect.ValueOf(true)
	vUint := reflect.ValueOf(uint8(10))
	vFloat := reflect.ValueOf(float64(1.00000001))

	i := valueToString(vInt)
	s := valueToString(vString)
	b := valueToString(vBool)
	u := valueToString(vUint)
	f := valueToString(vFloat)

	want := "42 foo true 10 1.00000001"
	got := fmt.Sprintf("%s %s %s %s %s", i, s, b, u, f)

	if want != got {
		t.Errorf("TestValueToString gagal!\nWant: %q\nGot : %q", want, got)
	}

}
