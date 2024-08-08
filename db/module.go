package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

type Patient struct {
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

func PatientFromJson(data string) (patient Patient, err error) {
	err = json.Unmarshal([]byte(data), &patient)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal patient data: %s", err)
	}
	return
}

type PatientMap map[string]interface{}

func PatientMapFromJson(data string) (PatientMap, error) {
	var patient PatientMap
	err := json.Unmarshal([]byte(data), &patient)
	if err != nil {
		return patient, fmt.Errorf("failed to unmarshal current patient data: %s", err)
	}
	return patient, nil
}

func (pm PatientMap) Update(newValues map[string]interface{}) error {
	for key, value := range newValues {
		if _, exists := pm[key]; exists {
			pm[key] = value
		} else {
			return fmt.Errorf("invalid field name %s", key)
		}
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
	patient, err := PatientMapFromJson(currentData)
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
	var currentPatient map[string]interface{}
	err = json.Unmarshal([]byte(currentData), &currentPatient)
	if err != nil {
		return fmt.Errorf("failed to unmarshal current patient data: %s", err)
	}

	// Find the procedure to update
	procedures, ok := currentPatient["Procedures"].([]interface{})
	if !ok {
		return fmt.Errorf("procedures field is missing or not an array")
	}

	var procedureToUpdate map[string]interface{}
	for _, proc := range procedures {
		procMap, ok := proc.(map[string]interface{})
		if !ok {
			continue
		}
		if procMap["ID"] == procedureId {
			procedureToUpdate = procMap
			break
		}
	}

	if procedureToUpdate == nil {
		return fmt.Errorf("procedure with id %s not found", procedureId)
	}

	// Apply changes to the procedure
	for key, value := range changes {
		procedureToUpdate[key] = value
	}

	// Update the modified_at field to the current time
	now := time.Now().Format(time.RFC3339)
	currentPatient["LastModified"] = now

	// Marshal the updated patient data into JSON
	updatedData, err := json.Marshal(currentPatient)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient data: %s", err)
	}

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
	var currentPatient Patient
	err = json.Unmarshal([]byte(currentData), &currentPatient)
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

	// Update the modified_at field to the current time
	now := time.Now().Format(time.RFC3339)
	currentPatient.LastModified = now

	// Marshal the updated patient data into JSON
	updatedData, err := json.Marshal(currentPatient)
	if err != nil {
		return fmt.Errorf("failed to marshal updated patient data: %s", err)
	}

	// Update the patient record in the database
	updateQuery := `UPDATE patients SET data = ?, modified_at = ? WHERE id = ?`
	_, err = db.Exec(updateQuery, string(updatedData), now, patientId)
	if err != nil {
		return fmt.Errorf("failed to update patient: %s", err)
	}

	return nil
}

func GetPatientTypes() (types string, err error) {
	// var settings Settings
	// result := connection.First(&settings)
	// if result.Error != nil {
	// 	err = fmt.Errorf("failed to get patient types: %s", result.Error)
	// }
	// types = settings.PatientTypes
	return
}

func UpdatePatientTypes(types string) {
	// var settings Settings
	// connection.First(&settings)
	// settings.PatientTypes = types
	// connection.Save(&settings)
}

func GetProcedureTypes() (types string, err error) {
	// var settings Settings
	// result := connection.First(&settings)
	// if result.Error != nil {
	// 	err = fmt.Errorf("failed to get procedure types: %s", result.Error)
	// }
	// types = settings.ProcedureTypes
	return
}

func UpdateProcedureTypes(types string) {
	// var settings Settings
	// connection.First(&settings)
	// settings.ProcedureTypes = types
	// connection.Save(&settings)
}
