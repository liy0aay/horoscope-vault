package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBURL             string
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	ServerPort        string
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}
	config := Config{
		DBURL:             os.Getenv("DATABASE_URL"),
		JWTPrivateKeyPath: os.Getenv("JWT_PRIVATE_KEY_PATH"),
		JWTPublicKeyPath:  os.Getenv("JWT_PUBLIC_KEY_PATH"),
		ServerPort:        os.Getenv("PORT"),
	}

	return config, nil
}
