// Package options provides flags and configuration for initializing the Onex Cache Server.
package options

import (
	"chunyu/pkg/app"
	"chunyu/pkg/log"
	genericoptions "chunyu/pkg/options"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
)

// ServerOptions contains the configuration options for the server.
type ServerOptions struct {
	GRPCOptions  *genericoptions.GRPCOptions  `json:"grpc" mapstructure:"grpc"`
	TLSOptions   *genericoptions.TLSOptions   `json:"tls" mapstructure:"tls"`
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	Logging      *log.Options                 `json:"log" mapstructure:"log"`
}

// Ensure ServerOptions implements the app.NamedFlagSetOptions interface.
var _ app.NamedFlagSetOptions = (*ServerOptions)(nil)

// NewServerOptions creates a ServerOptions instance with default values.
func NewServerOptions() *ServerOptions {
	o := &ServerOptions{
		DisableCache: false,
		GRPCOptions:  genericoptions.NewGRPCOptions(),
		TLSOptions:   genericoptions.NewTLSOptions(),
		RedisOptions: genericoptions.NewRedisOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
		Logging:      log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *ServerOptions) Flags() (fss cliflag.NamedFlagSets) {
	// Add flags for each option group with meaningful section names.
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.TLSOptions.AddFlags(fss.FlagSet("tls"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.Logging.AddFlags(fss.FlagSet("log"))

	// Add a miscellaneous flag for the cache control feature.
	miscFs := fss.FlagSet("misc")
	miscFs.BoolVar(&o.DisableCache, "disable-cache", o.DisableCache, "Disable the local memory cache.")

	return fss
}

// Validate checks whether the options in ServerOptions are valid.
func (o *ServerOptions) Validate() error {
	var errs []error

	// Perform validation for each option group, accumulating errors.
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Logging.Validate()...)

	// Aggregate all validation errors into a single error object.
	return utilerrors.NewAggregate(errs)
}

// Config builds an cacheserver.Config based on ServerOptions.
func (o *ServerOptions) Config() (*cacheserver.Config, error) {
	// Ensure the configuration includes all relevant fields from the options.
	return &cacheserver.Config{
		DisableCache:  o.DisableCache,
		GRPCOptions:   o.GRPCOptions,
		TLSOptions:    o.TLSOptions,
		RedisOptions:  o.RedisOptions,
		MySQLOptions:  o.MySQLOptions,
		JaegerOptions: o.JaegerOptions,
	}, nil
}
