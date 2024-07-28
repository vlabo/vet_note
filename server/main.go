package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"vet_note/db"
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

func main() {
	err := db.InitializeDB("test.db")
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to initialize db")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.GET("/v1/patient-list", getPatientList)
	e.GET("/v1/patient/:patientId", getPatient)
	e.POST("/v1/patient", updatePatient)
	e.DELETE("/v1/patient/:patientId", deletePatient)
	e.GET("/v1/procedure/:procedureId", getProcedure)
	e.POST("/v1/procedure/:patientId", updateProcedure)
	e.DELETE("/v1/procedure/:procedureId", deleteProcedure)
	e.Logger.Fatal(e.Start(":8080"))
}
