package config

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap/zapcore"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Mailer   MailerConfig
	Logger   LoggerConfig
	Cors     CorsConfig
}

type DatabaseConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type ServerConfig struct {
	Mode string
	Port string
}

type JWTConfig struct {
	SigningAlg string
	Key        string
}

type MailerConfig struct {
	Dir  string
	Pass string
	Host string
	Port string
}

type LoggerConfig struct {
	Filepath string
	Filename string
	Level    zapcore.Level
	Tz       string
}

type CorsConfig struct {
	AllowedCors string
}

var Cfg *Config

func LoadConfig() error {
	Cfg = &Config{}

	loadServerConfig(&Cfg.Server)

	loadLoggerConfig(&Cfg.Logger)

	if err := loadJWTConfig(&Cfg.JWT); err != nil {
		return fmt.Errorf("LoadConfig: %w", err)
	}

	if err := loadDatabaseConfig(&Cfg.Database); err != nil {
		return fmt.Errorf("LoadConfig: %w", err)
	}

	if err := loadMailerConfig(&Cfg.Mailer); err != nil {
		return fmt.Errorf("LoadConfig: %w", err)
	}

	if err := loadCorsConfig(&Cfg.Cors); err != nil {
		return fmt.Errorf("LoadConfig: %w", err)
	}

	return nil
}

func loadServerConfig(serverConfig *ServerConfig) {
	serverConfig.Mode = GetEnvOrDefault("GIN_MODE", "debug")
	serverConfig.Port = GetEnvOrDefault("PORT", "8000")

	log.Println("Server - Configuration done.")
}

func loadLoggerConfig(loggerConfig *LoggerConfig) {
	loggerConfig.Filepath = GetEnvOrDefault("FILEPATH", "/app/logs/")
	loggerConfig.Filename = GetEnvOrDefault("FILENAME", "template-logs")
	loggerConfig.Tz = GetEnvOrDefault("TZ", "America/Santiago")

	if Cfg.Server.Mode == "debug" {
		loggerConfig.Level = zapcore.DebugLevel
	} else {
		loggerConfig.Level = zapcore.InfoLevel
	}

	log.Println("Logger - Configuration done.")
}

func loadDatabaseConfig(dbConfig *DatabaseConfig) error {
	databaseVars := []string{"DB_USER_POSTGRES", "DB_PASS_POSTGRES", "DB_HOST_POSTGRES", "DB_PORT_POSTGRES", "DB_NAME_POSTGRES"}
	if err := CheckMissingEnv(databaseVars); err != nil {
		return err
	}

	dbConfig.User = os.Getenv("DB_USER_POSTGRES")
	dbConfig.Pass = os.Getenv("DB_PASS_POSTGRES")
	dbConfig.Host = os.Getenv("DB_HOST_POSTGRES")
	dbConfig.Port = os.Getenv("DB_PORT_POSTGRES")
	dbConfig.Name = os.Getenv("DB_NAME_POSTGRES")

	log.Println("Database - Configuration done.")
	return nil
}

func loadJWTConfig(jwtConfig *JWTConfig) error {
	jwtVars := []string{"JWT_KEY"}
	if err := CheckMissingEnv(jwtVars); err != nil {
		return nil
	}

	jwtKey := os.Getenv("JWT_KEY")

	jwtConfig.Key = jwtKey
	jwtConfig.SigningAlg = "HS256"

	log.Println("JWT - Configuration done.")

	return nil
}

func loadMailerConfig(mailerConfig *MailerConfig) error {
	mailerVars := []string{"EMAIL_DIR", "EMAIL_PASS", "EMAIL_HOST", "EMAIL_PORT"}
	if err := CheckMissingEnv(mailerVars); err != nil {
		return err
	}

	mailerConfig.Dir = os.Getenv("EMAIL_DIR")
	mailerConfig.Pass = os.Getenv("EMAIL_PASS")
	mailerConfig.Host = os.Getenv("EMAIL_HOST")
	mailerConfig.Port = os.Getenv("EMAIL_PORT")

	log.Println("Mailer - Configuration done.")

	return nil
}

func loadCorsConfig(corsConfig *CorsConfig) error {
	cors, set := os.LookupEnv("CORS_URLS")
	if !set {
		return fmt.Errorf("missing cors url")
	}

	corsConfig.AllowedCors = cors

	log.Println("Cors: Configuration done.")

	return nil
}
