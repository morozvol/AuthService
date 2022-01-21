package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DB       DB
	BindAddr string `mapstructure:"bind_addr"`
	LogLevel string `mapstructure:"log_level"`
	JWTKey   string
}

type DB struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string
}

func Init() (*Config, error) {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("db", &cfg.DB); err != nil {
		return nil, err
	}

	if err := cfg.ParseEnv(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
func (c *Config) ParseEnv() error {

	if err := viper.BindEnv("db_password"); err != nil {
		return err
	}
	if err := viper.BindEnv("jwt_key"); err != nil {
		return err
	}
	c.JWTKey = viper.GetString("jwt_key")
	c.DB.Password = viper.GetString("db_password")
	return nil
}
func (db *DB) GetConnactionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", db.User, db.Password, db.Host, db.Port, db.Name)
}
