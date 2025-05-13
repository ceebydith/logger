package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Logger interface provides basic logging methods.
// This is useful for holding your logger in a variable globally.
//
// Example:
//
// var log logger.Logger
type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Defer(err *error, v ...any) func()
	Deferf(err *error, format string, v ...any) func()
}

// LogHandler interface provides logging methods for a specific type.
// This allows your type to log messages anywhere in your code.
// There is no difference between Logger and LogHandler, but LogHandler is more convenient to use with your type.
//
// Example:
//
//	type YourType struct {
//	    logger.LogHandler
//	}
type LogHandler interface {
	Log(v ...any)
	Logf(format string, v ...any)
	LogDefer(err *error, v ...any) func()
	LogfDefer(err *error, format string, v ...any) func()
}

// Default is a basic logger with added functionalities.
type Default struct {
	logger *log.Logger
}

// Helper function for internal printing.
func (l *Default) print(v ...any) {
	lines := strings.ReplaceAll(strings.TrimSpace(fmt.Sprintln(v...)), "\n", "\n\t")
	l.logger.Println(lines)
}

// Print logs the message with formatting.
func (l *Default) Print(v ...any) {
	l.print(v...)
}

// Printf logs the formatted message.
func (l *Default) Printf(format string, v ...any) {
	l.print(fmt.Sprintf(format, v...))
}

// Defer logs the message and defers the final status.
func (l *Default) Defer(err *error, v ...any) func() {
	l.print(append(v, "...")...)
	return func() {
		if err == nil {
			l.print(append(v, "DONE")...)
		} else if *err == nil {
			l.print(append(v, "OK")...)
		} else {
			l.print(append(v, "ERROR\nError:", *err)...)
		}
	}
}

// Deferf logs the formatted message and defers the final status.
func (l *Default) Deferf(err *error, format string, v ...any) func() {
	return l.Defer(err, fmt.Sprintf(format, v...))
}

// Log logs the message.
func (l *Default) Log(v ...any) {
	l.print(v...)
}

// Logf logs the formatted message.
func (l *Default) Logf(format string, v ...any) {
	l.print(fmt.Sprintf(format, v...))
}

// LogDefer logs the message and defers the final status.
func (l *Default) LogDefer(err *error, v ...any) func() {
	return l.Defer(err, v...)
}

// LogfDefer logs the formatted message and defers the final status.
func (l *Default) LogfDefer(err *error, format string, v ...any) func() {
	return l.Defer(err, fmt.Sprintf(format, v...))
}

// Writer returns the output destination for the logger.
func (l *Default) Writer() io.Writer {
	return l.logger.Writer()
}

// New returns a new instance of Default.
func New(w ...io.Writer) *Default {
	if len(w) == 0 {
		return &Default{
			logger: log.New(os.Stderr, "", log.LstdFlags),
		}
	}
	return &Default{
		logger: log.New(io.MultiWriter(append([]io.Writer{os.Stderr}, w...)...), "", log.LstdFlags),
	}
}
