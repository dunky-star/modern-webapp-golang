package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxLogSize  = 5 * 1024 * 1024     // 5MB
	maxLogAge   = 14 * 24 * time.Hour // 2 weeks
	logDir      = "log"
	logFileName = "access.log"
)

// rotatingLogWriter handles log file rotation based on size and age
type rotatingLogWriter struct {
	mu          sync.Mutex
	currentFile *os.File
	currentSize int64
	createdAt   time.Time
	logPath     string
	baseWriter  io.Writer // For console output in dev mode
}

// newRotatingLogWriter creates a new rotating log writer
func newRotatingLogWriter(alsoWriteToConsole bool) (*rotatingLogWriter, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, logFileName)

	// Open or create the log file
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Get file info for size
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	// Determine file creation time for age tracking
	// Use modification time as proxy for file age (on Unix, creation time isn't reliable)
	fileAge := info.ModTime()
	// If file is empty (new file), use current time
	if info.Size() == 0 {
		fileAge = time.Now()
	}

	writer := &rotatingLogWriter{
		currentFile: file,
		currentSize: info.Size(),
		createdAt:   fileAge,
		logPath:     logPath,
	}

	if alsoWriteToConsole {
		writer.baseWriter = os.Stdout
	}

	// If existing file is already older than max age, rotate immediately
	if time.Since(fileAge) >= maxLogAge && info.Size() > 0 {
		if err := writer.rotate(); err != nil {
			return nil, fmt.Errorf("failed to rotate old log file: %w", err)
		}
	}

	return writer, nil
}

// Write implements io.Writer interface
func (r *rotatingLogWriter) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if rotation is needed
	if r.shouldRotate() {
		if err := r.rotate(); err != nil {
			// If rotation fails, still try to write to current file
			fmt.Fprintf(os.Stderr, "Failed to rotate log: %v\n", err)
		}
	}

	// Write to file
	n, err = r.currentFile.Write(p)
	if err != nil {
		return n, err
	}
	r.currentSize += int64(n)

	// Also write to console if configured
	if r.baseWriter != nil {
		r.baseWriter.Write(p)
	}

	return n, nil
}

// shouldRotate checks if log rotation is needed
func (r *rotatingLogWriter) shouldRotate() bool {
	// Rotate if file size exceeds max size
	if r.currentSize >= maxLogSize {
		return true
	}

	// Rotate if file age exceeds max age
	if time.Since(r.createdAt) >= maxLogAge {
		return true
	}

	return false
}

// rotate performs log file rotation
func (r *rotatingLogWriter) rotate() error {
	// Close current file
	if err := r.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	// Generate rotated filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	rotatedPath := fmt.Sprintf("%s.%s", r.logPath, timestamp)

	// Rename current file to rotated filename
	if err := os.Rename(r.logPath, rotatedPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// Create new log file
	file, err := os.OpenFile(r.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat new log file: %w", err)
	}

	// Update writer state with new file
	r.currentFile = file
	r.currentSize = info.Size()
	r.createdAt = time.Now() // New file, so creation time is now

	return nil
}

// Close closes the log file
func (r *rotatingLogWriter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentFile != nil {
		return r.currentFile.Close()
	}
	return nil
}
