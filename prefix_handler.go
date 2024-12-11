package logger

// prefixHandler struct combines a logging prefix with a LogHandler.
type prefixHandler struct {
	prefix  string
	handler LogHandler
}

// Log logs a message with the prefixed log handler.
func (l *prefixHandler) Log(v ...any) {
	l.handler.Log(append([]any{l.prefix + ":"}, v...)...)
}

// Logf logs a formatted message with the prefixed log handler.
func (l *prefixHandler) Logf(format string, v ...any) {
	l.handler.Logf(l.prefix+": "+format, v...)
}

// LogDefer logs a message with deferred error handling using the prefixed log handler.
func (l *prefixHandler) LogDefer(err *error, v ...any) func() {
	return l.handler.LogDefer(err, append([]any{l.prefix + ":"}, v...)...)
}

// LogfDefer logs a formatted message with deferred error handling using the prefixed log handler.
func (l *prefixHandler) LogfDefer(err *error, format string, v ...any) func() {
	return l.handler.LogfDefer(err, l.prefix+": "+format, v...)
}

// PrefixHandler initializes and returns a new prefixHandler instance with the given prefix and LogHandler.
func PrefixHandler(prefix string, handler LogHandler) *prefixHandler {
	return &prefixHandler{
		prefix:  prefix,
		handler: handler,
	}
}
