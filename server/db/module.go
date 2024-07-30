package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var connection *gorm.DB = nil

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
	Procedures   []Procedure `gorm:"foreignKey:PatientID"`
}

type Settings struct {
	gorm.Model
	PatientTypes   string
	ProcedureTypes string
}

func InitializeDB(path string) error {
	var err error
	connection, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %s", err)
	}

	// Migrate the schema
	err = connection.AutoMigrate(&Patient{}, &Procedure{}, &Settings{})
	if err != nil {
		return fmt.Errorf("failed to migrate db: %s", err)
	}

	return nil
}

func GetPatient(id string) (patient Patient, err error) {
	result := connection.Preload("Procedures").First(&patient, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = fmt.Errorf("patient with Id %s not found", id)
		} else {
			err = fmt.Errorf("Database error")
		}
	}
	return
}

func GetProcedure(id string) (procedure Procedure, err error) {
	result := connection.First(&procedure, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			err = fmt.Errorf("procedure with Id %s not found", id)
		} else {
			err = fmt.Errorf("Database error")
		}
	}
	return
}

func GetPatientList() []Patient {
	var patients []Patient
	connection.Select("Id", "Type", "Name", "ChipId", "Owner", "OwnerPhone").Find(&patients)
	return patients
}

func UpdatePatient(patient *Patient) (err error) {
	if err = connection.Save(patient).Error; err != nil {
		err = fmt.Errorf("failed to update patient: %s", err)
	}
	return
}

func DeletePatient(patientId string) (err error) {
	if err = connection.Delete(&Patient{}, patientId).Error; err != nil {
		err = fmt.Errorf("failed to delete patient: %s", err)
	}
	return
}

func UpdateProcedure(procedure *Procedure) (err error) {
	if err = connection.Save(procedure).Error; err != nil {
		err = fmt.Errorf("failed to update procedure: %s", err)
	}
	return
}

func DeleteProcedure(procedureId string) (err error) {
	if err = connection.Delete(&Procedure{}, procedureId).Error; err != nil {
		err = fmt.Errorf("failed to delete procedure: %s", err)
	}
	return
}

func GetPatientTypes() (types string, err error) {
	var settings Settings
	result := connection.First(&settings)
	if result.Error != nil {
		err = fmt.Errorf("failed to get patient types: %s", result.Error)
	}
	types = settings.PatientTypes
	return
}

func UpdatePatientTypes(types string) {
	var settings Settings
	connection.First(&settings)
	settings.PatientTypes = types
	connection.Save(&settings)
}

func GetProcedureTypes() (types string, err error) {
	var settings Settings
	result := connection.First(&settings)
	if result.Error != nil {
		err = fmt.Errorf("failed to get procedure types: %s", result.Error)
	}
	types = settings.ProcedureTypes
	return
}

func UpdateProcedureTypes(types string) {
	var settings Settings
	connection.First(&settings)
	settings.ProcedureTypes = types
	connection.Save(&settings)
}
