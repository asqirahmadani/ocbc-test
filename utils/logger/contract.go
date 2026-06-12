package logger

var Log logger

type logger interface {
	Debugf(template string, args ...any)
	Debugw(msg string, keysAndValues ...any)
	Infof(template string, args ...any)
	Infow(msg string, keysAndValues ...any)
	Warnf(template string, args ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorf(template string, args ...any)
	Errorw(msg string, keysAndValues ...any)
	Sync() error
}

func setLogger(logger logger) {
	Log = logger
}