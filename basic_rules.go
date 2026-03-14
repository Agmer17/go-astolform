package astolform

import (
	"reflect"
)

type ValidatorFunc func(val reflect.Value, param string, fieldName string) error

var (
	basicRules = map[string]ValidatorFunc{
		"required": requiredRules,
	}
)

func RegisterRules(key string, fn ValidatorFunc) {
	basicRules[key] = fn
}

func requiredRules(val reflect.Value, _ string, fieldName string) error {
	if val.IsZero() {
		return GenerateErrorMessage(fieldName, " is required!")
	}
	return nil
}

// func MaxRules(val reflect.Value, params string, fieldName string) error {

// 	limit, err := strconv.Atoi(params)

// 	if err != nil {
// 	  return GenerateErrorMessage(fieldName, )
// 	}

// }
