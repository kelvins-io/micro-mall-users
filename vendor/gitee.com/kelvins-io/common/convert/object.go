package convert

import (
	"reflect"
	"strings"
)

// parse struct of object
// tagType 1 => "db"	2=> "json"
func ParseObject(input interface{}, tagType int) map[string]string {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	var m = make(map[string]string)

	return parseObject(v, t, m, tagType)
}

// t.Field(i).Name : field name
// t.Field(i).Type : field type
// v.Field(i).Interface(): field value
// t.Field(i).Tag: struct tag
func parseObject(v reflect.Value, t reflect.Type, m map[string]string, tagType int) map[string]string {
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			if v.Field(i).Type().Kind() == reflect.Struct {
				parseObject(v.Field(i), v.Field(i).Type(), m, tagType)
				continue
			}

			m[getTagByType(tagType, t.Field(i).Tag)] = t.Field(i).Type.String()
		}
	}

	return m
}

// get tag string
// tagType 1 => "db"	2=> "json"
func getTagByType(tagType int, tag reflect.StructTag) string {
	switch tagType {
	case 1:
		return tag.Get("db")
	case 2:
		var tmp = tag.Get("json")
		return strings.Replace(tmp, ",omitempty", "", 1)
	default:
		return ""
	}
}
