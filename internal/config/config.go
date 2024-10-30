package config

import (
	"net"
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
	DSN        string
	ServerHost string
	ServerPort int
	Address    string
}

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

	address := net.JoinHostPort(os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	return &Env{
		Host:       os.Getenv("DB_HOST"),
		Database:   os.Getenv("DB_DATABASE"),
		Username:   os.Getenv("DB_USERNAME"),
		Password:   os.Getenv("DB_PASSWORD"),
		Port:       port,
		DSN:        os.Getenv("DB_DSN"),
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: serverPort,
		Address:    address,
	}, nil
}
