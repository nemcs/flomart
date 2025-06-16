package logger

import (
	"log/slog"
	"os"
)

// TODO Добавь logger.Init() и вызови в main — так можно будет переключать JSON/Text по окружению.
const (
	FieldErr = "error"
	FieldSQL = "sqlError"
	Handler  = "handler"
)

var Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
}))

//TODO Сделать корректную глубину вызова в стеке для Source через Stack Frames API
//func Info(msg string, args ...any) {
//	Log.Info(msg, args...)
//}
//func Warn(msg string, args ...any) {
//	Log.Warn(msg, args...)
//}
//func Error(msg string, args ...any) {
//	Log.Error(msg, args...)
//}
