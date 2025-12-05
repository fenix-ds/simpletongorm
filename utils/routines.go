package utils

import (
	"reflect"
	"regexp"
)

func IsStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func IsSnakeCase(s string) bool {
	snakeCaseRegex := regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)
	return snakeCaseRegex.MatchString(s)
}
