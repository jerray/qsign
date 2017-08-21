package qsign

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type field struct {
	name       string
	value      string
	idx        int
	stringable bool
}

type stringable interface {
	String() string
}

var (
	tags             = []string{"qsign", "json", "yaml", "xml"}
	typeInfoLock     sync.RWMutex
	typeInfoMap      = make(map[reflect.Type][]*field)
	typeOfStringable = reflect.TypeOf((*stringable)(nil)).Elem()
)

// getStructValues
func getStructValues(v interface{}) []*field {
	val := reflect.ValueOf(v)
	vs := []*field{}

	fields := parseStruct(val.Type())

	for _, f := range fields {
		value := getStringValue(val.Field(f.idx), f.stringable)
		vs = append(vs, &field{
			name:  f.name,
			value: value,
		})
	}

	return vs
}

func getStringValue(val reflect.Value, isStringable bool) (r string) {
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return r
		}
		val = val.Elem()
	}

	if isStringable {
		return val.String()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits())
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	}

	return r
}

// parseStruct
func parseStruct(typ reflect.Type) []*field {
	typeInfoLock.RLock()
	fields, ok := typeInfoMap[typ]
	typeInfoLock.RUnlock()
	if ok {
		return fields
	}

	fields = parseFieldsFromType(typ)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].name < fields[j].name
	})

	typeInfoLock.Lock()
	typeInfoMap[typ] = fields
	typeInfoLock.Unlock()
	return fields
}

// parseFieldsFromType
func parseFieldsFromType(typ reflect.Type) []*field {
	res := []*field{}
	typ = findFinalType(typ)

	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			ft := f.Type

			var name string
			if name = getFieldName(f); len(name) == 0 {
				continue
			}

			var canString, canConvert bool
			if isStringable(ft) {
				canString = true
			}

			if isConvertable(ft) {
				canConvert = true
			}

			if canString || canConvert {
				res = append(res, &field{name: name, idx: i, stringable: canString})
			}
		}
	}

	return res
}

func getFieldName(field reflect.StructField) string {
	var ok bool
	var v string

	for _, tag := range tags {
		v, ok = field.Tag.Lookup(tag)
		if ok {
			break
		}
	}

	if !ok {
		return field.Name
	}

	if v == "-" {
		return ""
	}

	return strings.Split(v, ",")[0]
}

func findFinalType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Interface || typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func isStringable(typ reflect.Type) bool {
	if typ.Kind() == reflect.String {
		return true
	}

	if typ.Implements(typeOfStringable) {
		return true
	}

	if typ.Kind() == reflect.Interface || typ.Kind() == reflect.Ptr {
		return isStringable(typ.Elem())
	}

	return false
}

func isConvertable(typ reflect.Type) bool {
	typ = findFinalType(typ)

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	default:
		return false
	}
}
