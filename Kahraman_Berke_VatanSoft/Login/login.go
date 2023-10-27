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

func main(){

	e := echo.New()

	var err error

	dsn := "root:kahraman1@tcp(127.0.0.1:3306)/first"

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database Connection Failed.")
	}
	log.Println("Database Connection Successfull.")

	e.Use(middleware.Static("static-files"))


	e.GET("/login", LoginPage)
	e.POST("/login", LoginUser)
	e.GET("/Anasayfa", Anasayfa)

	e.Start(":8080")

}



func LoginPage(c echo.Context) error {
	return c.File("static-files/login.html")
}

func LoginUser(c echo.Context) error {
    email := c.FormValue("Email")
    sifre := c.FormValue("Sifre")

    var user User

 
    result := db.Table("users").Where("Email = ?", email).First(&user)
    if result.Error != nil {
        return c.String(http.StatusUnauthorized, "Giriş başarısız. Lütfen geçerli email ve şifre girin.")
    }


    if user.Sifre != sifre {
        return c.String(http.StatusUnauthorized, "Giriş başarısız. Lütfen geçerli email ve şifre girin.")
    }

    

	return c.String(http.StatusOK, "Giriş başarılı")

	
}




func Anasayfa(c echo.Context) error {
    data := map[string]interface{}{
        "Username": "Kullanıcı Adı", 
    }
    
    return c.Render(http.StatusOK, "anasayfa", data)
}
