package conf

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server CfgServer `mapstructure:"server"`
	Minio  CfgMinio  `mapstructure:"minio"`
}

type CfgServer struct {
	Name string `mapstructure:"name"`
}

type CfgMinio struct {
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	RootPath  string `mapstructure:"-"`
}

func LoadConfig() Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("parse config file failed: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
