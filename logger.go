package golib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

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
	base := os.Getenv("STORAGE_DIR") + "/logs/"
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
