package logging

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LLvlAccess    = "ACCESS_LOG"
	LLvlService   = "SERVICE_LOG"
	LLvlMysql     = "MYSQL_LOG"
	LLvlHTTP      = "HTTP_LOG"
	LLvlInfo      = "INFO_LOG"
	LLvlPubSub    = "PUBSUB_LOG"
	LLvlRedis     = "REDIS_LOG"
	LLvlGRPC      = "GRPC_LOG"
	LLvlException = "EXCEPTION_LOG"
)

const (
	LLvlDevelopment = iota + 1
	LLvlProduction
	LStdDir  = "log"        // initial value for log's directory
	LStdFile = "2006-01-02" // log's default filename
)

var (
	instance      *zap.Logger
	sugarInstance *zap.SugaredLogger
	once          sync.Once
	initWithDir   bool // flag identify package declare with dir or not
)

func Logging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := uuid.NewV4().String()
			c.Set("request_id", reqID)

			defer func(now time.Time) {

				// flush cache
				// cache.Delete(reqID)

				message := LLvlAccess
				fields := []zap.Field{
					zap.String("at", now.Format(time.RFC3339)),
					zap.String("method", c.Request().Method),
					zap.String("uri", c.Request().URL.String()),
					zap.String("ip", c.RealIP()),
					zap.String("host", c.Request().Host),
					zap.String("user_agent", c.Request().UserAgent()),
					zap.Int("code", c.Response().Status),
					zap.Any("query_param", c.QueryParams()),
				}
				WithRequestID(reqID).Info(message, fields...)

			}(time.Now())

			return next(c)
		}
	}
}

// Options represent option to custom-zap logger
//
// Level set log's level logger, either development or production
// Time set log's time location being used, default is "Asia/Jakarta".
// Use according to Time Zone database, such as "America/New_York".
// Output file is another output file. If you want logger to write log
// to multiple file, add other source here. add "stdout" for console log.
type Options struct {
	Level int
	// Time       *time.Location
	OutputFile []string
}

// newLogger return new custom zap-logger
// set default logger: logs to os.stdout, production level,
// with Asia/Jakarta time. default log filename yyyy-mm-dd.log
func newLoggerWithDir(dir string, prefix string, opt *Options) *zap.Logger {
	var filename = "stdout"

	if opt == nil {
		opt = &Options{}
	}

	if opt.Level < 1 {
		opt.Level = LLvlProduction
	}

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("create folder in: %s\n", dir)
			if err = os.MkdirAll(dir, os.ModePerm|os.ModeAppend); err != nil {
				panic(fmt.Sprintf("[log] failed to create directory: %v", err))
			}
		}

		logFile := fmt.Sprintf("%s.%s", time.Now().Format(LStdFile), LStdDir)
		if prefix != "" {
			filename = fmt.Sprintf("%s/%s-%s", dir, prefix, logFile)
		}
		filename = fmt.Sprintf("%s/%s", dir, logFile)
	}

	logger, err := opt.newConfig(filename).Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

// newConfig set config for custom-zap logger
// set log's file to logFile
// set log's time with timeLocation
func (opt *Options) newConfig(logFile string) (cfg zap.Config) {
	if opt.Level > LLvlDevelopment {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.StacktraceKey = ""
	}

	cfg.OutputPaths = []string{logFile}
	cfg.ErrorOutputPaths = []string{logFile}

	if len(opt.OutputFile) > 0 {
		for _, out := range opt.OutputFile {
			cfg.OutputPaths = append(cfg.OutputPaths, out)
			cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, out)
		}
	}

	return cfg
}

func InitLoggerWithDir(dir string, prefix string, opt *Options) {
	initWithDir = true
	instance = newLoggerWithDir(dir, prefix, opt)
}

func getLogger() *zap.Logger {
	once.Do(func() {
		if !initWithDir {
			instance = newLogger()
		}
	})

	return instance
}

func GetSugaredLogger() *zap.SugaredLogger {
	if sugarInstance == nil {
		sugarInstance = getLogger().Sugar()
	}
	return sugarInstance
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	logger, err := loggerConfig.Build()
	if err != nil {
		Errorf("failed to create new logger with error: %s", err)
		panic(err)
	}
	return logger
}

func Debug(message string, fields ...zap.Field) {
	getLogger().Debug(message, fields...)
}

func Debugf(template string, args ...interface{}) {
	GetSugaredLogger().Debugf(template, args...)
}

func Panic(message string, fields ...zap.Field) {
	getLogger().Error(message, fields...)
	panic(message)
}

func Error(err error, fields ...zap.Field) {
	getLogger().Error(err.Error(), fields...)
}

func Errorf(template string, args ...interface{}) {
	GetSugaredLogger().Errorf(template, args...)
}

func Fatal(err error, fields ...zap.Field) {
	getLogger().Fatal(err.Error(), fields...)
}

func Fatalf(template string, args ...interface{}) {
	GetSugaredLogger().Fatalf(template, args...)
}

func Info(message string, fields ...zap.Field) {
	getLogger().Info(message, fields...)
}

func Infof(template string, args ...interface{}) {
	GetSugaredLogger().Infof(template, args...)
}

func Warn(message string, fields ...zap.Field) {
	getLogger().Warn(message, fields...)
}

func Warnf(template string, args ...interface{}) {
	GetSugaredLogger().Warnf(template, args...)
}

func AddHook(hook func(zapcore.Entry) error) {
	instance = getLogger().WithOptions(zap.Hooks(hook))
	sugarInstance = instance.Sugar()
}

func WithRequestID(reqID string) *zap.Logger {
	return getLogger().With(
		zap.String("request_id", reqID),
	)
}

func WithRequest(r *http.Request) *zap.Logger {
	return getLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}

func SugaredWithRequest(r *http.Request) *zap.SugaredLogger {
	return GetSugaredLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}
