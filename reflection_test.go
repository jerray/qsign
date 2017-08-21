package qsign

import (
	"reflect"
	"testing"
)

type myString string

func (s myString) String() string {
	return string(s)
}

type structTypeForTest struct {
	Name  string
	Value int64   `qsign:"value"`
	Skip  bool    `qsign:"-"`
	Addr  *string `qsign:"address"`
	MyStr myString
	Json  string `json:"support_json_tag,omitempty"`
}

func TestReflectionIsStringable(t *testing.T) {
	var s myString = "this is my string var"
	var realString string = "this is a string type var"

	cases := []struct {
		input  reflect.Type
		expect bool
	}{
		{reflect.TypeOf([3]int{1, 2, 3}), false},
		{reflect.TypeOf([]int{}), false},
		{reflect.TypeOf(struct{}{}), false},
		{reflect.TypeOf(0), false},
		{reflect.TypeOf(uint(0)), false},
		{reflect.TypeOf(0.0), false},
		{reflect.TypeOf(false), false},
		{reflect.TypeOf(structTypeForTest{}), false},
		{reflect.TypeOf(""), true},
		{reflect.TypeOf(&realString), true},
		{reflect.TypeOf(realString), true},
		{reflect.TypeOf(&s), true},
		{reflect.TypeOf(s), true},
	}

	for _, c := range cases {
		actual := isStringable(c.input)
		if actual != c.expect {
			t.Errorf("expect type `%s` is not stringable", c.input)
		}
	}
}

func TestReflectionGetStringValue(t *testing.T) {
	var s myString = "this is my string var"
	var realString string = "this is a string type var"
	var nilString *string

	cases := []struct {
		input        reflect.Value
		isStringable bool
		expect       string
	}{
		{reflect.ValueOf([3]int{1, 2, 3}), false, ""},
		{reflect.ValueOf([]int{}), false, ""},
		{reflect.ValueOf(struct{}{}), false, ""},
		{reflect.ValueOf(int(0)), false, "0"},
		{reflect.ValueOf(uint(1)), false, "1"},
		{reflect.ValueOf(0.1), false, "0.1"},
		{reflect.ValueOf(0.0), false, "0"},
		{reflect.ValueOf(false), false, "false"},
		{reflect.ValueOf(""), true, ""},
		{reflect.ValueOf(&realString), true, realString},
		{reflect.ValueOf(realString), true, realString},
		{reflect.ValueOf(&s), true, s.String()},
		{reflect.ValueOf(s), true, s.String()},
		{reflect.ValueOf(nilString), true, ""},
	}

	for i, c := range cases {
		actual := getStringValue(c.input, c.isStringable)
		if actual != c.expect {
			t.Errorf("expect index %d value is `%s`, actual is `%s`", i, actual)
		}
	}
}

func TestReflectionIsConvertable(t *testing.T) {
	cases := []struct {
		input  reflect.Type
		expect bool
	}{
		{reflect.TypeOf([3]int{1, 2, 3}), false},
		{reflect.TypeOf([]int{}), false},
		{reflect.TypeOf(struct{}{}), false},
		{reflect.TypeOf(""), true},
		{reflect.TypeOf(0), true},
		{reflect.TypeOf(uint(0)), true},
		{reflect.TypeOf(0.0), true},
		{reflect.TypeOf(false), true},
	}

	for _, c := range cases {
		actual := isConvertable(c.input)
		if actual != c.expect {
			t.Errorf("expect type %s is not convertable", c.input)
		}
	}
}

func TestReflectionGetFieldName(t *testing.T) {
	s := structTypeForTest{}

	typ := reflect.TypeOf(s)
	cases := []struct {
		input  reflect.StructField
		expect string
	}{
		{typ.Field(0), "Name"},
		{typ.Field(1), "value"},
		{typ.Field(2), ""},
		{typ.Field(3), "address"},
		{typ.Field(4), "MyStr"},
		{typ.Field(5), "support_json_tag"},
	}

	for _, c := range cases {
		actual := getFieldName(c.input)
		if actual != c.expect {
			t.Errorf("expect field name is `%s`, actual is `%s`", c.expect, actual)
		}
	}
}

func TestReflectionParseFieldsFromType(t *testing.T) {
	input := structTypeForTest{}
	expect := []*field{
		&field{name: "Name", idx: 0, stringable: true},
		&field{name: "value", idx: 1, stringable: false},
		&field{name: "address", idx: 3, stringable: true},
		&field{name: "MyStr", idx: 4, stringable: true},
		&field{name: "support_json_tag", idx: 5, stringable: true},
	}

	actual := parseFieldsFromType(reflect.TypeOf(input))
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("expect parse result equals, expect %#v, actual %#v", expect, actual)
	}
}

func TestReflectionParseStruct(t *testing.T) {
	input := structTypeForTest{}
	expect := []*field{
		&field{name: "MyStr", idx: 4, stringable: true},
		&field{name: "Name", idx: 0, stringable: true},
		&field{name: "address", idx: 3, stringable: true},
		&field{name: "support_json_tag", idx: 5, stringable: true},
		&field{name: "value", idx: 1, stringable: false},
	}

	typ := reflect.TypeOf(input)

	if _, ok := typeInfoMap[typ]; ok {
		t.Errorf("expect typeInfoMap has no items")
	}

	actual := parseStruct(typ)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("expect parse result equals, expect %#v, actual %#v", expect, actual)
	}

	if _, ok := typeInfoMap[typ]; !ok {
		t.Errorf("expect typeInfoMap has type cache")
	}

	actual = parseStruct(typ)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("expect parse result equals, expect %#v, actual %#v", expect, actual)
	}
}

func TestReflectionGetStructValues(t *testing.T) {
	var s myString = "this is my string var"
	var realString string = "this is a string type var"

	input := structTypeForTest{
		Name:  "",
		Value: 7,
		Skip:  true,
		Addr:  &realString,
		MyStr: s,
	}

	expect := []*field{
		&field{name: "MyStr", value: s.String()},
		&field{name: "Name", value: input.Name},
		&field{name: "address", value: realString},
		&field{name: "support_json_tag", value: ""},
		&field{name: "value", value: "7"},
	}

	actual := getStructValues(input)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("expect parse result equals, expect %#v, actual %#v", expect, actual)
	}
}
