package logger

import (
	"CesarRodriguezPardo/template-go/config"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	filePerm         = 0644        // user permit
	writeAndReadFile = os.O_RDWR   // escribir y leer un archivo
	createFile       = os.O_CREATE // crear un archivo
	addToFile        = os.O_APPEND // añadir al final del archivo sin sobreescribirlo

	defaultSize     = 1
	defaultBackups  = 0
	defaultLifeTime = 0
)

var logger *zap.Logger
var santiago *time.Location

func InitLogger() {

	filepath := config.Cfg.Logger.Filepath
	filename := config.Cfg.Logger.Filename
	logLevel := config.Cfg.Logger.Level
	timeZone := config.Cfg.Logger.Tz

	logFilePath := filepath + filename + ".log"

	consoleWriter := zapcore.AddSync(os.Stdout)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    defaultSize,
		MaxBackups: defaultBackups,
		MaxAge:     defaultLifeTime,
	})

	var err error
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		santiago = time.UTC
	} else {
		santiago = loc
	}

	encoderConfig := &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     SantiagoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	syncer := zapcore.NewMultiWriteSyncer(consoleWriter, fileWriter)
	encoder := zapcore.NewConsoleEncoder(*encoderConfig)

	core := zapcore.NewCore(encoder, syncer, logLevel)

	logger = zap.New(core)

	if err != nil {
		logger.Error("ERROR: No se cargo la zona horaria especificada. Usando UTC.", zap.String("timezone", timeZone), zap.Error(err))
	} else {
		logger.Info("Zona horaria cargada correctamente", zap.String("timezone", timeZone))
	}

	defer logger.Sync()
}

func Info(message string, fields ...zap.Field)  { logger.Info(message, fields...) }
func Debug(message string, fields ...zap.Field) { logger.Debug(message, fields...) }
func Warn(message string, fields ...zap.Field)  { logger.Warn(message, fields...) }
func Error(message string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	logger.Error(message, fields...)
}
func Fatal(message string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	logger.Fatal(message, fields...)
}

func SantiagoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.In(santiago).Format("2006-01-02 15:04:05.000 MST"))
}
