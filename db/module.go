package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

type SettingType string

const (
	PatientType   SettingType = "PatientType"
	ProcedureType SettingType = "ProcedureType"
)

type Setting struct {
	gorm.Model
	Value string
	Type  SettingType
	Idx   uint
}

func InitializeDB(path string, enableLogging bool) error {
	config := gorm.Config{}
	if enableLogging {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	}
	var err error
	connection, err = gorm.Open(sqlite.Open(path), &config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %s", err)
	}

	// Migrate the schema
	err = connection.AutoMigrate(&Patient{}, &Procedure{}, &Setting{})
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

func CreatePatient(patient *Patient) (err error) {
	if err = connection.Create(patient).Error; err != nil {
		err = fmt.Errorf("failed to update patient: %s", err)
	}
	return
}

func UpdatePatient(id uint, changes *Patient) (err error) {
	patient := Patient{Model: gorm.Model{ID: id}}
	if err = connection.Model(&patient).Updates(changes).Error; err != nil {
		err = fmt.Errorf("failed to update patient: %w", err)
	}
	return
}

func DeletePatient(patientId string) (err error) {
	if err = connection.Delete(&Patient{}, patientId).Error; err != nil {
		err = fmt.Errorf("failed to delete patient: %s", err)
	}
	return
}

func CreateProcedure(procedure *Procedure) (err error) {
	if err = connection.Create(procedure).Error; err != nil {
		err = fmt.Errorf("failed to update procedure: %s", err)
	}
	return
}

func UpdateProcedure(id uint, changes *Procedure) (err error) {
	procedure := Procedure{Model: gorm.Model{ID: id}}
	if err = connection.Model(&procedure).Updates(changes).Error; err != nil {
		err = fmt.Errorf("failed to update procedure: %w", err)
	}
	return
}

func DeleteProcedure(procedureId string) (err error) {
	if err = connection.Delete(&Procedure{}, procedureId).Error; err != nil {
		err = fmt.Errorf("failed to delete procedure: %s", err)
	}
	return
}

func GetPatientTypes() (settings []Setting, err error) {
	err = connection.Model(Setting{}).Where("type = ?", PatientType).Order("idx").Find(&settings).Error
	return
}

func GetProcedureTypes() (settings []Setting, err error) {
	err = connection.Model(Setting{}).Where("type = ?", ProcedureType).Order("idx").Find(&settings).Error
	return
}

func UpdateSetting(setting Setting) error {
	return connection.Save(&setting).Error
}

func DeleteSetting(settingId string) (err error) {
	if err = connection.Delete(&Setting{}, settingId).Error; err != nil {
		err = fmt.Errorf("failed to delete setting: %s", err)
	}
	return
}
