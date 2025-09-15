package log

import (
	"time"

	"github.com/spf13/pflag"
)

type Options struct {
	Dir          string        `json:"dir" mapstructure:"dir" comment:"日志目录"`
	ID           string        `json:"id" mapstructure:"id" comment:"实例ID"`
	Name         string        `json:"name" mapstructure:"name" comment:"实例名称"`
	Version      string        `json:"version" mapstructure:"version" comment:"实例版本"`
	Debug        bool          `json:"debug" mapstructure:"debug" comment:"是否开启调试模式"`
	MaxAge       time.Duration `json:"max_age" mapstructure:"max_age" comment:"日志最大保存时间"`
	RotationTime time.Duration `json:"rotation_time" mapstructure:"rotation_time" comment:"日志滚动时间"`
	RotationSize int64         `json:"rotation_size" mapstructure:"rotation_size" comment:"日志滚动大小"` // 单位字节
	Level        string        `json:"level" mapstructure:"level" comment:"日志级别"`                   // debug/info/warn/error
	TickSec      int           `json:"tick_sec" mapstructure:"tick_sec" comment:"时间窗口(秒)"`
	First        int           `json:"first" mapstructure:"first" comment:"每个时间窗口内记录的前N条日志"`
	Thereafter   int           `json:"thereafter" mapstructure:"thereafter" comment:"超过N条后每M条记录一次"` // 采样器
}

func NewOptions() *Options {
	return &Options{
		ID:           "test",
		Dir:          "./logs",
		Version:      "0.0.1",
		Debug:        true,
		MaxAge:       7 * 24 * time.Hour,
		RotationTime: 1 * time.Hour,
		RotationSize: 1 * 1024 * 1024,
		Level:        "debug",
		TickSec:      1,
		First:        5,
		Thereafter:   5,
	}
}

func (o *Options) Validate() []error {
	errs := []error{}
	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Dir, "log.dir", o.Dir, "日志目录")
	fs.StringVar(&o.ID, "log.id", o.ID, "实例ID")
	fs.StringVar(&o.Name, "log.name", o.Name, "实例名称")
	fs.StringVar(&o.Version, "log.version", o.Version, "实例版本")
	fs.BoolVar(&o.Debug, "log.debug", o.Debug, "是否开启调试模式")
	fs.DurationVar(&o.MaxAge, "log.max_age", o.MaxAge, "日志最大保存时间")
	fs.DurationVar(&o.RotationTime, "log.rotation_time", o.RotationTime, "日志滚动时间")
	fs.Int64Var(&o.RotationSize, "log.rotation_size", o.RotationSize, "日志滚动大小")
	fs.StringVar(&o.Level, "log.level", o.Level, "日志级别")
	fs.IntVar(&o.TickSec, "log.tick_sec", o.TickSec, "采样器时间窗口(秒)")
	fs.IntVar(&o.First, "log.first", o.First, "采样器每个时间窗口内记录的前N条日志")
	fs.IntVar(&o.Thereafter, "log.thereafter", o.Thereafter, "采样器超过N条后每M条记录一次")
}
