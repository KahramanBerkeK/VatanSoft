package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Ad    string
	Soyad string
	Email string
	Sifre string
}

var db *gorm.DB

func main() {
	e := echo.New()

	var err error

	dsn := "root:kahraman1@tcp(127.0.0.1:3306)/first"

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database Connection Failed.")
	}
	log.Println("Database Connection Successfull.")

	

	db.AutoMigrate(&User{})

	e.Use(middleware.Static("Register/static-files")) 
	e.GET("/register", RegisterPage)
	e.POST("/register", RegisterUser)

	e.Start(":8080")
}

func RegisterPage(c echo.Context) error {
	return c.File("static-files/index.html")
}

func RegisterUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	var existingUser User
	result := db.Where("Ad = ?", u.Ad).First(&existingUser)
	if result.Error == nil {
		return c.String(http.StatusBadRequest, "Kullanıcı adı zaten kullanılıyor!")
	}

	result = db.Create(u)
	if result.Error == nil {
		return c.String(http.StatusOK, "Kayıt başarılı!")
	}

	return c.String(http.StatusInternalServerError, "Kayıt işlemi sırasında bir hata oluştu.")
}
