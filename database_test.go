package golib

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDBLogFormatter_Format(t *testing.T) {
	t.Run("SUCCES FORMAT", func(t *testing.T) {
		data := make(log.Fields, 0)
		data["test"] = "test"
		formatter := &DBLogFormatter{}
		entry := &log.Entry{
			Message: "test",
			Data:    data,
		}
		b, err := formatter.Format(entry)
		assert.NoError(t, err)
		assert.Equal(t, "test\n", string(b))
	})
}

func TestInitDB(t *testing.T) {
	t.Run("PANIC InitDB", func(t *testing.T) {
		os.Setenv("DEBUG", "1")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		InitDB()
	})
}

func TestGetWriteDB(t *testing.T) {
	t.Parallel()

	t.Run("NIL & PANIC GetWriteDB", func(t *testing.T) {
		os.Setenv("DBW_HOST", "")
		os.Setenv("DBW_USER", "")
		os.Setenv("DBW_PASS", "")
		os.Setenv("DBW_NAME", "")
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		GetWriteDB()
	})

	t.Run("NOT NIL GetWriteDB", func(t *testing.T) {
		sql, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open("postgres", sql)
		dbWrite = gormDB
		assert.Equal(t, gormDB, GetWriteDB())
	})
}

func TestGetReadDB(t *testing.T) {
	t.Run("NIL & PANIC GetReadDB", func(t *testing.T) {
		os.Setenv("DBW_HOST", "")
		os.Setenv("DBW_USER", "")
		os.Setenv("DBW_PASS", "")
		os.Setenv("DBW_NAME", "")
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		GetReadDB()
	})

	t.Run("NOT NIL GetReadDB", func(t *testing.T) {
		sql, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open("postgres", sql)
		dbRead = gormDB
		assert.Equal(t, gormDB, GetReadDB())
	})
}

func TestCloseDb(t *testing.T) {

	t.Run("SUCCESS CLOSE DB", func(t *testing.T) {
		sqlRead, _, _ := sqlmock.New()
		gormDBRead, _ := gorm.Open("postgres", sqlRead)
		dbRead = gormDBRead

		sqlWrite, _, _ := sqlmock.New()
		gormDBWrite, _ := gorm.Open("postgres", sqlWrite)
		dbWrite = gormDBWrite

		CloseDb()

		assert.Nil(t, dbRead)
		assert.Nil(t, dbWrite)
	})
}
