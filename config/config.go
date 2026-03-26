package config

import (
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

func LoadConfig() {
	Cfg = &Config{}
	loadServerConfig(&Cfg.Server)
	loadJWTConfig(&Cfg.JWT)
	loadDatabaseConfig(&Cfg.Database)
	loadMailerConfig(&Cfg.Mailer)
	loadLoggerConfig(&Cfg.Logger)
	loadCorsConfig(&Cfg.Cors)
}

func loadServerConfig(serverConfig *ServerConfig) {
	serverConfig.Mode = GetEnvOrDefault("GIN_MODE", "debug")
	serverConfig.Port = GetEnvOrDefault("PORT", "8000")

	log.Println("Server: Configuration done.")
}

func loadDatabaseConfig(dbConfig *DatabaseConfig) {
	databaseVars := []string{"DB_USER_POSTGRES", "DB_PASS_POSTGRES", "DB_HOST_POSTGRES", "DB_PORT_POSTGRES", "DB_NAME_POSTGRES"}

	CheckMissingEnv(databaseVars)

	dbConfig.User = os.Getenv("DB_USER_POSTGRES")
	dbConfig.Pass = os.Getenv("DB_PASS_POSTGRES")
	dbConfig.Host = os.Getenv("DB_HOST_POSTGRES")
	dbConfig.Port = os.Getenv("DB_PORT_POSTGRES")
	dbConfig.Name = os.Getenv("DB_NAME_POSTGRES")

	log.Println("Database: Configuration done.")
}

func loadJWTConfig(jwtConfig *JWTConfig) {
	jwtVars := []string{"JWT_KEY"}
	CheckMissingEnv(jwtVars)

	jwtKey := os.Getenv("JWT_KEY")

	jwtConfig.Key = jwtKey
	jwtConfig.SigningAlg = "HS256"

	log.Println("JWT: Configuration done.")
}

func loadMailerConfig(mailerConfig *MailerConfig) {
	mailerVars := []string{"EMAIL_DIR", "EMAIL_PASS", "EMAIL_HOST", "EMAIL_PORT"}

	CheckMissingEnv(mailerVars)

	mailerConfig.Dir = os.Getenv("EMAIL_DIR")
	mailerConfig.Pass = os.Getenv("EMAIL_PASS")
	mailerConfig.Host = os.Getenv("EMAIL_HOST")
	mailerConfig.Port = os.Getenv("EMAIL_PORT")

	log.Println("Mailer: Configuration done.")
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

	log.Println("Logger: Configuration done.")
}

func loadCorsConfig(corsConfig *CorsConfig) {
	cors, set := os.LookupEnv("CORS_URLS")
	if !set {
		log.Fatal("Missing cors variable.")
	}
	corsConfig.AllowedCors = cors

	log.Println("Cors: Configuration done.")
}
