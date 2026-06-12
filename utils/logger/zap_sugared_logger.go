package logger

import "go.uber.org/zap"

type zapSugaredLogger struct {
	log *zap.SugaredLogger
}

func InitZapSugaredLogger() error {
	if Log != nil {
		return nil
	}

	opts := []zap.Option{}

	l, err := zap.NewDevelopment(opts...)
	if err != nil {
		return err
	}

	sl := l.Sugar()

	setLogger(&zapSugaredLogger{log: sl})

	return nil
}

func (l *zapSugaredLogger) Debugf(template string, args ...any) {
	l.log.Debugf(template, args...)
}

func (l *zapSugaredLogger) Debugw(msg string, keysAndValues ...any) {
	l.log.Debugw(msg, keysAndValues...)
}

func (l *zapSugaredLogger) Infof(template string, args ...any) {
	l.log.Infof(template, args...)
}

func (l *zapSugaredLogger) Infow(msg string, keysAndValues ...any) {
	l.log.Infow(msg, keysAndValues...)
}

func (l *zapSugaredLogger) Warnf(template string, args ...any) {
	l.log.Warnf(template, args...)
}

func (l *zapSugaredLogger) Warnw(msg string, keysAndValues ...any) {
	l.log.Warnw(msg, keysAndValues...)
}

func (l *zapSugaredLogger) Errorf(template string, args ...any) {
	l.log.Errorf(template, args...)
}

func (l *zapSugaredLogger) Errorw(msg string, keysAndValues ...any) {
	l.log.Errorw(msg, keysAndValues...)
}

func (l *zapSugaredLogger) Sync() error {
	return l.log.Sync()
}
