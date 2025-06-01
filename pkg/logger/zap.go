package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DefaultConfig = Config{
		Level:          -1,
		TimeLayout:     "2006-01-02 15:04",
		Caller:         true,
		Trace:          false,
		FilePath:       "sample.log",
		FileMaxBackups: 30,
		FileMaxAge:     10,
		FileMaxSize:    500,
		FileCompress:   true,
	}
	Logger *zap.Logger
	Enable bool
)

type Config struct {
	// debug=-1,info=0,warn=1,error=2,DPanic=3,panic=4,fatal=5
	Enable     bool   `json:"enable"`
	Level      int    `json:"level"`
	TimeLayout string `json:"time_layout"`
	Caller     bool   `json:"caller"`
	Trace      bool   `json:"trace"`
	FilePath   string `json:"file_path"`
	// MB
	FileMaxSize int `json:"file_max_size"`
	// maximum number of old log files to retain
	FileMaxBackups int `json:"file_max_backups"`
	// day
	FileMaxAge   int  `json:"file_max_age"`
	FileCompress bool `json:"file_compress"`
}

func Init(cfg Config) {
	if cfg.Enable {
		Logger = CreateLogger(cfg)
		Enable = true
	} else {
		Enable = false
	}
}

func CreateLogger(cfg Config) *zap.Logger {

	logLevel := zapcore.Level(cfg.Level)
	// rotate file
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.FileMaxSize,
		MaxBackups: cfg.FileMaxBackups,
		MaxAge:     cfg.FileMaxAge,
		Compress:   cfg.FileCompress,
		LocalTime:  true,
	})
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.TimeKey = "date"
	if cfg.TimeLayout != "" {
		encoderConf.EncodeTime = zapcore.TimeEncoderOfLayout(cfg.TimeLayout)
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConf),
		w,
		logLevel,
	)
	opts := []zap.Option{
		zap.AddCallerSkip(1),
	}
	if cfg.Caller {
		opts = append(opts, zap.AddCaller())
	}
	if cfg.Trace {
		opts = append(opts, zap.AddStacktrace(logLevel))
	}
	return zap.New(core, opts...)
}
func Error(msg ...any) {
	if Enable {
		Logger.Error(fmt.Sprint(msg...))
	}
}
func Errorf(format string, msg ...any) {
	if Enable {
		Logger.Error(fmt.Sprintf(format, msg...))
	}
}
func Warn(msg ...any) {
	if Enable {
		Logger.Warn(fmt.Sprint(msg...))
	}
}
func Warnf(format string, msg ...any) {
	if Enable {
		Logger.Warn(fmt.Sprintf(format, msg...))
	}
}
func Info(msg ...any) {
	if Enable {
		Logger.Info(fmt.Sprint(msg...))
	}
}
func Infof(format string, msg ...any) {
	if Enable {
		Logger.Info(fmt.Sprintf(format, msg...))
	}
}
func Debug(msg ...any) {
	if Enable {
		Logger.Debug(fmt.Sprint(msg...))
	}
}
func Debugf(format string, msg ...any) {
	if Enable {
		Logger.Debug(fmt.Sprintf(format, msg...))
	}
}
func Panic(msg ...any) {
	if Enable {
		Logger.Panic(fmt.Sprint(msg...))
	}
}
func Panicf(format string, msg ...any) {
	if Enable {
		Logger.Panic(fmt.Sprintf(format, msg...))
	}
}
func Fatal(msg ...any) {
	if Enable {
		Logger.Fatal(fmt.Sprint(msg...))
	}
}
func Fatalf(format string, msg ...any) {
	if Enable {
		Logger.Fatal(fmt.Sprintf(format, msg...))
	}
}
func DPanic(msg ...any) {
	if Enable {
		Logger.DPanic(fmt.Sprint(msg...))
	}
}
func DPanicf(format string, msg ...any) {
	if Enable {
		Logger.DPanic(fmt.Sprintf(format, msg...))
	}
}
