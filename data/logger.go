package data

type SmsboxLogger interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
	Printf(format string, v ...interface{})
	LevelPrintLn(level string, v ...interface{})
}
