package config

import (
	"fmt"
	"os"

	"github.com/cybre/salesforge-assignment/internal/database"
)

type Config struct {
	Port     string
	Database database.Config
}

func LoadConfig() Config {
	return Config{
		Port: GetEnv("PORT", "3000"),
		Database: database.Config{
			Host:     MustGetEnv("DATABASE_HOST"),
			Port:     MustGetEnv("DATABASE_PORT"),
			Name:     MustGetEnv("DATABASE_NAME"),
			User:     MustGetEnv("DATABASE_USER"),
			Password: ReadSecret(MustGetEnv("DATABASE_PASSWORD_FILE")),
		},
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func MustGetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("%s must be set", key))
}

func ReadSecret(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 100)
	n, err := file.Read(buf)
	if err != nil {
		panic(err)
	}

	return string(buf[:n])
}
