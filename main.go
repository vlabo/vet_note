package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"vet_note/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getPatientList(c echo.Context) error {
	patients := db.GetPatientList()

	view := make([]ViewListPatient, 0)
	for _, p := range patients {
		view = append(view, patientToListPatient(p))
	}

	return c.JSON(http.StatusCreated, view)
}

func getPatient(c echo.Context) error {
	// Get the patientId from the URL parameter
	id := c.Param("patientId")

	patient, err := db.GetPatient(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FmtError(err.Error()))
	}

	// Return the patient data as JSON
	return c.JSON(http.StatusOK, patientToViewPatient(patient))
}

func updatePatient(c echo.Context) error {
	var viewPatient ViewPatient

	// Bind the request body to the patient struct
	if err := c.Bind(&viewPatient); err != nil {
		return c.JSON(http.StatusBadRequest, FmtError("Invalid request body: %s", err))
	}
	patient := viewPatient.asPatient()
	err := db.UpdatePatient(&patient)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FmtError(err.Error()))
	}

	return c.JSON(http.StatusOK, patientToViewPatient(patient))
}

func deletePatient(c echo.Context) error {
	id := c.Param("patientId")
	err := db.DeletePatient(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FmtError("db error: %s", err))
	}
	return c.NoContent(http.StatusOK)
}

func getProcedure(c echo.Context) error {
	id := c.Param("procedureId")
	procedure, err := db.GetProcedure(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FmtError("Database error: %s", err))
	}

	return c.JSON(http.StatusOK, procedureToViewProcedure(procedure))
}

func updateProcedure(c echo.Context) error {
	id := c.Param("patientId")

	var viewProcedure ViewProcedure

	// Bind the request body to the procedure struct
	if err := c.Bind(&viewProcedure); err != nil {
		return c.JSON(http.StatusBadRequest, FmtError("Invalid request body: %s", err))
	}

	procedure := viewProcedure.asProcedure()
	patientId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FmtError("Invalid patient ID %d", err))
	}
	procedure.PatientID = uint(patientId)
	err = db.UpdateProcedure(&procedure)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FmtError("db error: %s", err))
	}

	return c.JSON(http.StatusOK, procedureToViewProcedure(procedure))
}

func deleteProcedure(c echo.Context) error {
	id := c.Param("procedureId")
	err := db.DeleteProcedure(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FmtError("db error: %s", err))
	}
	return c.NoContent(http.StatusOK)
}

func getPatientTypes(c echo.Context) error {
	types, err := db.GetPatientTypes()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read settings: %s", err))
	}
	return c.String(http.StatusOK, types)
}

func updatePatientTypes(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to update settings: %s", err))
	}
	db.UpdatePatientTypes(string(body))
	return c.NoContent(http.StatusOK)
}

func getProcedureTypes(c echo.Context) error {
	types, err := db.GetProcedureTypes()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read settings: %s", err))
	}
	return c.String(http.StatusOK, types)
}

func updateProcedureTypes(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to update settings: %s", err))
	}
	db.UpdateProcedureTypes(string(body))
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

//go:embed web/www/*
var embeddedFiles embed.FS

var (
	port         = flag.Int("port", 8000, "the port to listen on")
	databaseFile = flag.String("db", "", "path to the database")
	useCors      = flag.Bool("cors", false, "enable cors")
)

func main() {
	flag.Parse()
	if *databaseFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	err := db.InitializeDB(*databaseFile)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to initialize db")
	}
	username := os.Getenv("AUTH_USERNAME")
	password := os.Getenv("AUTH_PASSWORD")

	e := echo.New()

	e.Use(middleware.Logger())
	if *useCors {
		e.Use(middleware.CORS())
	}

	apiV1 := e.Group("/v1")
	apiV1.Use(basicAuthMiddleware(username, password))
	apiV1.GET("/authenticate", authenticate)
	apiV1.GET("/patient-list", getPatientList)
	apiV1.GET("/patient/:patientId", getPatient)
	apiV1.POST("/patient", updatePatient)
	apiV1.DELETE("/patient/:patientId", deletePatient)
	apiV1.GET("/procedure/:procedureId", getProcedure)
	apiV1.POST("/procedure/:patientId", updateProcedure)
	apiV1.DELETE("/procedure/:procedureId", deleteProcedure)
	apiV1.GET("/patient-types", getPatientTypes)
	apiV1.POST("/patient-types", updatePatientTypes)
	apiV1.GET("/procedure-types", getProcedureTypes)
	apiV1.POST("/procedure-types", updateProcedureTypes)

	subFS, err := fs.Sub(embeddedFiles, "web/www")
	if err != nil {
		fmt.Printf("%s", err)
		panic("failed to initialize www")
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(subFS))))

	e.Logger.Fatal(e.Start(fmt.Sprintf("127.0.0.1:%d", *port)))
}
