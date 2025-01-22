package config

import (
	"errors"
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type SMTP struct {
	Username string `toml:"username" env:"SMTP_USERNAME"`
	Password string `toml:"password" env:"SMTP_PASSWORD"`
	Host     string `toml:"host" env:"SMTP_HOST"`
}

type Redis struct {
	Host          string `toml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port          int    `toml:"port" env:"REDIS_PORT" env-default:"6379"`
	DB            int    `toml:"db" env:"REDIS_DB" env-default:"0"`
	User          string `toml:"user" env:"REDIS_USER"`
	Password      string `toml:"password" env:"REDIS_USER_PASSWORD"`
	RedisPassword string `toml:"redis_password" env:"REDIS_PASSWORD"`
}

type MailConfig struct {
	SMTP       SMTP   `toml:"smtp"`
	ServerURL  string `toml:"server" env:"SERVER_URL"`
	Redis      Redis  `toml:"redis"`
	configPath string
}

func NewConfigService() *MailConfig {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		if cfg_path, ok := os.LookupEnv("CONFIG_PATH"); ok {
			path = cfg_path
		}
	}

	cfg := &MailConfig{configPath: path}
	if err := cfg.load(); err != nil {
		panic(err)
	}
	return cfg
}

func (cfg *MailConfig) load() error {
	if cfg.configPath == "" {
		return errors.New("config path is not set")
	}

	if err := cleanenv.ReadConfig(cfg.configPath, cfg); err != nil {
		return err
	}

	return nil
}
