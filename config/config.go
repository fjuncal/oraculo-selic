package config

import "os"

type Config struct {
	DatabaseURL    string
	QueueURL       string
	QueueManager   string
	QueueName      string
	ConnectionName string
	Channel        string
	UserID         string
	Password       string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		QueueURL:       os.Getenv("QUEUE_URL"),
		QueueManager:   os.Getenv("QUEUE_MANAGER"),
		QueueName:      os.Getenv("QUEUE_NAME"),
		ConnectionName: os.Getenv("CONNECTION_NAME"),
		Channel:        os.Getenv("CHANNEL"),
		UserID:         os.Getenv("USER_ID"),
		Password:       os.Getenv("PASSWORD"),
	}
}
