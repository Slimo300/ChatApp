package orm

import (
	"fmt"
	"os"

	"github.com/Slimo300/ChatApp/backend/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

// Setup creates Database object and initializes connection between MySQL database
func Setup() (*Database, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("MYSQLUSERNAME"),
		os.Getenv("MYSQLPASSWORD"), os.Getenv("MYSQLHOST"), "3306", os.Getenv("MYSQLDBNAME"))), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{}, models.Group{}, models.Member{}, models.Message{}, models.Invite{})

	return &Database{DB: db}, nil
}
