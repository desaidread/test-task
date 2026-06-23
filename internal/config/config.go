package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConn string
	Port   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("unable to load .env: %v", err)
	}

	var Conf Config

	Conf.DBConn = os.Getenv("SQL_CONN")
	if Conf.DBConn == "" {
		return nil, fmt.Errorf("SQL_CONN is not set")
	}
	Conf.Port = os.Getenv("PORT")
	if Conf.Port == "" {
		Conf.Port = "8080"
	}

	return &Conf, nil

}
