package main

import (
	"embed"
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

//go:embed web/www/*
var embeddedFiles embed.FS

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("supply database file argument: %s test.db", args[0])
		return
	}

	err := db.InitializeDB(args[1])
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to initialize db")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/v1/patient-list", getPatientList)
	e.GET("/v1/patient/:patientId", getPatient)
	e.POST("/v1/patient", updatePatient)
	e.DELETE("/v1/patient/:patientId", deletePatient)
	e.GET("/v1/procedure/:procedureId", getProcedure)
	e.POST("/v1/procedure/:patientId", updateProcedure)
	e.DELETE("/v1/procedure/:procedureId", deleteProcedure)
	e.GET("/v1/patient-types", getPatientTypes)
	e.POST("/v1/patient-types", updatePatientTypes)
	e.GET("/v1/procedure-types", getProcedureTypes)
	e.POST("/v1/procedure-types", updateProcedureTypes)

	subFS, err := fs.Sub(embeddedFiles, "web/www")
	if err != nil {
		fmt.Printf("%s", err)
		panic("failed to initialize www")
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(subFS))))

	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
