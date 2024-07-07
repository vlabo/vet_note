package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/sahilm/fuzzy"
)

var db *badger.DB = nil

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func patient(c echo.Context) error {
	l := c.Logger()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		l.Errorf("invalid id: %s", err)
		return c.NoContent(http.StatusOK)
	}

	patient, err := getPatientById(db, id)
	if err != nil {
		l.Errorf("error getting patient: %v", err)
		return c.NoContent(http.StatusOK)
	}

	if htmx.IsHTMX(c.Request()) {
		return Render(c, http.StatusOK, renderPatient(patient))
	} else {
		return Render(c, http.StatusOK, mainWrapper(renderPatient(patient)))
	}
}

func patientEdit(c echo.Context) error {
	l := c.Logger()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		id = uuid.New()
	}
	patient, err := getPatientById(db, id)
	if err != nil {
		l.Errorf("error getting patient: %v", err)
		patient = &Patient{Id: id, Name: "", Owner: ""}
	}
	if htmx.IsHTMX(c.Request()) {
		return Render(c, http.StatusOK, renderPatientEdit(patient))
	} else {
		return Render(c, http.StatusOK, mainWrapper(renderPatientEdit(patient)))
	}
}

func patientUpdate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		id = uuid.New()
	}

	name := c.FormValue("name")
	owner := c.FormValue("owner")

	patient := &Patient{Id: id, Name: name, Owner: owner}
	if err := updatePatent(patient); err != nil {
		return err
	}

	return Render(c, http.StatusOK, renderPatient(patient))
}

func patientList(c echo.Context) error {
	return Render(c, http.StatusOK, renderPatientList(getPatientList(db)))
}

func newExamination(c echo.Context) error {
	l := c.Logger()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		id = uuid.New()
	}
	patient, err := getPatientById(db, id)
	if err != nil {
		l.Errorf("error getting patient: %v", err)
		return c.NoContent(http.StatusOK)
	}
	if htmx.IsHTMX(c.Request()) {
		return Render(c, http.StatusOK, renderNewExamination(patient))
	} else {
		return Render(c, http.StatusOK, mainWrapper(renderNewExamination(patient)))
	}
}

func saveExamination(c echo.Context) error {
	l := c.Logger()
	id, err := uuid.Parse(c.FormValue("Id"))
	if err != nil {
		l.Errorf("invalid patient id: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", err))
	}
	l.Debugf("UUID: %s", id)

	patient, err := getPatientById(db, id)
	if err != nil {
		l.Errorf("error getting patient: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("patient not found: %s", err))
	}
	// Get other form data
	dateStr := c.FormValue("Date")
	examType := c.FormValue("Type")
	description := c.FormValue("Description")

	// Parse the date
	date, err := time.Parse("2006-01-02T15:04:05.999Z", dateStr)
	if err != nil {
		l.Errorf("failed to parse data %s: %v", dateStr, err)
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	// Create a new examination
	examination := Examination{
		Date:        date,
		Type:        examType,
		Description: description,
	}
	patient.Examinations = append(patient.Examinations, examination)
	if err := updatePatent(patient); err != nil {
		return err
	}

	return Render(c, http.StatusOK, renderPatient(patient))
}

func patientSearch(c echo.Context) error {
	searchTerm := ""
	v, err := c.FormParams()
	if err == nil && v["search"] != nil && len(v["search"]) > 0 {
		searchTerm = v["search"][0]
	}

	if searchTerm == "" {
		return Render(c, http.StatusOK, renderPatientList(getPatientList(db)))
	}

	patients := Patients(getPatientList(db))
	matches := fuzzy.FindFrom(v["search"][0], patients)

	filtered := []Patient{}

	for _, r := range matches {
		filtered = append(filtered, patients[r.Index])
	}
	return Render(c, http.StatusOK, renderPatientList(filtered))
}

func index(c echo.Context) error {
	return Render(c, http.StatusOK, mainWrapper(renderMainView()))
}

func patients(c echo.Context) error {
	list := getPatientList(db)
	log.Printf("list: %+v", list)
	return c.JSON(http.StatusOK, list)
}

//go:embed static/index.html
var content embed.FS

func main() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	engine := echo.New()

	engine.Logger.SetLevel(log.DEBUG)
	engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	engine.GET("/", func(c echo.Context) error {
		data, err := content.ReadFile("static/index.html")
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "text/html", data)
	})
	engine.GET("/test", patients)
	engine.GET("/patient/new", patientEdit)
	engine.GET("/patient/search", patientSearch)
	engine.POST("/patient/search", patientSearch)
	engine.GET("/patient/:id", patient)
	engine.GET("/patient/:id/edit", patientEdit)
	engine.PUT("/patient/:id", patientUpdate)
	engine.GET("/examination/new/:id", newExamination)
	engine.POST("/examination/new", saveExamination)
	engine.Static("/assets", "assets")

	fmt.Println("Listening on :8080")
	engine.Logger.Fatal(engine.Start(":8080"))
}
