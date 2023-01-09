package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type tag string // a non-primitive to use as context keys
type internal string

// tags is the key for a []tag, used for contextual logging
const tagsKey = internal("tags")

// DEBUG is controlled by the LOG_LEVEL environment variable.
var DEBUG bool
var mu sync.Mutex

func SetDebug(val bool) {
	mu.Lock()
	DEBUG = val
	mu.Unlock()
}

// Logger is the main type for the logger package.
type Logger struct {
	tags tagList
}

// tagList contains metadata for the logger instance.
type tagList map[tag]interface{}

type message struct {
	Level     string  `json:"level"`
	EventTime string  `json:"event_time"`
	Message   string  `json:"message"`
	Trace     string  `json:"trace,omitempty"`
	Tags      tagList `json:"tags,omitempty"`
}

func init() {
	DEBUG = os.Getenv("DEBUG") == "TRUE"
}

// New instantiates and returns a Logger object
func New() *Logger {
	return &Logger{
		tags: make(tagList),
	}
}

func stripPathPrefix(s string) string {
	i := strings.LastIndex(s, "/")
	if i > 0 {
		s = s[i+1:]
	}
	return s
}

func stack() string {
	names := make([]string, 0, 10)
	pc := make([]uintptr, 100)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()
		if frame.Function == "main.index" {
			break
		}
		if strings.Contains(frame.File, "asm_amd64.s") {
			break
		}

		path := fmt.Sprintf("%s:%s:%d", stripPathPrefix(frame.File), stripPathPrefix(frame.Function), frame.Line)
		names = append(names, path)
		if !more {
			break
		}
	}
	return strings.Join(names, ", ")
}

func (l *Logger) With(ctx context.Context, k string, v interface{}) context.Context {
	// Add tag to logger.
	l.tags[tag(k)] = v
	ctx = context.WithValue(ctx, k, v)

	// Add value to context.
	tags := make([]tag, 0, len(l.tags))
	for k, v := range l.tags {
		tags = append(tags, k)
		ctx = context.WithValue(ctx, k, v)
	}
	return context.WithValue(ctx, tagsKey, tags)
}

// FromContext returns a new *Logger, automatically
// adding tags from ctx if ctx contains a
// logger.Tags key with a value of []string.
func FromContext(ctx context.Context) *Logger {
	log := New()
	tags, ok := ctx.Value(tagsKey).([]tag)
	if !ok {
		return log
	}
	for _, t := range tags {
		log.tags[t] = ctx.Value(t)
	}
	return log
}

func (l *Logger) ValueString(key string) string {
	v, found := l.tags[tag(key)]
	if !found {
		return ""
	}
	return fmt.Sprintf("%s", v)
}

func (l *Logger) log(level, text string) {
	msg := message{
		Level:     level,
		EventTime: time.Now().UTC().Format(time.RFC3339),
		Message:   text,
		Tags:      l.tags,
	}
	if level == "ERROR" {
		msg.Trace = stack()
	}
	b, err := json.Marshal(&msg)

	if err != nil {
		fmt.Printf("logger ERROR: cannot marshal payload: %s", err)
		return
	}

	fmt.Fprintln(os.Stdout, string(b))
}

// Debug prints out a message with DEBUG level.
func (l Logger) Debug(message string) {
	if !DEBUG {
		return
	}
	l.log("DEBUG", message)
}

// Debugf prints out a message with DEBUG level.
func (l Logger) Debugf(message string, args ...interface{}) {
	l.Debug(fmt.Sprintf(message, args...))
}

// Info prints out a message with INFO level.
func (l Logger) Info(message string) {
	l.log("INFO", message)
}

// Infof prints out a message with INFO level.
func (l Logger) Infof(message string, args ...interface{}) {
	l.Info(fmt.Sprintf(message, args...))
}

// Error prints out a message with ERROR level.
func (l Logger) Error(message string) {
	l.log("ERROR", message)
}

// Errorf prints out a message with INFO level.
func (l Logger) Errorf(message string, args ...interface{}) {
	l.Error(fmt.Sprintf(message, args...))
}
