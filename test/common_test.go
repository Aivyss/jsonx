package test

import (
	"github.com/aivyss/jsonx"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// before a test
	jsonx.Close()

	// run test
	exitCode := m.Run()

	// after a test
	jsonx.Close()
	os.Exit(exitCode)
}
