package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"strconv"

	"vet_note/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getPatientList(c echo.Context) error {
	patients, err := db.GetPatientList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, patients)
}

func getPatient(c echo.Context) error {
	// Get the patientId from the URL parameter
	id := c.Param("patientId")

	patient, err := db.GetPatient(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Return the patient data as JSON
	return c.JSON(http.StatusOK, patient)
}

func updatePatient(c echo.Context) error {
	var viewPatient db.ViewPatient

	// Bind the request body to the patient struct
	if err := c.Bind(&viewPatient); err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("Invalid request body: %s", err))
	}

	if _, ok := viewPatient.ID.Get(); ok {
		// Update entry
		err := db.UpdatePatient(viewPatient.AsSetter())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		// New entry
		id, err := db.CreatePatient(viewPatient.AsSetter())
		viewPatient.ID.Set(int32(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, viewPatient)

	}

	return c.NoContent(http.StatusOK)
}

func deletePatient(c echo.Context) error {
	id := c.Param("patientId")
	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("parse error: %s", err))
	}
	err = db.DeletePatient(int32(intID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("db error: %s", err))
	}
	return c.NoContent(http.StatusOK)
}

func getProcedure(c echo.Context) error {
	id := c.Param("procedureId")
	procedure, err := db.GetProcedure(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, db.FmtError("Database error: %s", err))
	}

	return c.JSON(http.StatusOK, procedure)
}

func updateProcedure(c echo.Context) error {
	var viewProcedure db.ViewProcedure

	// Bind the request body to the procedure struct
	if err := c.Bind(&viewProcedure); err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("Invalid request body: %s", err))
	}

	if _, ok := viewProcedure.ID.Get(); ok {
		// Update entry
		err := db.UpdateProcedure(viewProcedure.AsSetter())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		// New entry
		err := db.CreateProcedure(viewProcedure.AsSetter())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

func deleteProcedure(c echo.Context) error {
	id := c.Param("procedureId")

	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("parse error: %s", err))
	}
	err = db.DeleteProcedure(int32(intID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("db error: %s", err))
	}
	return c.NoContent(http.StatusOK)
}

func getPatientTypes(c echo.Context) error {
	types, err := db.GetPatientTypes()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read settings: %s", err))
	}

	return c.JSON(http.StatusOK, types)
}

func getProcedureTypes(c echo.Context) error {
	types, err := db.GetProcedureTypes()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read settings: %s", err))
	}

	return c.JSON(http.StatusOK, types)
}

func updateSetting(c echo.Context) error {
	viewSetting := db.ViewSetting{}
	if err := c.Bind(&viewSetting); err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("Invalid request body: %s", err))
	}
	if err := db.UpdateSetting(viewSetting); err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("Failed to update setting: %s", err))
	}
	return c.NoContent(http.StatusOK)
}

func updateSettings(c echo.Context) error {
	viewSettings := []db.ViewSetting{}
	if err := c.Bind(&viewSettings); err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("Invalid request body: %s", err))
	}
	for _, viewSetting := range viewSettings {
		if err := db.UpdateSetting(viewSetting); err != nil {
			return c.JSON(http.StatusBadRequest, db.FmtError("Failed to update setting: %s", err))
		}
	}
	return c.NoContent(http.StatusOK)
}

func deleteSetting(c echo.Context) error {
	id := c.Param("settingId")
	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("parse error: %s", err))
	}
	err = db.DeleteSetting(int32(intID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, db.FmtError("db error: %s", err))
	}
	return c.NoContent(http.StatusOK)
}

func authenticate(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func basicAuthMiddleware(username, password string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(providedUsername, providedPassword string, c echo.Context) (bool, error) {
		if providedUsername == username && providedPassword == password {
			return true, nil
		}
		return false, nil
	})
}

//go:embed ui/static/*
var embeddedFiles embed.FS

var (
	port         = flag.Int("port", 8000, "the port to listen on")
	databaseFile = flag.String("db", "", "path to the database")
	useCors      = flag.Bool("cors", false, "enable cors")
	dbLog        = flag.Bool("dbLog", false, "enable database logging")
)

func main() {
	flag.Parse()
	if *databaseFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	err := db.InitializeDB(*databaseFile, *dbLog)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to initialize db")
	}
	// username := os.Getenv("AUTH_USERNAME")
	// password := os.Getenv("AUTH_PASSWORD")

	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".mjs", "application/javascript")
	_ = mime.AddExtensionType(".cjs", "application/javascript")

	e := echo.New()
	// e.Use(middleware.Logger())
	if *useCors {
		e.Use(middleware.CORS())
	}

	apiV1 := e.Group("/v1")
	// FIXME: uncomment
	// apiV1.Use(basicAuthMiddleware(username, password))
	apiV1.GET("/authenticate", authenticate)
	apiV1.GET("/patient-list", getPatientList)
	apiV1.GET("/patient/:patientId", getPatient)
	apiV1.POST("/patient", updatePatient)
	apiV1.DELETE("/patient/:patientId", deletePatient)
	apiV1.GET("/procedure/:procedureId", getProcedure)
	apiV1.POST("/procedure", updateProcedure)
	apiV1.DELETE("/procedure/:procedureId", deleteProcedure)
	apiV1.GET("/patient-types", getPatientTypes)
	apiV1.GET("/procedure-types", getProcedureTypes)
	apiV1.POST("/setting", updateSetting)
	apiV1.POST("/settings", updateSettings)
	apiV1.DELETE("/setting/:settingId", deleteSetting)

	subFS, err := fs.Sub(embeddedFiles, "ui/static")
	if err != nil {
		fmt.Printf("%s", err)
		panic("failed to initialize www")
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(subFS))))
	// e.Static("/*", "ui/static")

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
