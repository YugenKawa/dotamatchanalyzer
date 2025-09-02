package config

import "fmt"

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Database DatabaseConfig `yaml:"database"`
	Settings SettingsConfig `yaml:"settings"`
	Teams    []TeamConfig   `yaml:"teams"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   int64  `yaml:"chat_id"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type SettingsConfig struct {
	CheckInterval      string `yaml:"check_interval"`
	MaxMatchesPerCheck int    `yaml:"max_matches_per_check"`
}

type TeamConfig struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
	Game string `yaml:"game"`
}

func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.Name,
	)
}
