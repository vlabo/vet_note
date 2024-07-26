package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

type Patient struct {
	Id           uuid.UUID `gorm:"type:uuid;primary_key;"`
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
}

func (p *Patient) asViewListPatient() ViewListPatient {
	return ViewListPatient{
		Id:     p.Id,
		Type:   p.Type,
		Name:   p.Name,
		ChipId: p.ChipId,
		Owner:  p.Owner,
		Phone:  p.OwnerPhone,
	}
}

func (p *Patient) asViewPatient() ViewPatient {
	return ViewPatient{
		Id:           p.Id,
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
	patientId := c.Param("patientId")

	// Parse the patientId to a UUID
	id, err := uuid.Parse(patientId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid patient ID"})
	}

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

	// Check if the patient ID is provided
	if patient.Id == uuid.Nil {
		// If no ID, generate a new one for a new patient
		patient.Id = uuid.New()
	}

	// Set the LastModified field to the current time
	patient.LastModified = time.Now().Format("2006-01-02 15:04:05")

	// Check if the patient exists in the database
	var existingPatient Patient
	result := db.First(&existingPatient, "id = ?", patient.Id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// If the patient does not exist, create a new record
			if err := db.Create(&patient).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create patient"})
			}
			return c.JSON(http.StatusCreated, patient)
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// If the patient exists, update the record
	if err := db.Save(&patient).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update patient"})
	}

	return c.JSON(http.StatusOK, patient)
}

// func generateMockData() {
// 	mockPatients := []Patient{
// 		{
// 			Id:           uuid.New(),
// 			Type:         "Dog",
// 			Name:         "Buddy",
// 			Gender:       "male",
// 			BirthDate:    "2015-06-01",
// 			ChipId:       "123456789",
// 			Weight:       20.5,
// 			Castrated:    true,
// 			LastModified: "2023-10-01",
// 			Note:         "Healthy",
// 			Owner:        "John Doe",
// 			OwnerPhone:   "123-456-7890",
// 		},
// 		{
// 			Id:           uuid.New(),
// 			Type:         "Cat",
// 			Name:         "Whiskers",
// 			Gender:       "female",
// 			BirthDate:    "2018-09-15",
// 			ChipId:       "987654321",
// 			Weight:       10.2,
// 			Castrated:    false,
// 			LastModified: "2023-10-01",
// 			Note:         "Needs vaccination",
// 			Owner:        "Jane Smith",
// 			OwnerPhone:   "098-765-4321",
// 		},
// 		// Add more mock patients as needed
// 	}

// 	for _, patient := range mockPatients {
// 		db.Create(&patient)
// 	}
// }

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Patient{})
	if err != nil {
		fmt.Printf("Failed to migrate db: %s", err)
	}
	// generateMockData()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.GET("/v1/patient-list", get)
	e.GET("/v1/patient/:patientId", get_patient)
	e.POST("/v1/patient", update_patient)
	e.Logger.Fatal(e.Start(":8080"))
}
