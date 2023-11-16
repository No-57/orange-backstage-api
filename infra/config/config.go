package config

import (
	"orange-backstage-api/infra/util/convert"
	"time"
)

type App struct {
	Server Server `mapstructure:"server"`
	Log    Log    `mapstructure:"log"`
	DB     DB     `mapstructure:"db"`
}

type Server struct {
	RunMode      string        `mapstructure:"run_mode"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	CertFilePath string        `mapstructure:"cert_file_path"`
	KeyFilePath  string        `mapstructure:"key_file_path"`
	EnableDoc    bool          `mapstructure:"enable_doc"`

	JWT JWT `mapstructure:"jwt"`
}

type JWT struct {
	Secret             string        `mapstructure:"secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire"`
}

func (jwt JWT) SecretBytes() []byte {
	return convert.StrToBytes(jwt.Secret)
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

type DB struct {
	Engine   string `mapstructure:"engine"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
	TimeZone string `mapstructure:"timezone"`

	Verbose bool `mapstructure:"verbose"`
}
