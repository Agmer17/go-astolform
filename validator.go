package astolform

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const tagKey = "should"

type Validator struct {
	refCache      sync.Map
	customRules   map[string]ValidatorFunc
	customMessage map[string]*Template
}

func NewDefault() *Validator {
	return &Validator{
		customRules:   make(map[string]ValidatorFunc),
		customMessage: make(map[string]*Template),
	}
}

func (v *Validator) Validate(s any) (bool, map[string]string) {

	mtd := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	if mtd.Kind() == reflect.Pointer {
		mtd = mtd.Elem()
		val = val.Elem()
	}

	meta := v.getCachedStruct(mtd)

	var errsMap map[string]string

	success := true

	for _, field := range meta.fields {

		value := val.Field(field.index)

		for _, rule := range field.rules {

			if err := rule.fn(value, rule.param, field.fieldName); err != nil {
				if errsMap == nil {
					errsMap = make(map[string]string, 1)
				}
				success = false

				if t, ok := v.customMessage[rule.ruleName]; ok {
					errsMap[field.fieldName] = t.Render(field.fieldName, valueToString(value), rule.param)
					continue
				}
				errsMap[field.fieldName] = err.Error()
			}
		}
	}

	return success, errsMap
}

// parsing struct to be cache and use later!
func parseStruct(sType reflect.Type) *structCache {
	nField := sType.NumField()
	meta := &structCache{
		fields: make([]fieldMetaData, 0, nField),
	}

	for i := range nField {
		field := sType.Field(i)
		tags := field.Tag.Get(tagKey)
		if tags == "" {
			continue
		}

		fieldMeta := fieldMetaData{
			index:     i,
			rules:     make([]compiledRule, 0, 4), // default is 4. you can change this if you want
			fieldName: field.Name,
		}

		remaining := tags
		for {
			var token string
			idx := strings.IndexByte(remaining, ',')
			if idx == -1 {
				token = remaining
				remaining = ""
			} else {
				token = remaining[:idx]
				remaining = remaining[idx+1:]
			}

			ruleName, params, found := strings.Cut(token, "=")
			if !found {
				params = ""
			}

			if f := basicRules[ruleName]; f != nil {
				fieldMeta.rules = append(fieldMeta.rules, compiledRule{
					fn:       f,
					param:    params,
					ruleName: ruleName,
				})
			}

			if remaining == "" {
				break
			}
		}

		meta.fields = append(meta.fields, fieldMeta)
	}

	return meta
}

func (v *Validator) getCachedStruct(t reflect.Type) *structCache {

	if meta, ok := v.refCache.Load(t); ok {
		return meta.(*structCache)
	}

	meta := parseStruct(t)

	act, _ := v.refCache.LoadOrStore(t, meta)

	return act.(*structCache)
}

func (v *Validator) RegisterRules(key string, fn ValidatorFunc) {
	v.customRules[key] = fn
}

func (v *Validator) RegisterCustomMessage(key string, tmpl string) {
	tpl := CompileTemplate(tmpl)

	v.customMessage[key] = tpl

}

func valueToString(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	switch v.Kind() {

	case reflect.String:
		return v.String()

	case reflect.Bool:
		return strconv.FormatBool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)

	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)

	case reflect.Slice:
		// handle []byte
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes())
		}
	}

	if v.CanInterface() {
		return fmt.Sprint(v.Interface())
	}

	return ""
}
