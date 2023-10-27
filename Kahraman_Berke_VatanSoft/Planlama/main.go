package main

import (
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Plan struct {
	ID          uint

    Title       string
    Description string
    Date        time.Time `gorm:"type:datetime"`
    Status      string
}

var db *gorm.DB

func main() {
    e := echo.New()

    var err error

    dsn := "root:kahraman1@tcp(127.0.0.1:3306)/first"

    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Database'e Bağlanılamadı.")
    }
    log.Println("Database'e Bağlanıldı.")
    db.AutoMigrate(&Plan{})

    tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Student Daily Plans</title>
    </head>
    <body>
        <h1>Öğrenci Planları</h1>

        <form method="POST" action="/plan">
            <label for="Title">Title:</label>
            <input type="text" name="Title" required><br>

            <label for="Description">Description:</label>
            <input type="text" name="Description" required><br>

            <label for="Date">Date:</label>
            <input type="datetime-local" name="Date" required><br>

            <button type="submit">Add Plan</button>
        </form>

        <h2>Planlar:</h2>
        <ul>
            {{range .Plans}}
                <li>
				{{.Title}}, {{.Description}}, {{.Date.Format "2006-01-02T15:04"}}, {{.Status}}
				<a href="/edit/{{.ID}}">Düzenle</a>
				<a href="/delete/{{.ID}}">Sil</a>
                </li>
            {{end}}
        </ul>
    </body>
    </html>
    `

    t := &Template{
        templates: template.Must(template.New("index").Parse(tmpl)),
    }
    e.Renderer = t

    e.GET("/", PlanPage)
    e.GET("/plans", PlanPage)
    e.POST("/plan", PlanUser)
    e.GET("/edit/:id", EditPlan)
    e.POST("/update/:id", UpdatePlan)
    e.GET("/delete/:id", DeletePlan)
    e.Start(":8080")
}

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func PlanPage(c echo.Context) error {
    var plans []Plan
    db.Find(&plans)
    return c.Render(http.StatusOK, "index", map[string]interface{}{
        "Plans": plans,
    })
}

func PlanUser(c echo.Context) error {
    title := c.FormValue("Title")
    description := c.FormValue("Description")
    dateStr := c.FormValue("Date")

    date, err := parseDatetime(dateStr)
    if err != nil {
        return c.String(http.StatusBadRequest, "Geçersiz tarih formatı")
    }

    newPlan := Plan{
        Title:       title,
        Description: description,
        Date:        date,
        Status:      "Yapılacak",
    }

    db.Create(&newPlan)

    return c.Redirect(http.StatusSeeOther, "/plans")
}

func EditPlan(c echo.Context) error {
    id := c.Param("id")
    var plan Plan
    db.First(&plan, id)
    return c.Render(http.StatusOK, "edit", plan)
}

func UpdatePlan(c echo.Context) error {
    id := c.Param("id")
    var plan Plan
    db.First(&plan, id)

    title := c.FormValue("Title")
    description := c.FormValue("Description")
    dateStr := c.FormValue("Date")

    date, err := parseDatetime(dateStr)
    if err != nil {
        return c.String(http.StatusBadRequest, "Geçersiz tarih formatı")
    }

    plan.Title = title
    plan.Description = description
    plan.Date = date
    db.Save(&plan)

    return c.Redirect(http.StatusSeeOther, "/plans")
}

func DeletePlan(c echo.Context) error {
    id := c.Param("id")
    var plan Plan
    db.First(&plan, id)
    db.Delete(&plan)
    return c.Redirect(http.StatusSeeOther, "/plans")
}

func parseDatetime(datetimeStr string) (time.Time, error) {
    layout := "2006-01-02T15:04"
    t, err := time.Parse(layout, datetimeStr)
    if err != nil {
        return time.Time{}, err
    }
    return t, nil
}
