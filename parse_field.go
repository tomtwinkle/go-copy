package copy

import (
	"reflect"
	"strings"
)

type StructField struct {
	reflect.StructField
	depth int
}

type FieldParseFunc func(field StructField) string

func ParseFiledByName(field StructField) string {
	return field.Name
}

func ParseFieldByJSONTag(field StructField) string {
	tag := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if tag == "-" {
		return ""
	}
	return tag
}

func ParseFieldByCopyTag(field StructField) string {
	tag := strings.SplitN(field.Tag.Get("copy"), ",", 2)[0]
	if tag == "-" {
		return ""
	}
	return tag
}
