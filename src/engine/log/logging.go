package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"time"
)

const LOGGING_DIR = "./logs"

// Helper function for log file management.
func compressAndRemove(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open old file, %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create new file, %w", err)
	}
	defer dstFile.Close()

	gw := gzip.NewWriter(dstFile)
	if _, err := io.Copy(gw, srcFile); err != nil {
		return fmt.Errorf("failed to compress old file, %w", err)
	}
	defer gw.Close()
	srcFile.Close()

	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("failed to remove old file, %w", err)
	}

	return nil
}

// Init initializes the main logger.
// Should be called once at startup.
func Init(noFile bool) (err error) {
	// By the way, these comments are human made and are here cause when I was half asleep debugging the code, I kept thinking parts of this were mistakes lol
	latestPath := path.Join(LOGGING_DIR, "latest.log")

	var file *os.File
	var old_log_error error
	if !noFile {
		// Ignore the mkdir error intentionally
		os.MkdirAll(LOGGING_DIR, 0755)

		// Check if latest exists, if yes move it
		if _, err := os.Stat(latestPath); err == nil {
			archiveName := path.Join(LOGGING_DIR, fmt.Sprintf("log-%s.log.gz", time.Now().Format("2006-01-02-150405")))

			if err := compressAndRemove(latestPath, archiveName); err != nil {
				old_log_error = fmt.Errorf("failed to archive old log, %v", err)
			}
		}

		// No defer file.Close() because this stays open for the entire time the process is running
		file, err = os.OpenFile(latestPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
	}

	// handler := slog.NewTextHandler(io.MultiWriter(os.Stdout, file), nil)

	handler := NewBracketHandler(io.MultiWriter(os.Stdout, file))
	slog.SetDefault(slog.New(handler))

	if old_log_error != nil {
		slog.Warn(old_log_error.Error())
	}

	return nil
}

func Infof(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}
func Warnf(format string, args ...any) {
	slog.Warn(fmt.Errorf(format, args...).Error())
}
func Errorf(format string, args ...any) {
	slog.Error(fmt.Errorf(format, args...).Error())
}
func Info(msg string) {
	slog.Info(msg)
}
func WarnStr(err string) {
	slog.Warn(err)
}
func ErrorStr(err string) {
	slog.Error(err)
}
func Warn(err error) {
	slog.Warn(err.Error())
}
func Error(err error) {
	slog.Error(err.Error())
}
