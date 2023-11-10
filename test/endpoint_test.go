package test

import (
	"errors"
	"github.com/aivyss/jsonx"
	"strings"
	"testing"
)

type testStruct struct {
	Value string `json:"value"`
}

type testStructValidator struct{}

func (v *testStructValidator) Validate(t testStruct) error {
	if strings.TrimSpace(t.Value) == "" {
		return errors.New("blank string is not allowed")
	}

	return nil
}

type orderedValidator1 int
type orderedValidator2 int

func (v *orderedValidator1) Validate(o testStruct) error {
	if !strings.Contains(o.Value, "apple") {
		return errors.New("apple is not contained")
	}

	return nil
}
func (v *orderedValidator1) Order() int {
	return 1
}

func (v *orderedValidator2) Validate(o testStruct) error {
	if !strings.Contains(o.Value, "banana") {
		return errors.New("banana is not contained")
	}

	return nil
}
func (v *orderedValidator2) Order() int {
	return 2
}

func TestUnmarshal(t *testing.T) {
	t.Run("[pass validation - normal]", func(t *testing.T) {
		jsonx.RegisterValidator[testStruct](&testStructValidator{})
		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": "test_string"}`))

		if err != nil {
			t.Fatal(err)
		}

		if "test_string" != o.Value {
			t.Fatal("not equal string")
		}
	})

	t.Run("[fail to validate - normal]", func(t *testing.T) {
		jsonx.RegisterValidator[testStruct](&testStructValidator{})
		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": ""}`))

		if err == nil || o != nil {
			t.Fatal(err)
		}
	})

	t.Run("[pass validation - ordered]", func(t *testing.T) {
		validator1 := orderedValidator1(1)
		validator2 := orderedValidator2(1)
		jsonx.RegisterOrderedValidator[testStruct](&validator1)
		jsonx.RegisterOrderedValidator[testStruct](&validator2)
		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": "apple,banana"}`))

		if err != nil {
			t.Fatal(err)
		}

		if "apple,banana" != o.Value {
			t.Fatal("not equal string")
		}
	})

	t.Run("[fail to validate - ordered1]", func(t *testing.T) {
		validator1 := orderedValidator1(1)
		validator2 := orderedValidator2(1)
		jsonx.RegisterOrderedValidator[testStruct](&validator1)
		jsonx.RegisterOrderedValidator[testStruct](&validator2)
		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": "graph,banana"}`))

		if err == nil || o != nil || err.Error() != "apple is not contained" {
			t.Fatal("unexpected result")
		}
	})

	t.Run("[fail to validate - ordered2]", func(t *testing.T) {
		validator1 := orderedValidator1(1)
		validator2 := orderedValidator2(1)
		jsonx.RegisterOrderedValidator[testStruct](&validator1)
		jsonx.RegisterOrderedValidator[testStruct](&validator2)
		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": "apple,graph"}`))

		if err == nil || o != nil || err.Error() != "banana is not contained" {
			t.Fatal("unexpected result")
		}
	})
}
