package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBC struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	SSLMode  string
}

type JWT struct {
	KEY string
}

type Mail struct {
	Server   string
	Port     string
	Email    string
	Password string
}

type TOKEN struct {
	ACCESS_TOKEN_MAX_AGE     string
	ACCESS_TOKEN_EXPIRED_IN  string
	REFRESH_TOKEN_EXPIRED_IN string
	REFRESH_TOKEN_MAX_AGE    string
}

type Config struct {
	DBC   DBC
	JWT   JWT
	Mail  Mail
	TOKEN TOKEN
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
		return Config{}
	}
	return Config{
		DBC{
			os.Getenv("DB_USER"),
			os.Getenv("DB_PWD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL"),
		},
		JWT{
			os.Getenv("SECRET_TOKEN_KEY"),
		},
		Mail{
			os.Getenv("EMAIL_HOST"),
			os.Getenv("EMAIL_PORT"),
			os.Getenv("EMAIL_HOST_USER"),
			os.Getenv("EMAIL_HOST_PASSWORD"),
		},
		TOKEN{
			os.Getenv("ACCESS_TOKEN_MAX_AGE"),
			os.Getenv("ACCESS_TOKEN_EXPIRED_IN"),
			os.Getenv("REFRESH_TOKEN_EXPIRED_IN"),
			os.Getenv("REFRESH_TOKEN_MAX_AGE"),
		},
	}
}
