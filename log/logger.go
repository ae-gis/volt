package log

type logging interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Panic(msg string, fields ...interface{})

	// for type fields
	Field(key string, value interface{}) interface{}
}

var logger logging

func Debug(msg string, fields ...interface{}) {
	validator()
	logger.Debug(msg, fields...)
}

// Info add log entry with info level
func Info(msg string, fields ...interface{}) {
	validator()
	logger.Info(msg, fields...)
}

// Warn add log entry with warn level
func Warn(msg string, fields ...interface{}) {
	validator()
	logger.Warn(msg, fields...)
}

// Error add log entry with error level
func Error(msg string, fields ...interface{}) {
	validator()
	logger.Error(msg, fields...)
}

// Fatal add log entry with fatal level
func Fatal(msg string, fields ...interface{}) {
	validator()
	logger.Fatal(msg, fields...)
}

// Panic add log entry with panic level
func Panic(msg string, fields ...interface{}) {
	validator()
	logger.Panic(msg, fields...)
}

func Field(key string, value interface{}) interface{} {
	validator()
	return logger.Field(key, value)
}

func validator() {
	if logger == nil {
		logger = NewZap(ProductionCore())
	}
}
