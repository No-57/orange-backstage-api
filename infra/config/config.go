package config

import "time"

type App struct {
	Server Server `mapstructure:"server"`
	Log    Log    `mapstructure:"log"`
}

type Server struct {
	RunMode      string        `mapstructure:"run_mode"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	CertFilePath string        `mapstructure:"cert_file_path"`
	KeyFilePath  string        `mapstructure:"key_file_path"`

	JWT JWT `mapstructure:"jwt"`
}

type JWT struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
}

type Log struct {
	Level        string `mapstructure:"level"`
	FileName     string `mapstructure:"file_name"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Compress     bool   `mapstructure:"compress"`
	ConsoleDebug bool   `mapstructure:"console_debug"`
}
