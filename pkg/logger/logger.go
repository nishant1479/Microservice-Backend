package logger

import (
	"html/template"
	"os"

	"github.com/nishant1479/Microservice-Backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/tools/go/cfg"
)

type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string,args ...interface{})
	Info(args ...interface{})
	Infof(template string,args ...interface{})
	Warn(args ...interface{})
	Warnf(template string,args ...interface{})
	Error(args ...interface{})
	Errorf(template string,args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string,args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string,args ...interface{})
	Printf(template string,args ...interface{})	
}

type apiLogger struct{
	cfg			*config.Config
	sugarLogger	*zap.SugaredLogger
}


// App logger Constructor
func NewApiLogger(cfg *config.Config) *apiLogger{
	return &apiLogger{cfg: cfg}
}


var loggerLevelMap = map[string]zapcore.Level{
	"debug":	zapcore.DebugLevel,
	"info":		zapcore.InfoLevel,
	"warn":		zapcore.WarnLevel,
	"error":	zapcore.ErrorLevel,
	"dpanic":	zapcore.DPanicLevel,
	"panic":	zap.PanicLevel,
	"fatal":	zapcore.FatalLevel,
}

func (l *apiLogger) getLoggerLevel(cfg *config.Config) zapcore.Level{
	level,exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}

func (l *apiLogger) InitLogger(){
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.Development{
		encodercfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey	= "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.Newcore(encoder,logWriter,zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core,zap.AddCaller(),zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err:= l.sugarLogger.Sync();err != nil {
		l.sugarLogger.Error(err)
	}
}