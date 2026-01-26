package logging


func WithFields(args ...any) []any {
	return args
}

func LogError(err error, msg string, args ...any) []any {
	return append([]any{"error", err, "message", msg}, args ...)
}

func LogDuration(duration float64, args ...any) []any {
	return append([]any{"duration_ms", duration}, args ...)
}

func LogSize(size int64, args ...any) []any {
	return append([]any{"size_bytes", size}, args ...)
}