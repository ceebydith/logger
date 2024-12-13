package logger

// MustHandler struct implements a no-op LogHandler for default cases.
type MustHandler struct{}

// Log is a no-op logging method.
func (l *MustHandler) Log(v ...any) {}

// Logf is a no-op formatted logging method.
func (l *MustHandler) Logf(format string, v ...any) {}

// LogDefer is a no-op deferred logging method.
func (l *MustHandler) LogDefer(err *error, v ...any) func() { return func() {} }

// LogfDefer is a no-op deferred formatted logging method.
func (l *MustHandler) LogfDefer(err *error, format string, v ...any) func() { return func() {} }

// NewMustHandler creates a LogHandler using the provided creation function or returns a no-op handler if none is provided.
func NewMustHandler(create func(LogHandler) LogHandler, logger ...LogHandler) LogHandler {
	if len(logger) == 0 || logger[0] == nil {
		return &MustHandler{}
	}
	if create == nil {
		return logger[0]
	}
	return create(logger[0])
}
