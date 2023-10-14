package initialization

import (
	"fmt"
	"go-admin/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const (
	Level = "info"
)

func init() {
	fmt.Println("开始获取日志")
	InitLogger()
}
func InitLogger() {

	hookInfo := lumberjack.Logger{
		Filename:   global.StdOut, // 日志文件路径
		MaxSize:    128,           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,            // 日志文件最多保存多少个备份
		MaxAge:     14,             // 文件最多保存多少天
		Compress:   true,          // 是否压缩
	}
	hookError :=lumberjack.Logger{
		Filename:   global.StdError, // 日志文件路径
		MaxSize:    128,             // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,              // 日志文件最多保存多少个备份
		MaxAge:     14,               // 文件最多保存多少天
		Compress:   true,            // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:    func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	var level zapcore.Level
	switch Level {
	case "info":
		level = zap.InfoLevel
	case "debug":
		level = zap.DebugLevel
	default:
		level = zap.InfoLevel
	}
	// 设置日志级别
	atom := zap.NewAtomicLevelAt(level)

	config := zap.Config{
		Level:         atom,                                            // 日志级别
		Development:   true,                                            // 开发模式，堆栈跟踪
		Encoding:      "json",                                          // 输出格式 console 或 json
		EncoderConfig: encoderConfig,                                   // 编码器配置
		InitialFields: map[string]interface{}{"serviceName": "dynamic-store-api"}, // 初始化字段，如：添加一个服务器名称

		OutputPaths:      []string{"stdout", global.StdOut}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr", global.StdError},
	}


	encoder :=zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(&hookInfo), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(&hookError), errorLevel),
	)
	//config = config
	_,_ = config.Build()
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("api", "dynamic-weapp-api"))
	logger := zap.New(core, caller, development, filed)

	zap.ReplaceGlobals(logger)
	zap.S().Infof("日志开始记录！")
}
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}


