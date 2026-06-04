package game

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Helper function for log file management.
func compressAndRemove(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return ErrConcat("failed to open old file", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return ErrConcat("failed to create new file", err)
	}
	defer dstFile.Close()

	gw := gzip.NewWriter(dstFile)
	if _, err := io.Copy(gw, srcFile); err != nil {
		return ErrConcat("failed to compress old file", err)
	}
	defer gw.Close()
	srcFile.Close()

	if err := os.Remove(srcPath); err != nil {
		return ErrConcat("failed to remove old file", err)
	}

	return nil
}

// ErrConcat concatenates an error and a string message into a new error.
func ErrConcat(msg string, errs error) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("%s: %w", msg, errs)
	}
	// Formats as: "failed to initialize... (main.go:42): original error"
	return fmt.Errorf("%s (%s:%d): %w", msg, file, line, errs)
}

// InitLogging initializes the main logger.
// Should be called once at startup.
func InitLogging(logDir, latestLogName string) error {
	latestPath := logDir + "/" + latestLogName + ".log"

	// Ignore the mkdir error intentionally
	os.MkdirAll(logDir, 0755)

	// Check if latest exists, if yes move it
	if _, err := os.Stat(latestPath); err == nil {
		archiveName := fmt.Sprintf("%s/log-%s.log.gz", logDir, time.Now().Format("2006-01-02-150405"))

		if err := compressAndRemove(latestPath, archiveName); err != nil {
			fmt.Printf("Warning: failed to archive old log: %v\n", err)
		}
	}

	// No defer file.Close() because this stays open for the entire time the process is running
	file, err := os.OpenFile(latestPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	handler := slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))
	return nil
}

// LInfo is shorthand for
//
//	slog.Info(msg, args...)
func LInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}

// LWarn is shorthand for
//
//	slog.Warn(msg, args...)
func LWarn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// LError is shorthand for
//
//	slog.Error(msg, args...)
func LError(err any, args ...any) {
	switch t := err.(type) {
	case error:
		slog.Error(t.Error(), args...)
	case string:
		slog.Error(t, args...)
	default:
		slog.Error(fmt.Sprint(t), args...)
	}
}

// LPanic is shorthand for
//
//	slog.Error(err, args...)
//	panic(err)
func LPanic(err any, args ...any) {
	LError(err)
	panic(err)
}
