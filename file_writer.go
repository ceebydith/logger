package logger

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ceebydith/curly"
)

// fileWriter struct combines the Buffer interface with file writing capabilities.
type fileWriter struct {
	Buffer
	fileformat  string
	currentfile string
	file        *os.File
}

// onstop closes the file when the buffer processing stops.
func (w *fileWriter) onstop() {
	if w.file != nil {
		w.file.Close()
	}
}

// ondata handles data received in the buffer and writes it to the file.
func (w *fileWriter) ondata(buffer []byte) {
	if err := w.open(); err != nil {
		os.Stderr.WriteString("logger.FileWriter error: " + err.Error() + "\n")
		return
	}
	if _, err := w.file.Write(buffer); err != nil {
		os.Stderr.WriteString("logger.FileWriter error: " + err.Error() + "\n")
	}
}

// open opens a new file if the current one has changed, and creates necessary directories.
func (w *fileWriter) open() error {
	currentfile, err := curly.Format(w.fileformat)
	if err != nil {
		return err
	}
	if currentfile == w.currentfile {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(currentfile), 0744); err != nil {
		return err
	}

	if w.file != nil {
		w.file.Close()
	}

	w.file, err = os.OpenFile(currentfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.currentfile = currentfile
	return nil
}

// FileWriter initializes and returns a new fileWriter instance with the given parameters.
func FileWriter(ctx context.Context, fileformat string, buffer int) *fileWriter {
	w := &fileWriter{
		fileformat: fileformat,
	}
	w.Buffer = NewBuffer(ctx, buffer, w.ondata, nil, w.onstop)
	return w
}
