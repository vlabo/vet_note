package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

type Procedure struct {
	gorm.Model
	Type      string
	Date      string
	Details   string
	PatientID uint
}

type Patient struct {
	gorm.Model
	Type         string
	Name         string
	Gender       string
	BirthDate    string
	ChipId       string
	Weight       float64
	Castrated    bool
	LastModified string
	Note         string
	Owner        string
	OwnerPhone   string
	Procedures   []Procedure
}

func (p *Patient) asViewListPatient() ViewListPatient {
	return ViewListPatient{
		Id:     strconv.FormatUint(uint64(p.ID), 10),
		Type:   p.Type,
		Name:   p.Name,
		ChipId: p.ChipId,
		Owner:  p.Owner,
		Phone:  p.OwnerPhone,
	}
}

func (p *Patient) asViewPatient() ViewPatient {
	return ViewPatient{
		Id:           strconv.FormatUint(uint64(p.ID), 10),
		Type:         p.Type,
		Name:         p.Name,
		Gender:       p.Gender,
		BirthDate:    p.BirthDate,
		ChipId:       p.ChipId,
		Weight:       p.Weight,
		Castrated:    p.Castrated,
		LastModified: p.LastModified,
		Note:         p.Note,
		Owner:        p.Owner,
		OwnerPhone:   p.OwnerPhone,
		Procedures:   []ViewProcedure{},
	}
}

func get(c echo.Context) error {
	var patients []Patient
	db.Select("Id", "Type", "Name", "ChipId", "Owner", "OwnerPhone").Find(&patients)

	view := make([]ViewListPatient, 0)
	for _, p := range patients {
		view = append(view, p.asViewListPatient())
	}

	return c.JSON(http.StatusCreated, view)
}

func get_patient(c echo.Context) error {
	// Get the patientId from the URL parameter
	id := c.Param("patientId")

	// Parse the patientId to a UUID
	// id, err := uuid.Parse(patientId)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid patient ID"})
	// }

	// Find the patient by ID
	var patient Patient
	result := db.First(&patient, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Patient not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Return the patient data as JSON
	return c.JSON(http.StatusOK, patient.asViewPatient())
}

func update_patient(c echo.Context) error {
	var viewPatient ViewPatient

	// Bind the request body to the patient struct
	if err := c.Bind(&viewPatient); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	patient := viewPatient.asPatient()
	// Check if the patient exists in the database
	var existingPatient Patient
	result := db.First(&existingPatient, "id = ?", patient.ID)

	exists := true

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			exists = false
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
	}

	if exists {
		// If the patient exists, update the record
		if err := db.Save(&patient).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update patient"})
		}
	} else {
		if err := db.Create(&patient).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create patient"})
		}
	}

	return c.JSON(http.StatusOK, patient.asViewPatient())
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Procedure{})
	if err != nil {
		fmt.Printf("Failed to migrate db: %s", err)
	}
	err = db.AutoMigrate(&Patient{})
	if err != nil {
		fmt.Printf("Failed to migrate db: %s", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.GET("/v1/patient-list", get)
	e.GET("/v1/patient/:patientId", get_patient)
	e.POST("/v1/patient", update_patient)
	e.Logger.Fatal(e.Start(":8080"))
}
