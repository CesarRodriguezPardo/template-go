package config

import (
	"citiaps/golang-backend-template/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
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
	UsePostgres bool
	UseMongo    bool
	Postgres    PostgresConfig
	Mongo       MongoConfig
}

type PostgresConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
	URI  string
}

type MongoConfig struct {
	User       string
	Pass       string
	Host       string
	Port       string
	Name       string
	URI        string
	Collection string
}

type ServerConfig struct {
	Mode string
	Port string
}

type JWTConfig struct {
	FlagAlgType bool // true simetrico - falso asimetrico
	SigningAlg  string
	Key         string
	PrivKeyPath string
	PubKeyPath  string
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
	CorsUrl string
}

var Cfg *Config

func LoadConfig() {
	Cfg = &Config{}

	// Gin config
	Cfg.Server.Mode = utils.GetEnvWithDefault("GIN_MODE", "debug")
	gin.SetMode(Cfg.Server.Mode)

	// port
	Cfg.Server.Port = utils.GetEnvWithDefault("PORT", "8080")

	//jwt config
	loadJWTConfig(&Cfg.JWT)

	// db config
	loadPostgresConfig(&Cfg.Database)
	loadMongoConfig(&Cfg.Database)

	// mailer config
	loadMailerConfig(&Cfg.Mailer)

	// logger config
	loadLoggerConfig(&Cfg.Logger)

	// cors config
	loadCorsConfig(&Cfg.Cors)

	// que al menos una de las 2 esten configuradas, sino no sirve
	if !Cfg.Database.UsePostgres && !Cfg.Database.UseMongo {
		log.Fatalf("CONFIG DATABASE: Al menos una base de datos (postgres o mongo) debe estar configurada.")
	}
}

func loadPostgresConfig(dbConfig *DatabaseConfig) {
	user := os.Getenv("DB_USER_POSTGRES")
	pass := os.Getenv("DB_PASS_POSTGRES")
	host := os.Getenv("DB_HOST_POSTGRES")
	port := os.Getenv("DB_PORT_POSTGRES")
	name := os.Getenv("DB_NAME_POSTGRES")

	if user != "" && pass != "" && host != "" && port != "" && name != "" {
		dbConfig.UsePostgres = true
		dbConfig.Postgres = PostgresConfig{
			User: user,
			Pass: pass,
			Host: host,
			Port: port,
			Name: name,
			URI:  fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name),
		}
	}
}

func loadMongoConfig(dbConfig *DatabaseConfig) {
	user := os.Getenv("DB_USER_MONGO")
	pass := os.Getenv("DB_PASS_MONGO")
	host := os.Getenv("DB_HOST_MONGO")
	port := os.Getenv("DB_PORT_MONGO")
	name := os.Getenv("DB_NAME_MONGO")

	if user != "" && pass != "" && host != "" && port != "" && name != "" {
		dbConfig.UseMongo = true
		dbConfig.Mongo = MongoConfig{
			User: user,
			Pass: pass,
			Host: host,
			Port: port,
			Name: name,
			URI:  fmt.Sprintf("mongodb://%s:%s/%s", host, port, name),
		}
	}
}

// por default implementa un jwt con algoritmo de firma simetrico y una clave default.
func loadJWTConfig(jwtConfig *JWTConfig) {
	jwtAlgorithm := utils.GetEnvWithDefault("SIGNING_ALG", "HS256")
	jwtConfig.SigningAlg = jwtAlgorithm

	// caso es simetrico
	symmetricAlgs := []string{"HS256", "HS384", "HS52"}
	if slices.Contains(symmetricAlgs, jwtAlgorithm) {
		jwtConfig.Key = utils.GetEnvWithDefault("JWT_KEY", "clave_super_privada")
		jwtConfig.FlagAlgType = true

		log.Println("CONFIG JWT: Algoritmo simetrico detectado y seteado en config.")
	}

	// caso es asimetrico
	asymmetricAlgs := []string{"RS256"} // por añadir mas?
	if slices.Contains(asymmetricAlgs, jwtAlgorithm) {

		privKeyPath := utils.GetEnvWithDefault("PRIVATE_KEY_PATH", "/app/keys/priv.pem")
		pubkeyPath := utils.GetEnvWithDefault("PUBLIC_KEY_PATH", "/app/keys/pub.pem")

		if !utils.FileExists(privKeyPath) && !utils.FileExists(pubkeyPath) {
			log.Println("CONFIG JWT: Ambas llaves, privada y publica, deben existir.", errors.New("No existen las llaves."))
		}

		jwtConfig.PrivKeyPath = privKeyPath
		jwtConfig.PubKeyPath = pubkeyPath
		jwtConfig.FlagAlgType = false

		log.Println("CONFIG JWT: Algoritmo asimetrico detectado y seteado en config.")
	}
}

func loadMailerConfig(mailerConfig *MailerConfig) {
	mailerVars := []string{"EMAIL_DIR", "EMAIL_PASS", "EMAIL_HOST", "EMAIL_PORT"}
	missing := ""

	for _, v := range mailerVars {
		_, set := os.LookupEnv(v)
		if !set {
			missing = missing + v
		}
	}

	// por ahora solamente se dejaran como variables criticas, puesto que no he preguntado si se pueden dejar como default
	if len(missing) != 0 {
		log.Fatal("CONFIG: Variables criticas no definidas: "+missing, errors.New("Variables faltantes"))
	}

	mailerConfig.Dir = os.Getenv("EMAIL_DIR")
	mailerConfig.Pass = os.Getenv("EMAIL_PASS")
	mailerConfig.Host = os.Getenv("EMAIL_HOST")
	mailerConfig.Port = os.Getenv("EMAIL_PORT")
}

func loadLoggerConfig(loggerConfig *LoggerConfig) {
	loggerConfig.Filepath = utils.GetEnvWithDefault("FILEPATH", "/app/logs/")
	loggerConfig.Filename = utils.GetEnvWithDefault("FILENAME", "template-logs")

	loggerConfig.Tz = utils.GetEnvWithDefault("TZ", "America/Santiago")

	if Cfg.Server.Mode == "debug" {
		loggerConfig.Level = zapcore.DebugLevel
	}
	loggerConfig.Level = zapcore.InfoLevel
}

func loadCorsConfig(corsConfig *CorsConfig) {
	cors, set := os.LookupEnv("CORS_URLS")
	if !set {
		log.Fatal("CONFIG: CORS no definidos.")
	}

	corsConfig.CorsUrl = cors
}
