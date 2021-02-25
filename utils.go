package copy

import (
	"reflect"
)

func deepFields(typ reflect.Type, depth int) []StructField {
	typ = indirectType(typ)
	num := typ.NumField()
	var fields []StructField
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		structField := StructField{
			field,
			depth,
		}
		if field.Anonymous && field.Type.Kind() != reflect.Interface {
			fields = append(fields, deepFields(field.Type, depth+1)...)
		} else {
			fields = append(fields, structField)
		}
	}
	return fields
}

func indirectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	return typ
}
