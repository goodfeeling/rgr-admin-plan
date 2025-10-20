package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain/constants"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	Log *zap.Logger
}

var encodeLevelMap = map[string]zapcore.LevelEncoder{
	constants.EncodeLevelCapitalcolorlevelencoder:   zapcore.CapitalColorLevelEncoder,
	constants.EncodeLevelCapitallevelencoder:        zapcore.CapitalLevelEncoder,
	constants.EncodeLevelLowercaselevelencoder:      zapcore.LowercaseLevelEncoder,
	constants.EncodeLevelLowercasecolorlevelencoder: zapcore.LowercaseColorLevelEncoder,
}

var logLevelMap = map[string]zapcore.Level{
	constants.LogLevelDebug:   zapcore.DebugLevel,
	constants.LogLevelInfo:    zapcore.InfoLevel,
	constants.LogLevelWarning: zapcore.WarnLevel,
	constants.LogLevelError:   zapcore.ErrorLevel,
	constants.LogLevelClose:   zapcore.PanicLevel,
	constants.LogLevelFatal:   zapcore.FatalLevel,
}

func NewLogger() (*Logger, error) {

	encodeLevel := zapcore.CapitalLevelEncoder
	if v, ok := encodeLevelMap[getenv("ZAP_ENCODE_LEVEL", "LowercaseColorLevelEncoder")]; ok {
		encodeLevel = v
	}

	logLevel := zapcore.InfoLevel
	if v, ok := logLevelMap[getenv("ZAP_LOG_LEVEL", "info")]; ok {
		logLevel = v
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  getenv("ZAP_STACKTRACE_KEY", "stacktrace"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 构建 WriteSyncer 列表
	var writeSyncers []zapcore.WriteSyncer
	enableConsoleOutput := getenv("ZAP_LOG_IN_CONSOLE", "true") != "false"
	if enableConsoleOutput {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		getEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		zap.NewAtomicLevelAt(logLevel),
	)

	logger := zap.New(core)

	return &Logger{Log: logger}, nil
}

// NewDevelopmentLogger crea un logger para desarrollo con más información de debug
func NewDevelopmentLogger() (*Logger, error) {

	encodeLevel := zapcore.CapitalLevelEncoder
	if v, ok := encodeLevelMap[getenv("ZAP_ENCODE_LEVEL", "LowercaseColorLevelEncoder")]; ok {
		encodeLevel = v
	}
	logLevel := zapcore.InfoLevel
	if v, ok := logLevelMap[getenv("ZAP_LOG_LEVEL", "info")]; ok {
		logLevel = v
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  getenv("ZAP_STACKTRACE_KEY", "stacktrace"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	maxAge, err := strconv.Atoi(getenv("ZAP_ENCODE_LEVEL", "30"))
	if err != nil {
		maxAge = 30
	}

	// 使用 lumberjack 进行日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/app.log", getenv("ZAP_DIRNAME", "logs")),
		MaxSize:    100,    // 每个文件最大100MB
		MaxBackups: 5,      // 最多保留5个备份文件
		MaxAge:     maxAge, // 保留30天
		Compress:   true,   // 压缩旧文件
	}

	// 构建 WriteSyncer 列表
	var writeSyncers []zapcore.WriteSyncer
	writeSyncers = append(writeSyncers, zapcore.AddSync(lumberJackLogger))
	enableConsoleOutput := getenv("ZAP_LOG_IN_CONSOLE", "true") != "false"
	if enableConsoleOutput {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		getEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		zap.NewAtomicLevelAt(logLevel),
	)

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))

	return &Logger{Log: logger}, nil
}
func getenv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
func getEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	if os.Getenv("ZAP_ENCODER") == "json" {
		return zapcore.NewJSONEncoder(cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Log.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Log.Panic(msg, fields...)
}
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Log.Debug(msg, fields...)
}

// SetupGinWithZapLogger configura Gin para usar el logger de Zap
func (l *Logger) SetupGinWithZapLogger() {
	// Configurar Gin para usar el modo release por defecto
	gin.SetMode(gin.ReleaseMode)

	// Crear un writer personalizado que use Zap
	gin.DefaultWriter = &ZapWriter{logger: l.Log}
	gin.DefaultErrorWriter = &ZapErrorWriter{logger: l.Log}
}

// SetupGinWithZapLoggerInDevelopment configura Gin para usar el logger de Zap en modo desarrollo
func (l *Logger) SetupGinWithZapLoggerInDevelopment() {
	// Configurar Gin para usar el modo debug en desarrollo
	gin.SetMode(gin.DebugMode)

	// Crear un writer personalizado que use Zap
	gin.DefaultWriter = &ZapWriter{logger: l.Log}
	gin.DefaultErrorWriter = &ZapErrorWriter{logger: l.Log}
}

// SetupGinWithZapLoggerWithMode configura Gin para usar el logger de Zap con un modo específico
func (l *Logger) SetupGinWithZapLoggerWithMode(mode string) {
	// Configurar Gin para usar el modo especificado
	gin.SetMode(mode)

	// Crear un writer personalizado que use Zap
	gin.DefaultWriter = &ZapWriter{logger: l.Log}
	gin.DefaultErrorWriter = &ZapErrorWriter{logger: l.Log}
}

// ZapWriter implementa io.Writer para usar con Gin
type ZapWriter struct {
	logger *zap.Logger
}

func (w *ZapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info("Gin-log", zap.String("message", string(p)))
	return len(p), nil
}

// ZapErrorWriter implementa io.Writer para errores de Gin
type ZapErrorWriter struct {
	logger *zap.Logger
}

func (w *ZapErrorWriter) Write(p []byte) (n int, err error) {
	w.logger.Error("Gin-error", zap.String("error", string(p)))
	return len(p), nil
}

func (l *Logger) GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		l.Log.Info("HTTP request", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()), zap.Duration("latency", latency), zap.String("client_ip", c.ClientIP()))
	}
}

type GormZapLogger struct {
	zap    *zap.SugaredLogger
	config gormlogger.Config
}

func NewGormLogger(base *zap.Logger) *GormZapLogger {
	sugar := base.Sugar()
	return &GormZapLogger{
		zap: sugar,
		config: gormlogger.Config{
			SlowThreshold:             time.Second, // umbral para destacar consultas lentas
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true, // no loguear "record not found"
			Colorful:                  false,
		},
	}
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newCfg := l.config
	newCfg.LogLevel = level
	return &GormZapLogger{zap: l.zap, config: newCfg}
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Info {
		l.zap.Infof(msg, data...)
	}
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Warn {
		l.zap.Warnf(msg, data...)
	}
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Error &&
		(!l.config.IgnoreRecordNotFoundError || msg != gormlogger.ErrRecordNotFound.Error()) {
		l.zap.Errorf(msg, data...)
	}
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	if err != nil {
		if l.config.IgnoreRecordNotFoundError && errors.Is(err, gormlogger.ErrRecordNotFound) {
			return
		}
		if l.config.LogLevel >= gormlogger.Error {
			sql, rows := fc()
			l.zap.Errorf("Error: %v | %.3fms | rows:%d | %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		return
	}

	if elapsed > l.config.SlowThreshold && l.config.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.zap.Warnf("SLOW ≥ %s | %.3fms | rows:%d | %s", l.config.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
