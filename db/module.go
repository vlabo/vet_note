package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	_ "modernc.org/sqlite"
)

type Base struct {
	ID         string
	CreatedAt  time.Time
	ModifiedAt time.Time
	Data       string
}

type Procedure struct {
	ID      string
	Type    string
	Date    string
	Details string
}

func (p *Procedure) Update(newValues map[string]interface{}) error {
	v := reflect.ValueOf(p).Elem()
	for key, value := range newValues {
		field := v.FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("invalid field name %s", key)
		}
		if !field.CanSet() {
			return fmt.Errorf("cannot set field %s", key)
		}
		val := reflect.ValueOf(value)
		if field.Type() != val.Type() {
			return fmt.Errorf("provided value type didn't match struct field type for field %s", key)
		}
		field.Set(val)
	}
	return nil
}

type Patient struct {
	Type       string
	Name       string
	Gender     string
	BirthDate  string
	ChipId     string
	Weight     float64
	Castrated  bool
	Note       string
	Owner      string
	OwnerPhone string
	Procedures []Procedure
}

func PatientFromJson(data string) (patient Patient, err error) {
	err = json.Unmarshal([]byte(data), &patient)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal patient data: %s", err)
	}
	return
}

func (p *Patient) Update(newValues map[string]interface{}) error {
	v := reflect.ValueOf(p).Elem()
	for key, value := range newValues {
		field := v.FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("invalid field name %s", key)
		}
		if !field.CanSet() {
			return fmt.Errorf("cannot set field %s", key)
		}
		val := reflect.ValueOf(value)
		if field.Type() != val.Type() {
			return fmt.Errorf("provided value type didn't match struct field type for field %s", key)
		}
		field.Set(val)
	}
	return nil
}

type SettingValue struct {
	Id    string
	Value string
}

type Settings struct {
	PatientTypes   []SettingValue
	ProcedureTypes []SettingValue
}

var db *sql.DB

func InitializeDB(path string) error {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("failed to connect database: %s", err)
	}

	_, err = db.Exec(getCreateTableQuery("patients"))
	if err != nil {
		return fmt.Errorf("failed to create table patients: %s", err)
	}

	_, err = db.Exec(getCreateTableQuery("settings"))
	if err != nil {
		return fmt.Errorf("failed to create table settings: %s", err)
	}
	return nil
}

func getCreateTableQuery(tableName string) string {
	return fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id TEXT PRIMARY KEY,
        deleted_at DATETIME,
        created_at DATETIME NOT NULL,
        modified_at DATETIME NOT NULL,
        data TEXT
    );`, tableName)
}

func GetPatient(id string) (patient Patient, err error) {
	// Query the database for the patient with the specified ID
	query := `SELECT data FROM patients WHERE id = ?`
	row := db.QueryRow(query, id)
	var data string
	err = row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return patient, fmt.Errorf("patient with id %s not found", id)
		}
		return patient, fmt.Errorf("failed to query patient: %s", err)
	}
	patient, err = PatientFromJson(data)
	if err != nil {
		return patient, fmt.Errorf("corrupted database: %s", err)
	}

	return
}

func GetPatientList() ([]Patient, error) {
	// Query the database for all patients where deleted_at is NULL
	query := `SELECT data FROM patients`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query patients: %s", err)
	}
	defer rows.Close()

	var patients []Patient

	// Iterate over the result set
	for rows.Next() {
		var data string
		err := rows.Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan patient data: %s", err)
		}

		patient, err := PatientFromJson(data)
		if err != nil {
			return nil, fmt.Errorf("corrupted database: %s", err)
		}

		patients = append(patients, patient)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %s", err)
	}
	return patients, nil
}

func UpdatePatient(id string, changes map[string]interface{}) (err error) {
	// Fetch the current state of the patient record
	var currentData string
	query := `SELECT data FROM patients WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&currentData)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("patient with id %s not found", id)
		}
		return fmt.Errorf("failed to query patient: %s", err)
	}

	// Unmarshal the current data into a map
	patient, err := PatientFromJson(currentData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal current patient data: %s", err)
	}
	_ = patient.Update(changes)

	// Marshal the updated patient data into JSON
	updatedData, err := json.Marshal(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient data: %s", err)
	}

	// Update the modified_at field to the current time
	now := time.Now().Format(time.RFC3339)

	// Update the patient record in the database
	updateQuery := `UPDATE patients SET data = ?, modified_at = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, string(updatedData), now, id)
	if err != nil {
		return fmt.Errorf("failed to update patient: %s", err)
	}

	return nil
}

func DeletePatient(patientId string) (err error) {
	// Delete the patient record from the database
	query := `DELETE FROM patients WHERE id = ?`
	_, err = db.Exec(query, patientId)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %s", err)
	}
	return nil
}

func UpdateProcedure(patientId string, procedureId string, changes map[string]interface{}) (err error) {
	// Fetch the current state of the patient record
	var currentData string
	query := `SELECT data FROM patients WHERE id = ?`
	err = db.QueryRow(query, patientId).Scan(&currentData)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("patient with id %s not found", patientId)
		}
		return fmt.Errorf("failed to query patient: %s", err)
	}

	// Unmarshal the current data into a Patient struct
	patient, err := PatientFromJson(currentData)
	// var currentPatient map[string]interface{}
	// err = json.Unmarshal([]byte(currentData), &currentPatient)
	if err != nil {
		return fmt.Errorf("failed to unmarshal current patient data: %s", err)
	}

	// Find the procedure to update
	for i := range patient.Procedures {
		if patient.Procedures[i].ID == procedureId {
			err := patient.Procedures[i].Update(changes)
			if err != nil {
				return fmt.Errorf("failed to update procedure: %s", err)
			}
			break
		}
	}
	// Marshal the updated patient data into JSON
	updatedData, err := json.Marshal(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient data: %s", err)
	}

	// Update the modified_at field to the current time
	now := time.Now().Format(time.RFC3339)
	// Update the patient record in the database
	updateQuery := `UPDATE patients SET data = ?, modified_at = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, string(updatedData), now, patientId)
	if err != nil {
		return fmt.Errorf("failed to update patient: %s", err)
	}

	return nil
}

func DeleteProcedure(patientId string, procedureId string) (err error) {
	// Fetch the current state of the patient record
	var currentData string
	query := `SELECT data FROM patients WHERE id = ?`
	err = db.QueryRow(query, patientId).Scan(&currentData)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("patient with id %s not found", patientId)
		}
		return fmt.Errorf("failed to query patient: %s", err)
	}

	// Unmarshal the current data into a Patient struct
	currentPatient, err := PatientFromJson(currentData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal current patient data: %s", err)
	}

	// Find and delete the procedure
	found := false
	for i, proc := range currentPatient.Procedures {
		if proc.ID == procedureId {
			currentPatient.Procedures = append(currentPatient.Procedures[:i], currentPatient.Procedures[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("procedure with id %s not found", procedureId)
	}

	// Marshal the updated patient data into JSON
	updatedData, err := json.Marshal(currentPatient)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient data: %s", err)
	}

	// Update the modified_at field to the current time
	now := time.Now().Format(time.RFC3339)

	// Update the patient record in the database
	updateQuery := `UPDATE patients SET data = ?, modified_at = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, string(updatedData), now, patientId)
	if err != nil {
		return fmt.Errorf("failed to update patient: %s", err)
	}

	return nil
}

func GetSettings() (Settings, error) {
	var settings Settings
	var data string

	query := `SELECT data FROM settings WHERE id = 'settings'`
	err := db.QueryRow(query).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return settings, fmt.Errorf("settings not found")
		}
		return settings, fmt.Errorf("failed to query settings: %s", err)
	}

	err = json.Unmarshal([]byte(data), &settings)
	if err != nil {
		return settings, fmt.Errorf("failed to unmarshal settings data: %s", err)
	}

	return settings, nil
}

func UpdatePatientSettings(newPatientTypes []SettingValue) error {
	settings, err := GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get current settings: %s", err)
	}

	// TODO: dont replace, update
	settings.PatientTypes = newPatientTypes

	return saveSettings(settings)
}

func DeletePatientSettings(idsToDelete []string) error {
	settings, err := GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get current settings: %s", err)
	}

	updatedPatientTypes := []SettingValue{}
	for _, setting := range settings.PatientTypes {
		if !contains(idsToDelete, setting.Id) {
			updatedPatientTypes = append(updatedPatientTypes, setting)
		}
	}
	settings.PatientTypes = updatedPatientTypes

	return saveSettings(settings)
}

func UpdateProcedureSettings(newProcedureTypes []SettingValue) error {
	settings, err := GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get current settings: %s", err)
	}
	// TODO: dont replace, update
	settings.ProcedureTypes = newProcedureTypes

	return saveSettings(settings)
}

func DeleteProcedureSettings(idsToDelete []string) error {
	settings, err := GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get current settings: %s", err)
	}

	updatedProcedureTypes := []SettingValue{}
	for _, setting := range settings.ProcedureTypes {
		if !contains(idsToDelete, setting.Id) {
			updatedProcedureTypes = append(updatedProcedureTypes, setting)
		}
	}
	settings.ProcedureTypes = updatedProcedureTypes

	return saveSettings(settings)
}

func saveSettings(settings Settings) error {
	updatedData, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal updated settings data: %s", err)
	}

	now := time.Now().Format(time.RFC3339)
	updateQuery := `UPDATE settings SET data = ?, modified_at = ? WHERE id = 'settings'`
	_, err = db.Exec(updateQuery, string(updatedData), now)
	if err != nil {
		return fmt.Errorf("failed to update settings: %s", err)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
