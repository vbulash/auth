package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Env Структура для распаковки .env
type Env struct {
	Host       string
	Database   string
	Username   string
	Password   string
	Port       int
	ServerHost string
	ServerPort int
}

// Config Экземпляр конфигурации .env
var Config *Env

// LoadConfig Загрузка конфигурации из .env
func LoadConfig() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	return &Env{
		Host:       os.Getenv("DB_HOST"),
		Database:   os.Getenv("DB_DATABASE"),
		Username:   os.Getenv("DB_USERNAME"),
		Password:   os.Getenv("DB_PASSWORD"),
		Port:       port,
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: serverPort,
	}, nil
}
