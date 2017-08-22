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
	idx        []int
	stringable bool
}

// stringable interface is used to check if a type has String() function.
type stringable interface {
	String() string
}

var (
	tags             = []string{"qsign", "json", "yaml", "xml"}
	typeInfoLock     sync.RWMutex
	typeInfoMap      = make(map[reflect.Type][]*field)
	typeOfStringable = reflect.TypeOf((*stringable)(nil)).Elem()
)

// getStructValues parses interface v, returns its field list with fields' string value.
func getStructValues(v interface{}) (vs []*field) {
	vs = []*field{}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if !val.IsValid() {
		return
	}

	fields := parseStruct(val.Type())

	for _, f := range fields {
		value := getStringValue(val, f.idx, f.stringable)
		vs = append(vs, &field{
			name:  f.name,
			value: value,
		})
	}

	return
}

// getStringValue returns the string value of val. If depth is great than 0, val must be a struct,
// it will find string value from the next depth level.
func getStringValue(val reflect.Value, depth []int, isStringable bool) (r string) {
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return r
		}
		val = val.Elem()
	}

	if len(depth) > 0 {
		return getStringValue(val.Field(depth[0]), depth[1:], isStringable)
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

// parseStruct parses input type, store its field list in a map for use in the next time.
func parseStruct(typ reflect.Type) []*field {
	typeInfoLock.RLock()
	fields, ok := typeInfoMap[typ]
	typeInfoLock.RUnlock()
	if ok {
		return fields
	}

	fields = parseFieldsFromType(typ, []int{})
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].name < fields[j].name
	})

	typeInfoLock.Lock()
	typeInfoMap[typ] = fields
	typeInfoLock.Unlock()
	return fields
}

// parseFieldsFromType parses the input type, returns its field list.
func parseFieldsFromType(typ reflect.Type, idx []int) []*field {
	res := []*field{}
	typ = findFinalType(typ)

	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			ft := findFinalType(f.Type)

			name, skip := getFieldName(f)
			if skip || len(name) == 0 {
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
				res = append(res, &field{name: name, idx: append(idx, i), stringable: canString})
			} else if ft.Kind() == reflect.Struct && f.Anonymous {
				res = append(res, parseFieldsFromType(ft, append(idx, i))...)
			}
		}
	}

	return res
}

// getFieldName returns a struct field's name according to field's tag.
func getFieldName(field reflect.StructField) (v string, skip bool) {
	var ok bool

	for _, tag := range tags {
		v, ok = field.Tag.Lookup(tag)
		if ok {
			break
		}
	}

	if !ok {
		v = field.Name
		return
	}

	if v == "-" {
		return "", true
	}

	v = strings.Split(v, ",")[0]
	return
}

func findFinalType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Interface || typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

// isStringable checks if a type is string or implements stringable interface.
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

// isConvertable checks if a type can be converted to string.
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
