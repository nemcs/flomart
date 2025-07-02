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
