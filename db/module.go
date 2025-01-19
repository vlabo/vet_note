package db

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"vet_note/db/models"

	"github.com/aarondl/opt/omitnull"
	_ "modernc.org/sqlite"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
)

type SettingType string

const (
	PatientType   SettingType = "PatientType"
	ProcedureType SettingType = "ProcedureType"
)

var db bob.DB

//go:embed scheme.sql
var sqlScheme string

func InitializeDB(path string, _ bool) error {
	var err error
	db, err = bob.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("failed to connect database: %s", err)
	}

	// Create tables
	_, err = db.ExecContext(context.Background(), sqlScheme)
	if err != nil {
		return fmt.Errorf("failed to create tables: %s", err)
	}

	return nil
}

func GetPatient(id string) (patient ViewPatient, err error) {
	slog.Info("GetPatient", "id", id)
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}
	dbPatient, err := models.FindPatient(context.Background(), db, int32(idInt))
	if err != nil {
		return
	}

	patient = ViewPatient{
		ID:         dbPatient.ID,
		Type:       dbPatient.Type.GetOrZero(),
		Name:       dbPatient.Name.GetOrZero(),
		Gender:     dbPatient.Gender.GetOrZero(),
		Age:        dbPatient.Age.GetOrZero().String(),
		ChipID:     dbPatient.ChipID.GetOrZero(),
		Weight:     dbPatient.Weight.GetOrZero(),
		Castrated:  !dbPatient.Castrated.GetOrZero().IsZero(),
		Note:       dbPatient.Note.GetOrZero(),
		Owner:      dbPatient.Owner.GetOrZero(),
		OwnerPhone: dbPatient.OwnerPhone.GetOrZero(),
		Procedures: make([]ViewProcedure, 0, 2),
	}

	procedures, err := dbPatient.Procedures().All(context.Background(), db)
	if err != nil {
		return
	}

	for _, proc := range procedures {
		procedure := ViewProcedure{
			ID:        proc.ID,
			Type:      proc.Type.GetOrZero(),
			Date:      proc.Date.GetOrZero(),
			Details:   proc.Details.GetOrZero(),
			PatientID: idInt,
		}
		patient.Procedures = append(patient.Procedures, procedure)
	}

	return
}

func GetProcedure(id string) (procedure ViewProcedure, err error) {
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}
	dbProcedure, err := models.FindProcedure(context.Background(), db, int32(idInt))
	procedure = ViewProcedure{
		ID:        dbProcedure.ID,
		Type:      dbProcedure.Type.GetOrZero(),
		Date:      dbProcedure.Date.GetOrZero(),
		Details:   dbProcedure.Details.GetOrZero(),
		PatientID: int64(dbProcedure.PatientID.GetOrZero()),
	}

	return
}

func GetPatientList() ([]ViewListPatient, error) {
	dbPatients, err := models.Patients.Query().All(context.Background(), db)
	if err != nil {
		return nil, err
	}

	var patients []ViewListPatient
	for _, p := range dbPatients {
		patient := ViewListPatient{
			ID:     p.ID,
			Type:   p.Type.GetOrZero(),
			Name:   p.Name.GetOrZero(),
			ChipId: p.ChipID.GetOrZero(),
			Owner:  p.Owner.GetOrZero(),
			Phone:  p.OwnerPhone.GetOrZero(),
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

func setIfPresent(target func(string), patient map[string]string, key string) {
	if value, ok := patient[key]; ok {
		target(value)
	}
}

// Helper function for decimal values
func setDecimalIfPresent(target func(decimal.Decimal), patient map[string]string, key string) {
	if value, ok := patient[key]; ok {
		decimalValue, err := decimal.NewFromString(value) // Convert string to decimal
		if err == nil {                                   // Only set if the conversion is successful
			target(decimalValue)
		}
	}
}

// Helper function for decimal values
func setFloatIfPresent(target func(float32), patient map[string]string, key string) {
	if value, ok := patient[key]; ok {
		decimalValue, err := strconv.ParseFloat(value, 32) // Convert string to decimal
		if err == nil {                                    // Only set if the conversion is successful
			target(float32(decimalValue))
		}
	}
}

func CreatePatient(patient map[string]string) error {
	slog.Info("CreatePatient", "patient", patient)

	setter := models.PatientSetter{
		CreatedAt: omitnull.From(time.Now()),
		UpdatedAt: omitnull.From(time.Now()),
	}
	setIfPresent(setter.Type.Set, patient, "type")
	setIfPresent(setter.Name.Set, patient, "name")
	setIfPresent(setter.Gender.Set, patient, "gender")
	setDecimalIfPresent(setter.Age.Set, patient, "name")
	setIfPresent(setter.ChipID.Set, patient, "chip_id")
	setFloatIfPresent(setter.Weight.Set, patient, "weight")
	setDecimalIfPresent(setter.Castrated.Set, patient, "castrated")
	setIfPresent(setter.Note.Set, patient, "note")
	setIfPresent(setter.Owner.Set, patient, "owner")
	setIfPresent(setter.OwnerPhone.Set, patient, "owner_phone")
	_, err := models.Patients.Insert(&setter).Exec(context.Background(), db)
	return err
}

func UpdatePatient(id int64, changes *ViewPatient) error {
	slog.Info("UpdatePatient", "id", id, "changes", changes)

	// // Initialize the base query and parameters slice
	// query := "UPDATE patients SET "
	// params := []interface{}{}

	// // Dynamically build the query based on non-null fields
	// if changes.Type != nil {
	// 	query += "type = ?, "
	// 	params = append(params, *changes.Type)
	// }
	// if changes.Name != nil {
	// 	query += "name = ?, "
	// 	params = append(params, *changes.Name)
	// }
	// if changes.Gender != nil {
	// 	query += "gender = ?, "
	// 	params = append(params, *changes.Gender)
	// }
	// if changes.Age != nil {
	// 	query += "age = ?, "
	// 	params = append(params, *changes.Age)
	// }
	// if changes.ChipID != nil {
	// 	query += "chip_id = ?, "
	// 	params = append(params, *changes.ChipID)
	// }
	// if changes.Weight != nil {
	// 	query += "weight = ?, "
	// 	params = append(params, *changes.Weight)
	// }
	// if changes.Castrated != nil {
	// 	query += "castrated = ?, "
	// 	params = append(params, *changes.Castrated)
	// }
	// if changes.Note != nil {
	// 	query += "note = ?, "
	// 	params = append(params, *changes.Note)
	// }
	// if changes.Owner != nil {
	// 	query += "owner = ?, "
	// 	params = append(params, *changes.Owner)
	// }
	// if changes.OwnerPhone != nil {
	// 	query += "owner_phone = ?, "
	// 	params = append(params, *changes.OwnerPhone)
	// }

	// // Add the updated_at field and the WHERE clause
	// query += "updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	// params = append(params, id)

	// // Execute the query
	// _, err := db.Exec(query, params...)
	// if err != nil {
	// 	return fmt.Errorf("failed to update patient: %s", err)
	// }
	return nil
}

func DeletePatient(patientId string) error {
	// query := `DELETE FROM patients WHERE id = ?`
	// _, err := db.Exec(query, patientId)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete patient: %s", err)
	// }
	return nil
}

func CreateProcedure(procedure *ViewProcedure) error {
	// slog.Info("CreateProcedure", "procedure", procedure)
	// query := `INSERT INTO procedures (type, date, details, patient_id) VALUES (?, ?, ?, ?)`
	// result, err := db.Exec(query, procedure.Type, procedure.Date, procedure.Details, procedure.PatientID)
	// if err != nil {
	// 	return fmt.Errorf("failed to create procedure: %s", err)
	// }

	// // Get the ID of the newly inserted record
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return fmt.Errorf("failed to retrieve last insert ID: %s", err)
	// }
	// procedure.ID = &id
	return nil
}

func UpdateProcedure(id int64, changes *ViewProcedure) error {
	// slog.Info("UpdateProcedure", "id", id, "changes", changes)

	// // Initialize the base query and parameters slice
	// query := "UPDATE procedures SET "
	// params := []interface{}{}

	// // Dynamically build the query based on non-null fields
	// if changes.Type != nil {
	// 	query += "type = ?, "
	// 	params = append(params, *changes.Type)
	// }
	// if changes.Date != nil {
	// 	query += "date = ?, "
	// 	params = append(params, *changes.Date)
	// }
	// if changes.Details != nil {
	// 	query += "details = ?, "
	// 	params = append(params, *changes.Details)
	// }

	// // Add the updated_at field and the WHERE clause
	// query += "updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	// params = append(params, id)

	// // Execute the query
	// _, err := db.Exec(query, params...)
	// if err != nil {
	// 	return fmt.Errorf("failed to update procedure: %s", err)
	// }
	return nil
}

func DeleteProcedure(procedureId string) error {
	// query := `DELETE FROM procedures WHERE id = ?`
	// _, err := db.Exec(query, procedureId)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete procedure: %s", err)
	// }
	return nil
}

func GetPatientTypes() ([]ViewSetting, error) {
	// query := `SELECT id, value, type, idx FROM settings WHERE type = ? ORDER BY idx`
	// rows, err := db.Query(query, PatientType)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var settings []ViewSetting
	// for rows.Next() {
	// 	var setting ViewSetting
	// 	err = rows.Scan(&setting.ID, &setting.Value, &setting.Type, &setting.Index)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	settings = append(settings, setting)
	// }
	return nil, nil
}

func GetProcedureTypes() ([]ViewSetting, error) {
	// query := `SELECT id, value, type, idx FROM settings WHERE type = ? ORDER BY idx`
	// rows, err := db.Query(query, ProcedureType)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var settings []ViewSetting
	// for rows.Next() {
	// 	var setting ViewSetting
	// 	err = rows.Scan(&setting.ID, &setting.Value, &setting.Type, &setting.Index)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	settings = append(settings, setting)
	// }
	return nil, nil
}

func CreateSetting(setting ViewSetting) error {
	// query := `INSERT INTO settings (value, type, idx) VALUES (?, ?, ?)`
	// _, err := db.Exec(query, setting.Value, setting.Type, setting.Index, setting.ID)
	// if err != nil {
	// 	return fmt.Errorf("failed to update setting: %s", err)
	// }
	return nil
}

func UpdateSetting(setting ViewSetting) error {
	// TODO: make fields optional.
	// query := `UPDATE settings SET value = ?, type = ?, idx = ?`
	// _, err := db.Exec(query, setting.Value, setting.Type, setting.Index, setting.ID)
	// if err != nil {
	// 	return fmt.Errorf("failed to update setting: %s", err)
	// }
	return nil
}

func DeleteSetting(settingId string) error {
	// query := `DELETE FROM settings WHERE id = ?`
	// _, err := db.Exec(query, settingId)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete setting: %s", err)
	// }
	return nil
}
