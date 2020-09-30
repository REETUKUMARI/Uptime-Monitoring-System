package models

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var err error

//ConnectDataBase connects to database
func ConnectDataBase() {

	DB, err = gorm.Open("mysql", "reetu:Password@(localhost:3306)/test")

	if err != nil {
		fmt.Printf("failed to connect to database!")
		os.Exit(1)
	}

	DB.AutoMigrate(&Pingdom{})
}
