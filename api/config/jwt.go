package config

import "time"

type Jwt struct {
	Key            string        `mapstructure:"key"`
	KeyFile        string        `mapstructure:"key-file"`
	ExpireDuration time.Duration `mapstructure:"expire-duration"`
}
