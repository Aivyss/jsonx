package test

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/jsonx/errors"
	"testing"
)

func TestFieldError(t *testing.T) {
	msg := "test_msg"
	name := "error_name"
	fieldErr := errors.NewFieldErr(name, msg)

	if fieldErr.DefaultMsg() != msg {
		t.Fatal("unexpected result1")
	}

	if fieldErr.Name() != name {
		t.Fatal("unexpected result2")
	}

	j, err := jsonx.Marshal(struct {
		FrameworkName string `json:"framework"`
		Name          string `json:"errorName"`
		Msg           string `json:"Msg"`
	}{
		FrameworkName: "jsonx",
		Name:          name,
		Msg:           msg,
	})
	if err != nil {
		t.Fatal("unexpected result3")
	}
	if fieldErr.Error() != string(j) {
		t.Fatal("unexpected result4")
	}
}
