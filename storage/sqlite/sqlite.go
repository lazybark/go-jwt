package sqlite

import (
	"log"
	"os"

	"github.com/lazybark/go-jwt/storage"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

type SQLite struct {
	db *gorm.DB
}

func NewSQLiteStorage(dbName string) (storage.Storage, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{Logger: gLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gLogger.Config{LogLevel: gLogger.Silent},
	)})
	return SQLite{db: db}, err
}

func (s SQLite) UserAdd(storage.User) (int, error) {
	return 0, nil
}

func (s SQLite) UserGetData(login string) (storage.User, error) {
	return storage.User{}, nil
}

func (s SQLite) Flush() error {
	return nil
}

func (s SQLite) Init() error {
	return nil
}
