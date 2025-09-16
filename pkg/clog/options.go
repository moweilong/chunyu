package clog

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	jsonFormat    = "json"
	consoleFormat = "console"

	flagName              = "log.name"
	flagLevel             = "log.level"
	flagFormat            = "log.format"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagEnableColor       = "log.enable-color"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
)

// Options contains configuration options for logging.
type Options struct {
	Name string `json:"name,omitempty"               mapstructure:"name"`
	// Level specifies the minimum log level. Valid values are: debug, info, warn, error, dpanic, panic, and fatal.
	Level string `json:"level,omitempty" mapstructure:"level"`
	// Format specifies the log output format. Valid values are: console and json.
	Format string `json:"format,omitempty" mapstructure:"format"`
	// OutputPaths specifies the output paths for the logs.
	OutputPaths      []string `json:"output-paths,omitempty" mapstructure:"output-paths"`
	ErrorOutputPaths []string `json:"error-output-paths,omitempty" mapstructure:"error-output-paths"`
	Development      bool     `json:"development,omitempty"        mapstructure:"development"`
	// EnableColor specifies whether to output colored logs.
	EnableColor bool `json:"enable-color"       mapstructure:"enable-color"`
	// DisableCaller specifies whether to include caller information in the log.
	DisableCaller bool `json:"disable-caller,omitempty" mapstructure:"disable-caller"`
	// DisableStacktrace specifies whether to record a stack trace for all messages at or above panic level.
	DisableStacktrace bool `json:"disable-stacktrace,omitempty" mapstructure:"disable-stacktrace"`
}

// NewOptions creates a new Options object with default values.
func NewOptions() *Options {
	return &Options{
		Name:              "default",
		Level:             zapcore.InfoLevel.String(),
		Format:            consoleFormat,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Development:       false,
		EnableColor:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
	}
}

// Validate verifies flags passed to LogsOptions.
func (o *Options) Validate() []error {
	errs := []error{}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds command line flags for the configuration.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace, o.DisableStacktrace, ""+
		"Disable the log to record a stack trace for all messages at or above panic level.")
}
