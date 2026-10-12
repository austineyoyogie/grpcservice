package database

import (
	"gocgrpcarchitecture/internal/pkg/configs"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var lc = configs.LoadConfig()

func CDriver() *gorm.DB {
	DC, err := gorm.Open(mysql.Open(lc.DBC.Username+":"+lc.DBC.Password+
		"@tcp("+lc.DBC.Hostname+":"+lc.DBC.Port+")/"+lc.DBC.Database+
		"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logFatal(err)
		return nil
	}
	return DC
}

func logFatal(err error) {
	if err != nil {
		log.Fatal("\nError database connection mode.")
	}
}

/*
import (
	"database/sql"
	"gocgrpcarchitecture/backend/configs"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var lc = configs.LoadConfig()

func CDriver() *sql.DB {
	DC, err := sql.Open("mysql", lc.DBC.Username+":"+lc.DBC.Password+
		"@tcp("+lc.DBC.Hostname+":"+lc.DBC.Port+")/"+lc.DBC.Database+
		"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logFatal(err)
		return nil
	}
	return DC
}

func logFatal(err error) {
	if err != nil {
		log.Fatal("\nError database connection mode.")
	}
}
*/
