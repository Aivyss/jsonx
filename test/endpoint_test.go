package test

import (
	"errors"
	"github.com/aivyss/jsonx"
	jsonxErr "github.com/aivyss/jsonx/errors"
	"github.com/aivyss/typex/util"
	"strings"
	"testing"
	"time"
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

func TestAnnotation(t *testing.T) {
	t.Run("[Required]", func(t *testing.T) {
		type testStruct struct {
			Value *string `json:"value" annotation:"@Required"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "test_value" }`))
		if err != nil {
			t.Fatal("unexpected result")
		}
	})

	t.Run("[NotBlank]", func(t *testing.T) {
		type testStruct struct {
			Value *string `json:"value" annotation:"@NotBlank"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "" }`))
		if err == nil {
			t.Fatal("unexpected result2")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "    " }`))
		if err == nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "test_value" }`))
		if err != nil {
			t.Fatal("unexpected result4")
		}
	})

	t.Run("[NotEmpty]", func(t *testing.T) {
		type testStruct struct {
			Value *string `json:"value" annotation:"@NotEmpty"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "" }`))
		if err == nil {
			t.Fatal("unexpected result2")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "    " }`))
		if err != nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "test_value" }`))
		if err != nil {
			t.Fatal("unexpected result4")
		}
	})

	t.Run("[Email]", func(t *testing.T) {
		type testStruct1 struct {
			Value *string `json:"value" annotation:"@Email"`
		}
		type testStruct2 struct {
			Value string `json:"value" annotation:"@Email"`
		}

		_, err := jsonx.Unmarshal[testStruct1]([]byte(`{ "value": "test@example.com" }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct1]([]byte(`{ "value": "test@example.co.jp" }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct1]([]byte(`{ "value": "test@sub.example.co.kr" }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": "test@example.com" }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": "test@example.co.jp" }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": "test@sub.example.co.kr" }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}
	})

	t.Run("[NotContainsNil]", func(t *testing.T) {
		type testStruct struct {
			Value []*string `json:"value" annotation:"@NotContainsNil"`
		}
		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": [null, "a", "b"] }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}

		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": ["a", "b"] }`))
		if err != nil {
			t.Fatal("unexpected result2")
		}

		for _, v := range o.Value {
			if *v != "a" && *v != "b" {
				t.Fatal("unexpected result3")
			}
		}
	})

	t.Run("[NotContainsEmpty]", func(t *testing.T) {
		type testStruct struct {
			Value []*string `json:"value" annotation:"@NotContainsEmpty"`
		}
		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": [null, "a", "b"] }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": ["", "a", "b"] }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": ["c", "a", "b"] }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": [" ", "a", "b"] }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}
	})

	t.Run("[NotContainsBlank]", func(t *testing.T) {
		type testStruct struct {
			Value []*string `json:"value" annotation:"@NotContainsBlank"`
		}
		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": [null, "a", "b"] }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": ["", "a", "b"] }`))
		if err == nil {
			t.Fatal("unexpected result2")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": ["c", "a", "b"] }`))
		if err != nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": [" ", "a", "b"] }`))
		if err == nil {
			t.Fatal("unexpected result4")
		}
	})

	t.Run("[customAnnotation]", func(t *testing.T) {
		type testStruct struct {
			Value string `json:"value" annotation:"@Banana"`
		}

		if err := jsonx.RegisterCustomAnnotation("Banana", func(v any) error {
			if v.(string) != "banana" {
				return errors.New("test error")
			}

			return nil
		}); err != nil {
			t.Fatal("unexpected result1")
		}

		o, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": "banana" }`))

		if err != nil || o.Value != "banana" {
			t.Fatal("unexpected result2")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": "apple" }`))
		if err == nil {
			t.Fatal("unexpected result3")
		}
	})

	t.Run("[Positive]", func(t *testing.T) {
		type testStruct struct {
			Value int `json:"value" annotation:"@Positive"`
		}
		type testStruct2 struct {
			Value *float32 `json:"value" annotation:"@Positive"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": 1 }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": 0 }`))
		if err == nil {
			t.Fatal("unexpected result2")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": -1 }`))
		if err == nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 1 }`))
		if err != nil {
			t.Fatal("unexpected result4")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 0 }`))
		if err == nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": -1 }`))
		if err == nil {
			t.Fatal("unexpected result6")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result7")
		}
	})

	t.Run("[PositiveOrZero]", func(t *testing.T) {
		type testStruct struct {
			Value int `json:"value" annotation:"@PositiveOrZero"`
		}
		type testStruct2 struct {
			Value *float32 `json:"value" annotation:"@PositiveOrZero"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": 1 }`))
		if err != nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": 0 }`))
		if err != nil {
			t.Fatal("unexpected result2")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": -1 }`))
		if err == nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 1 }`))
		if err != nil {
			t.Fatal("unexpected result4")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 0 }`))
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": -1 }`))
		if err == nil {
			t.Fatal("unexpected result6")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result7")
		}
	})

	t.Run("[Negative]", func(t *testing.T) {
		type testStruct struct {
			Value int `json:"value" annotation:"@Negative"`
		}
		type testStruct2 struct {
			Value *float32 `json:"value" annotation:"@Negative"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": 1 }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": 0 }`))
		if err == nil {
			t.Fatal("unexpected result2")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": -1 }`))
		if err != nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 1 }`))
		if err == nil {
			t.Fatal("unexpected result4")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 0 }`))
		if err == nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": -1 }`))
		if err != nil {
			t.Fatal("unexpected result6")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result7")
		}
	})

	t.Run("[NegativeOrZero]", func(t *testing.T) {
		type testStruct struct {
			Value int `json:"value" annotation:"@NegativeOrZero"`
		}
		type testStruct2 struct {
			Value *float32 `json:"value" annotation:"@NegativeOrZero"`
		}

		_, err := jsonx.Unmarshal[testStruct]([]byte(`{ "value": 1 }`))
		if err == nil {
			t.Fatal("unexpected result1")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": 0 }`))
		if err != nil {
			t.Fatal("unexpected result2")
		}
		_, err = jsonx.Unmarshal[testStruct]([]byte(`{ "value": -1 }`))
		if err != nil {
			t.Fatal("unexpected result3")
		}

		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 1 }`))
		if err == nil {
			t.Fatal("unexpected result4")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": 0 }`))
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": -1 }`))
		if err != nil {
			t.Fatal("unexpected result6")
		}
		_, err = jsonx.Unmarshal[testStruct2]([]byte(`{ "value": null }`))
		if err == nil {
			t.Fatal("unexpected result7")
		}
	})

	t.Run("[Future]", func(t *testing.T) {
		present := time.Now()
		future := present.Add(1 * time.Second)
		past := future.Add(-2 * time.Second)
		type testStruct struct {
			Value time.Time `json:"value" annotation:"@Future"`
		}
		j, err := jsonx.Marshal(testStruct{Value: future})
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result2")
		}

		j, err = jsonx.Marshal(testStruct{Value: past})
		if err != nil {
			t.Fatal("unexpected result3")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result4")
		}

		j, err = jsonx.Marshal(testStruct{Value: present})
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result6")
		}
	})

	t.Run("[Present]", func(t *testing.T) {
		present := time.Now()
		future := present.Add(1 * time.Second)
		past := future.Add(-2 * time.Second)
		type testStruct struct {
			Value time.Time `json:"value" annotation:"@Present"`
		}
		// future
		j, err := jsonx.Marshal(testStruct{Value: future})
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result2")
		}

		// past
		j, err = jsonx.Marshal(testStruct{Value: past})
		if err != nil {
			t.Fatal("unexpected result3")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result4")
		}

		// present
		j, err = jsonx.Marshal(testStruct{Value: present})
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result6")
		}
	})

	t.Run("[Past]", func(t *testing.T) {
		present := time.Now()
		future := present.Add(1 * time.Second)
		past := present.Add(-1 * time.Second)
		type testStruct struct {
			Value time.Time `json:"value" annotation:"@Past"`
		}
		// future
		j, err := jsonx.Marshal(testStruct{Value: future})
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result2")
		}

		// past
		j, err = jsonx.Marshal(testStruct{Value: past})
		if err != nil {
			t.Fatal("unexpected result3")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result4")
		}

		// present
		j, err = jsonx.Marshal(testStruct{Value: present})
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result6")
		}
	})

	t.Run("[FutureOrPresent]", func(t *testing.T) {
		present := time.Now()
		future := present.Add(1 * time.Second)
		past := present.Add(-1 * time.Second)
		type testStruct struct {
			Value time.Time `json:"value" annotation:"@FutureOrPresent"`
		}
		// future
		j, err := jsonx.Marshal(testStruct{Value: future})
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result2")
		}

		// past
		j, err = jsonx.Marshal(testStruct{Value: past})
		if err != nil {
			t.Fatal("unexpected result3")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result4")
		}

		// present
		j, err = jsonx.Marshal(testStruct{Value: present})
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result6")
		}
	})

	t.Run("[PastOrPresent]", func(t *testing.T) {
		present := time.Now()
		future := present.Add(1 * time.Second)
		past := present.Add(-1 * time.Second)
		type testStruct struct {
			Value time.Time `json:"value" annotation:"@PastOrPresent"`
		}
		// future
		j, err := jsonx.Marshal(testStruct{Value: future})
		if err != nil {
			t.Fatal("unexpected result1")
		}

		_, err = jsonx.Unmarshal[testStruct](j)
		if err == nil {
			t.Fatal("unexpected result2")
		}

		// past
		j, err = jsonx.Marshal(testStruct{Value: past})
		if err != nil {
			t.Fatal("unexpected result3")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result4")
		}

		// present
		j, err = jsonx.Marshal(testStruct{Value: present})
		if err != nil {
			t.Fatal("unexpected result5")
		}
		_, err = jsonx.Unmarshal[testStruct](j)
		if err != nil {
			t.Fatal("unexpected result6")
		}
	})
}

func TestPatternTag(t *testing.T) {
	type testStruct struct {
		Email string `json:"value" pattern:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
	}

	type testStruct2 struct {
		Email *string `json:"value" pattern:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
	}

	j, err := jsonx.Marshal(testStruct{Email: "hklee6417@gmail.com"})
	if err != nil {
		t.Fatal("unexpected result1")
	}
	_, err = jsonx.Unmarshal[testStruct](j)
	if err != nil {
		t.Fatal("unexpected result2")
	}

	j, err = jsonx.Marshal(testStruct{Email: "wrong_string"})
	if err != nil {
		t.Fatal("unexpected result3")
	}
	_, err = jsonx.Unmarshal[testStruct](j)
	if err == nil {
		t.Fatal("unexpected result4")
	}

	j, err = jsonx.Marshal(testStruct2{Email: util.MustPointer("hklee6417@gmail.com")})
	if err != nil {
		t.Fatal("unexpected result5")
	}
	_, err = jsonx.Unmarshal[testStruct2](j)
	if err != nil {
		t.Fatal("unexpected result6")
	}

	j, err = jsonx.Marshal(testStruct2{Email: util.MustPointer("wrong_string")})
	if err != nil {
		t.Fatal("unexpected result7")
	}
	_, err = jsonx.Unmarshal[testStruct2](j)
	if err == nil {
		t.Fatal("unexpected result8")
	}
}

func TestValidateFunction(t *testing.T) {
	type testStruct struct {
		Value string
	}

	v := testStruct{Value: "test_string"}
	err := jsonx.Validate(v)
	if err != nil {
		t.Fatal("unexpected result1")
	}

	err = jsonx.Validate(&v)
	if err == nil {
		t.Fatal("unexpected result2")
	}

	present := time.Now()
	future := present.Add(1 * time.Second)
	past := future.Add(-2 * time.Second)
	type testStruct2 struct {
		Value time.Time `json:"value" annotation:"@Present"`
	}
	// future
	err = jsonx.Validate[testStruct2](testStruct2{Value: future})
	if err == nil {
		t.Fatal("unexpected result3")
	}

	// past
	err = jsonx.Validate[testStruct2](testStruct2{Value: past})
	if err == nil {
		t.Fatal("unexpected result4")
	}

	// present
	err = jsonx.Validate[testStruct2](testStruct2{Value: present})
	if err != nil {
		t.Fatal("unexpected result5")
	}
}

func TestFixFieldErr(t *testing.T) {
	jsonx.Close()
	errName := "testErr"
	msg := "test_msg"
	jsonx.RegisterFieldError(errName, msg)
	present := time.Now()
	future := present.Add(1 * time.Second)
	past := future.Add(-2 * time.Second)
	type testStruct struct {
		Value time.Time `json:"value" annotation:"@Present" fieldErr:"testErr"`
	}
	// future
	err := jsonx.Validate[testStruct](testStruct{Value: future})
	fieldErr, ok := err.(*jsonxErr.FieldError)
	if !ok || err == nil {
		t.Fatal("unexpected result1")
	}
	if fieldErr.DefaultMsg() != msg || fieldErr.Name() != errName {
		t.Fatal("unexpected result2")
	}

	// past
	err = jsonx.Validate[testStruct](testStruct{Value: past})
	fieldErr, ok = err.(*jsonxErr.FieldError)
	if !ok || err == nil {
		t.Fatal("unexpected result3")
	}
	if fieldErr.DefaultMsg() != msg || fieldErr.Name() != errName {
		t.Fatal("unexpected result4")
	}

	// present
	err = jsonx.Validate[testStruct](testStruct{Value: present})
	if err != nil {
		t.Fatal("unexpected result5")
	}
}
