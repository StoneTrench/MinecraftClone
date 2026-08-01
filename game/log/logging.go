package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/StoneTrench/go-mc-clone/metadata"
)

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
func Init(logDir, latestLogName string) error {
	latestPath := logDir + "/" + latestLogName + ".log"

	// Ignore the mkdir error intentionally
	os.MkdirAll(logDir, 0755)

	var old_log_error error = nil
	// Check if latest exists, if yes move it
	if _, err := os.Stat(latestPath); err == nil {
		archiveName := fmt.Sprintf("%s/log-%s.log.gz", logDir, time.Now().Format("2006-01-02-150405"))

		if err := compressAndRemove(latestPath, archiveName); err != nil {
			old_log_error = fmt.Errorf("failed to archive old log, %v", err)
		}
	}

	// No defer file.Close() because this stays open for the entire time the process is running
	file, err := os.OpenFile(latestPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	handler := NewBracketHandler(io.MultiWriter(os.Stdout, file))
	slog.SetDefault(slog.New(handler))

	if old_log_error != nil {
		slog.Warn(old_log_error.Error())
	}

	return nil
}

// Info is shorthand for
//
//	slog.Info(msg, args...)
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Infof is shorthand for
//
//	slog.Info(fmt.Sprintf(format, s...))
func Infof(format string, s ...any) {
	slog.Info(fmt.Sprintf(format, s...))
}

// Warn is shorthand for
//
//	slog.Warn(msg, args...)
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error is shorthand for
//
//	slog.Error(msg, args...)
func Error(err any, args ...any) {
	switch t := err.(type) {
	case error:
		slog.Error(t.Error(), args...)
	case string:
		slog.Error(t, args...)
	default:
		slog.Error(fmt.Sprint(t), args...)
	}
}

// Panic is shorthand for
//
//	slog.Error(err, args...)
//	panic(err)
func Panic(err any, args ...any) {
	Error(err)
	panic(err)
}

// Assert is LogError, but with a condition.
func Assert(cond bool, err any, args ...any) {
	if !cond && metadata.IsInDebugMode() {
		switch t := err.(type) {
		case error:
			slog.Error(t.Error(), args...)
		case string:
			slog.Error(t, args...)
		default:
			slog.Error(fmt.Sprint(t), args...)
		}
	}
}
