package data

type SmsboxLogger interface {
	Info(message string, v ...interface{})
	Warn(message string, v ...interface{})
	Error(message string, v ...interface{})
	Debug(message string, v ...interface{})
	Printf(format string, v ...interface{})
	LevelPrintLn(level string, message string, v ...interface{})
}
