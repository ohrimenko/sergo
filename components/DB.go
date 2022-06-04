package components

import (
	"errors"
	"log"

	"github.com/ohrimenko/sergo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbStruct struct {
	connect *gorm.DB
	err     error
	valid   bool
}

var dbApp dbStruct

func (n *dbStruct) db() (*gorm.DB, error) {
	if !n.valid {
		if config.Env("DB_DRIVER") == "mysql" {
			dsn := config.Env("DB_USERNAME") + ":" + config.Env("DB_PASSWORD") + "@tcp(" + config.Env("DB_HOST") + ":" + config.Env("DB_PORT") + ")/" + config.Env("DB_DATABASE") + "?charset=" + config.Env("DB_CHARSET") + "&parseTime=True&loc=Local"

			n.connect, n.err = gorm.Open(mysql.New(mysql.Config{
				DSN:                       dsn,   // data source name
				DefaultStringSize:         256,   // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
			}), &gorm.Config{})

			if n.err == nil {
				sqlDB, err := n.connect.DB()

				if err == nil {
					// SetMaxIdleConns устанавливает максимальное количество соединений в пуле незанятых соединений.
					sqlDB.SetMaxIdleConns(100)

					// SetMaxOpenConns устанавливает максимальное количество открытых подключений к базе данных.
					sqlDB.SetMaxOpenConns(0)

					// SetConnMaxLifetime устанавливает максимальное количество времени, в течение которого соединение может быть повторно использовано.
					sqlDB.SetConnMaxLifetime(0)

					sqlDB.SetConnMaxIdleTime(0)
				}
			}

			n.valid = true
		} else {
			n.err = errors.New("Error Select Driver DB")
		}
	}

	return n.connect, n.err
}

func DB() (*gorm.DB, error) {
	return dbApp.db()
}

func CloseDB() {
	DB, err := DB()

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, errDB := DB.DB()

	if errDB == nil {
		sqlDB.Close()
	}
}
