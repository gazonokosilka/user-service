package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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

	// Переопределение из переменных окружения (приоритет над файлом)
	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		cfg.Postgres.Host = host
		fmt.Printf("Override Postgres Host from ENV: %s\n", host)
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Postgres.Port = p
			fmt.Printf("Override Postgres Port from ENV: %d\n", p)
		}
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.Postgres.User = user
		fmt.Printf("Override Postgres User from ENV: %s\n", user)
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cfg.Postgres.Password = password
		fmt.Println("Override Postgres Password from ENV: ***")
	}
	if db := os.Getenv("POSTGRES_DB"); db != "" {
		cfg.Postgres.DbName = db
		fmt.Printf("Override Postgres DB from ENV: %s\n", db)
	}
	if sslmode := os.Getenv("POSTGRES_SSLMODE"); sslmode != "" {
		cfg.Postgres.SslMode = sslmode
		fmt.Printf("Override Postgres SSLMode from ENV: %s\n", sslmode)
	}

	fmt.Printf("Final Postgres Config: %s:%d/%s (user: %s, sslmode: %s)\n",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DbName, cfg.Postgres.User, cfg.Postgres.SslMode)

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
