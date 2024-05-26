package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    string `gorm:"primaryKey"` // Defines ID as the primary key.
	Name  string `gorm:"not null"`   // Username must be unique and not null.
	Email string `gorm:"unique"`
}

type DBConfig struct {
	DB      *gorm.DB
	Name    string
	Server  string
	Port    int
	DSN     string
	Secrets struct {
		username string
		password string
	}
}

type Application struct {
	FiberApp *fiber.App
	Config   *DBConfig
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	wallet := &Application{
		FiberApp: fiber.New(),
		Config: &DBConfig{
			Name:   "test",
			Server: "localhost",
			Port:   3306,
			Secrets: struct {
				username string
				password string
			}{
				username: "root",
				password: "BlackNigga", // Replace with your actual password
			},
		},
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	wallet.Config.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", wallet.Config.Secrets.username, wallet.Config.Secrets.password, wallet.Config.Server, wallet.Config.Port, wallet.Config.Name)

	var err error
	wallet.Config.DB, err = gorm.Open(mysql.Open(wallet.Config.DSN), &gorm.Config{})
	if err != nil {
		wallet.ErrorLog.Println(err)
		return
	} else {
		wallet.InfoLog.Printf("Connected to MariaDB to Database: %s Port: %d", wallet.Config.Name, wallet.Config.Port)
	}

	if err := wallet.Config.DB.AutoMigrate(&User{}); err != nil {
		wallet.ErrorLog.Println("Error during AutoMigrate:", err)
	} else {
		wallet.InfoLog.Println("Database migration completed successfully.")
	}

	wallet.FiberApp.Listen(":3000")
}
