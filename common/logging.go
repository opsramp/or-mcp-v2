package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DEBUG level for detailed troubleshooting
	DEBUG LogLevel = iota
	// INFO level for general operational information
	INFO
	// WARN level for potentially harmful situations
	WARN
	// ERROR level for error events
	ERROR
	// FATAL level for very severe error events that will lead to application termination
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// CustomLogger is a custom logger that writes to both stdout and a file
type CustomLogger struct {
	level      LogLevel
	stdLogger  *log.Logger
	fileLogger *log.Logger
	mu         sync.Mutex
	file       *os.File
}

var (
	// Global logger instance
	globalLogger *CustomLogger
	// Once ensures the global logger is initialized only once
	once sync.Once
)

// InitLogger initializes the global logger
func InitLogger(level LogLevel, logDir, logFileName string) (*CustomLogger, error) {
	var err error
	once.Do(func() {
		globalLogger, err = newLogger(level, logDir, logFileName)
	})
	return globalLogger, err
}

// GetLogger returns the global logger instance
func GetLogger() *CustomLogger {
	if globalLogger == nil {
		// If not initialized, create a default logger to stdout
		globalLogger = &CustomLogger{
			level:     INFO,
			stdLogger: log.New(os.Stdout, "", log.LstdFlags),
		}
	}
	return globalLogger
}

// newLogger creates a new logger instance
func newLogger(level LogLevel, logDir, logFileName string) (*CustomLogger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create or open log file
	logFilePath := filepath.Join(logDir, logFileName)
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, file)

	return &CustomLogger{
		level:      level,
		stdLogger:  log.New(os.Stdout, "", log.LstdFlags),
		fileLogger: log.New(multiWriter, "", log.LstdFlags),
		file:       file,
	}, nil
}

// Close closes the log file
func (l *CustomLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetLevel sets the log level
func (l *CustomLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log logs a message with the given level
func (l *CustomLogger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// Extract just the filename
	file = filepath.Base(file)

	// Format the message
	msg := fmt.Sprintf(format, args...)
	logMsg := fmt.Sprintf("[%s] [%s:%d] %s", level.String(), file, line, msg)

	// Log to file (and stdout via multiwriter)
	if l.fileLogger != nil {
		l.fileLogger.Println(logMsg)
	} else {
		// Fallback to stdout only if file logger is not available
		l.stdLogger.Println(logMsg)
	}
}

// Debug logs a debug message
func (l *CustomLogger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *CustomLogger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *CustomLogger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *CustomLogger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits the application
func (l *CustomLogger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// FormatRequest formats an HTTP request for logging
func FormatRequest(r *os.File) string {
	if r == nil {
		return ""
	}

	var request []string
	request = append(request, fmt.Sprintf("Request: %s", time.Now().Format(time.RFC3339)))

	return strings.Join(request, "\n")
}

// LogRequest logs an HTTP request
func (l *CustomLogger) LogRequest(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// LogResponse logs an HTTP response
func (l *CustomLogger) LogResponse(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// LogToolExecution logs a tool execution
func (l *CustomLogger) LogToolExecution(toolName, action string, args map[string]interface{}) {
	// Convert args to a string
	argsStr := ""
	for k, v := range args {
		argsStr += fmt.Sprintf("%s=%v ", k, v)
	}
	l.Info("Tool Execution: %s, Action: %s, Args: %s", toolName, action, argsStr)
}

// LogToolResult logs a tool execution result
func (l *CustomLogger) LogToolResult(toolName, action string, result interface{}, err error) {
	if err != nil {
		l.Error("Tool Result: %s, Action: %s, Error: %v", toolName, action, err)
	} else {
		l.Info("Tool Result: %s, Action: %s, Result: %v", toolName, action, result)
	}
}
