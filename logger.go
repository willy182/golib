package golib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Level type
type Level uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

var (
	// TOPIC for setting topic of log
	TOPIC string
	// LogTag default log tag
	LogTag string
	// Env default environment
	Env string
)

// InitLogger function init logger
func InitLogger(topic, tag, env string) {
	TOPIC = topic
	LogTag = tag
	Env = env
}

// LogContext function for logging the context of echo
// c string context
// s string scope
func LogContext(c string, s string, customeTags []map[string]interface{}) *log.Entry {
	maps := make(map[string]interface{})
	for _, m := range customeTags {
		for k, v := range m {
			maps[k] = v
		}
	}

	map1 := log.Fields{
		"topic":      TOPIC,
		"context":    c,
		"scope":      s,
		"server_env": Env,
	}

	result := MergeMaps(map1, maps)

	return log.WithFields(result)
}

// Log function for returning entry type
// level log.Level
// message string message of log
// context string context of log
// scope string scope of log
func Log(level Level, message string, context string, scope string, customeTags ...map[string]interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		entry := LogContext(context, scope, customeTags)
		switch level {
		case DebugLevel:
			entry.Debug(message)
		case InfoLevel:
			entry.Info(message)
		case WarnLevel:
			entry.Warn(message)
		case ErrorLevel:
			entry.Error(message)
		case FatalLevel:
			entry.Fatal(message)
		case PanicLevel:
			entry.Panic(message)
		}
	}()
}

// LogError logging error
func LogError(err error, context string, messageData interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		entry := log.WithFields(log.Fields{
			"topic":      TOPIC,
			"context":    context,
			"error":      err,
			"server_env": Env,
		})

		jsonStr, _ := json.Marshal(messageData)
		entry.Error(string(jsonStr))
	}()
}

// ResultLogger result logger interface
type ResultLogger interface {
	Store(c string, d []byte) string
	Get(p string) string
	RequestResponse(c string, dt string) string
	LastError() error
}

// FileResultLogger file based storage
type FileResultLogger struct {
	baseDir   string
	lastError error
}

// newFileResultLogger private function for creating log file
// base string directory
func newFileResultLogger(base string) *FileResultLogger {
	this := new(FileResultLogger)
	if err := this.createOrIgnore(base); err != nil {
		this.lastError = err
		return this
	}
	this.baseDir = base
	return this
}

// LastError function for getting last error
func (flo *FileResultLogger) LastError() error {
	return flo.lastError
}

// createOrIgnore function for creating or ignoring error
// p string
func (flo *FileResultLogger) createOrIgnore(p string) error {
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(p, 0775); err != nil {
				flo.lastError = err
				return err
			}
		} else {
			flo.lastError = err
			return err
		}
	}

	flo.lastError = nil
	return nil
}

// GetFileName function to get file name
// c string code name
func (flo *FileResultLogger) GetFileName(c string) string {
	t := time.Now()
	u := strconv.Itoa(int(time.Now().Unix()))
	f := fmt.Sprintf("%s_%s.%s", t.Format("150405"), u, c)
	return flo.baseDir + "/" + c + "/" + f
}

// Get  function to get log
// p string file name
func (flo *FileResultLogger) Get(p string) string {
	base := os.Getenv("STORAGE_DIR")
	f := base + "/archive/" + p
	data, err := ioutil.ReadFile(f)
	if err != nil {
		flo.lastError = err
		return ""
	}

	flo.lastError = nil
	return string(data)
}

// Store function to store log
// c string code name
// d []byte json data which is about to stored
func (flo *FileResultLogger) Store(c string, d []byte) string {
	fileName := flo.GetFileName(c)
	if err := flo.createOrIgnore(flo.baseDir + "/" + c); err != nil {
		flo.lastError = err
		return fileName
	}

	flo.lastError = nil
	if err := ioutil.WriteFile(fileName, d, 0775); err != nil {
		flo.lastError = err
	}

	return fileName
}

// RequestResponse function for storing request and response into file
// c string content
// dt string date time
func (flo *FileResultLogger) RequestResponse(c string, dt string) string {
	t := time.Now()

	// set the value of data
	dir := os.Getenv("STORAGE_DIR") + "/logs/"
	filename := fmt.Sprintf("%s%s.%s", dir, t.Format("20060102"), c)
	val := fmt.Sprintf("%s : %s", t.Format("15:04:05"), dt)

	if err := flo.createOrIgnore(dir); err == nil {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
		if err == nil {
			if _, err = f.WriteString(val); err != nil {
				panic(err)
			}
		}

		defer f.Close()
	}

	return ""
}

// GetResultLogger function for getting log result
func GetResultLogger() ResultLogger {
	base := os.Getenv("LOG_DIR")
	if base == "" {
		base = os.Getenv("STORAGE_DIR") + "/logs/"
	}
	return newFileResultLogger(base)
}

// StoreRequestResponse function for storing request and response into file
// code string code of file
// req []byte request json byte data
// res []byte response json byte data
func StoreRequestResponse(code string, req []byte, res []byte) string {
	// set data to save/append into log
	data := "REQUEST: " + string(req[:]) + " RESPONSE: " + string(res[:]) + "\n"

	var fileLogger ResultLogger
	fileLogger = GetResultLogger()

	if fileLogger.LastError() != nil {
		return ""
	}

	return fileLogger.RequestResponse(code, data)
}
