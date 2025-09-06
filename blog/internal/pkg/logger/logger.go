package logger

import (
	"log/slog"
	"os"
)

// 全局 Logger 实例
var Log *slog.Logger

// InitLogger 初始化日志记录器
// inProduction: 是否为生产环境
func InitLogger(inProduction bool) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		// AddSource: true, // 如果需要记录文件名和行号，可以开启这个选项
		Level: slog.LevelDebug, // 默认级别为 Debug，所有级别的日志都会被记录
	}

	if inProduction {
		// 生产环境使用 JSON 格式，方便机器解析
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// 开发环境使用 Text 格式，方便人类阅读
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)

	// 将我们的 logger 设置为标准库的默认 logger
	slog.SetDefault(Log)
}
