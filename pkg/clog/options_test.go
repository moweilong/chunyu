package clog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Options_Validate(t *testing.T) {
	opts := &Options{
		Name:              "test",
		Level:             "test",
		Format:            "test",
		Development:       true,
		EnableColor:       true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     false,
		DisableStacktrace: false,
	}

	errs := opts.Validate()
	expected := `[unrecognized level: "test" not a valid log format: "test"]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
