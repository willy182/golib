package golib

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {

	t.Parallel()

	var l Level

	t.Run("TraceLevel String", func(t *testing.T) {
		l = TraceLevel
		assert.Equal(t, "trace", l.String())
	})

	t.Run("DebugLevel String", func(t *testing.T) {
		l = DebugLevel
		assert.Equal(t, "debug", l.String())
	})

	t.Run("InfoLevel String", func(t *testing.T) {
		l = InfoLevel
		assert.Equal(t, "info", l.String())
	})

	t.Run("WarnLevel String", func(t *testing.T) {
		l = WarnLevel
		assert.Equal(t, "warning", l.String())
	})

	t.Run("ErrorLevel String", func(t *testing.T) {
		l = ErrorLevel
		assert.Equal(t, "error", l.String())
	})

	t.Run("FatalLevel String", func(t *testing.T) {
		l = FatalLevel
		assert.Equal(t, "fatal", l.String())
	})

	t.Run("PanicLevel String", func(t *testing.T) {
		l = PanicLevel
		assert.Equal(t, "panic", l.String())
	})

	t.Run("Unknown String", func(t *testing.T) {
		l = 9999
		assert.Equal(t, "unknown", l.String())
	})
}

func TestInitLogger(t *testing.T) {
	t.Run("InitLogger", func(t *testing.T) {
		InitLogger("test", "test", "test")
		assert.Equal(t, "test", TOPIC)
		assert.Equal(t, "test", LogTag)
		assert.Equal(t, "test", Env)
	})
}

func TestLogContext(t *testing.T) {
	c := "test"
	s := "test"
	customTags := make([]map[string]interface{}, 1)
	customTag := make(map[string]interface{})

	t.Run("SUCCESS LOGCONTEXT", func(t *testing.T) {
		customTag["test"] = "test"
		customTags = append(customTags, customTag)
		assert.NotNil(t, LogContext(c, s, customTags))
	})
}

func TestLog(t *testing.T) {

	t.Parallel()

	var l Level

	m := "test"
	c := "test"
	s := "test"
	customTag := make(map[string]interface{})
	customTag["test"] = "test"

	t.Run("DebugLevel LOG", func(t *testing.T) {
		l = DebugLevel
		Log(l, m, c, s, customTag)
	})

	t.Run("InfoLevel LOG", func(t *testing.T) {
		l = InfoLevel
		Log(l, m, c, s, customTag)
	})

	t.Run("WarnLevel LOG", func(t *testing.T) {
		l = WarnLevel
		Log(l, m, c, s, customTag)
	})

	t.Run("ErrorLevel LOG", func(t *testing.T) {
		l = ErrorLevel
		Log(l, m, c, s, customTag)
	})

	t.Run("PanicLevel LOG", func(t *testing.T) {
		l = PanicLevel
		Log(l, m, c, s, customTag)
	})

}

func TestLogError(t *testing.T) {
	err := errors.New("test")
	ctx := "test"

	t.Parallel()

	t.Run("LOG ERROR", func(t *testing.T) {
		msg := "test"
		LogError(err, ctx, msg)
	})
}

func Test_newFileResultLogger(t *testing.T) {

	t.Run("SUCCESS newFileResultLogger", func(t *testing.T) {
		re := regexp.MustCompile(`^(.*golib)`)
		cwd, _ := os.Getwd()
		rootPath := string(re.Find([]byte(cwd)))
		res := newFileResultLogger(rootPath)
		assert.Equal(t, rootPath, res.baseDir)
	})

	t.Run("ERROR newFileResultLogger", func(t *testing.T) {
		re := regexp.MustCompile(`^(.*gox)`)
		cwd, _ := os.Getwd()
		rootPath := string(re.Find([]byte(cwd)))
		res := newFileResultLogger(rootPath)
		assert.Equal(t, rootPath, res.baseDir)
	})
}

func TestFileResultLogger_LastError(t *testing.T) {
	t.Run("NIL LastError", func(t *testing.T) {
		f := &FileResultLogger{}
		assert.Nil(t, f.LastError())
	})

	t.Run("NIL LastError", func(t *testing.T) {
		f := &FileResultLogger{
			lastError: errors.New("error"),
		}
		assert.Error(t, f.LastError())
	})
}

func TestFileResultLogger_GetFileName(t *testing.T) {
	t.Run("SUCCESS GetFileName", func(t *testing.T) {
		f := &FileResultLogger{}
		assert.NotEqual(t, "", f.GetFileName("test"))
	})
}

func TestFileResultLogger_Get(t *testing.T) {
	t.Run("ERROR Get", func(t *testing.T) {
		f := &FileResultLogger{}
		assert.Equal(t, "", f.Get(""))
	})
}

func TestFileResultLogger_Store(t *testing.T) {
	t.Run("SUCCESS Store", func(t *testing.T) {
		re := regexp.MustCompile(`^(.*golib)`)
		cwd, _ := os.Getwd()
		rootPath := string(re.Find([]byte(cwd)))
		f := &FileResultLogger{}
		f.baseDir = rootPath
		s := f.Store("go.sum", []byte("test"))
		assert.NotEqual(t, "", s)
	})

	t.Run("ERROR Store", func(t *testing.T) {
		f := &FileResultLogger{}
		s := f.Store("test", []byte("test"))
		assert.NotEqual(t, "", s)
	})
}

func TestFileResultLogger_RequestResponse(t *testing.T) {
	t.Run("ERROR RequestResponse", func(t *testing.T) {
		f := &FileResultLogger{}
		s := f.RequestResponse("test", "test")
		assert.Equal(t, "", s)
	})
}

func TestGetResultLogger(t *testing.T) {
	t.Run("NULL BASE GetResultLogger", func(t *testing.T) {
		assert.NotNil(t, GetResultLogger())
	})
}

func TestStoreRequestResponse(t *testing.T) {
	t.Run("", func(t *testing.T) {
		s := StoreRequestResponse("200", []byte("test"), []byte("test"))
		fmt.Println(s)
	})
}
