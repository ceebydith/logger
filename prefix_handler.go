package logger

// PrefixHandler struct combines a logging prefix with a LogHandler.
type PrefixHandler struct {
	prefix  string
	handler LogHandler
}

// Log logs a message with the prefixed log handler.
func (l *PrefixHandler) Log(v ...any) {
	l.handler.Log(append([]any{l.prefix + ":"}, v...)...)
}

// Logf logs a formatted message with the prefixed log handler.
func (l *PrefixHandler) Logf(format string, v ...any) {
	l.handler.Logf(l.prefix+": "+format, v...)
}

// LogDefer logs a message with deferred error handling using the prefixed log handler.
func (l *PrefixHandler) LogDefer(err *error, v ...any) func() {
	return l.handler.LogDefer(err, append([]any{l.prefix + ":"}, v...)...)
}

// LogfDefer logs a formatted message with deferred error handling using the prefixed log handler.
func (l *PrefixHandler) LogfDefer(err *error, format string, v ...any) func() {
	return l.handler.LogfDefer(err, l.prefix+": "+format, v...)
}

// NewPrefixHandler initializes and returns a new PrefixHandler instance with the given prefix and LogHandler.
func NewPrefixHandler(prefix string, handler LogHandler) *PrefixHandler {
	return &PrefixHandler{
		prefix:  prefix,
		handler: handler,
	}
}
