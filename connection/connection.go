package connection

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnection(stringConnection string) gorm.DB {
	sqlDB, _ := sql.Open("mysql", stringConnection)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return *gormDB
}
