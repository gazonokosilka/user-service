package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server    ServerConfig   `yaml:"server"`
	Postgres  PostgresConfig `yaml:"postgres"`
	SecretKey string         `yaml:"secret_key"`
	//RedisConfig RedisConfig    `yaml:"redis"`
	TokenTTL time.Duration `yaml:"token_ttl"`
}

type ServerConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Path     string `yaml:"path"` // TODO: may be delete
	SslMode  string `yaml:"ssl_mode"`
}

//type RedisConfig struct {
//	Addr     string `yaml:"addr"`
//	Password string `yaml:"password"`
//	DB       int    `yaml:"db"`
//}

func MustLoad() *Config {

	path := fetchConfigPath()
	fmt.Println(path)
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flag or env.
// Priority: env > flag
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
