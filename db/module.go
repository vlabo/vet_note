package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type SettingType string

const (
	PatientType   SettingType = "PatientType"
	ProcedureType SettingType = "ProcedureType"
)

var db *sql.DB

func InitializeDB(path string, _ bool) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to connect database: %s", err)
	}

	// Create tables
	createTablesSQL := `
    CREATE TABLE IF NOT EXISTS patients (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        type TEXT,
        name TEXT,
        gender TEXT,
        birth_date TEXT,
        age TEXT,
        chip_id TEXT,
        weight REAL,
        castrated NUMERIC DEFAULT 0,
        note TEXT,
        owner TEXT,
        owner_phone TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME
    );

    CREATE TABLE IF NOT EXISTS procedures (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        type TEXT,
        date TEXT,
        details TEXT,
        patient_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME,
        FOREIGN KEY(patient_id) REFERENCES patients(id)
    );

    CREATE TABLE IF NOT EXISTS settings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        value TEXT,
        type TEXT,
        idx INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME
    );

    `
	_, err = db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %s", err)
	}

	migration := `ALTER TABLE patients ADD COLUMN age TEXT;`

	_, err = db.Exec(migration)
	if err != nil {
		slog.Warn("failed to run migration", "err", err)
	}
	return nil
}

func GetPatient(id string) (patient ViewPatient, err error) {
	slog.Info("GetPatient", "id", id)
	query := `SELECT id, type, name, gender, age, chip_id, weight, castrated, note, owner, owner_phone FROM patients WHERE id = ?`
	row := db.QueryRow(query, id)
	err = row.Scan(&patient.ID, &patient.Type, &patient.Name, &patient.Gender, &patient.Age, &patient.ChipId, &patient.Weight, &patient.Castrated, &patient.Note, &patient.Owner, &patient.OwnerPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("patient with Id %s not found", id)
		} else {
			err = fmt.Errorf("database error: %s", err)
		}
		return
	}

	// Load procedures
	query = `SELECT id, type, date, details, patient_id FROM procedures WHERE patient_id = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var procedure ViewProcedure
		err = rows.Scan(&procedure.ID, &procedure.Type, &procedure.Date, &procedure.Details, &procedure.PatientID)
		if err != nil {
			return
		}
		if patient.Procedures == nil {
			procedures := make([]ViewProcedure, 0, 2)
			patient.Procedures = &procedures
		}
		*patient.Procedures = append(*patient.Procedures, procedure)
	}

	return
}

func GetProcedure(id string) (procedure ViewProcedure, err error) {
	slog.Info("GetProcedure", "id", id)
	query := `SELECT id, type, date, details, patient_id FROM procedures WHERE id = ?`
	row := db.QueryRow(query, id)
	err = row.Scan(&procedure.ID, &procedure.Type, &procedure.Date, &procedure.Details, &procedure.PatientID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("procedure with Id %s not found", id)
		} else {
			err = fmt.Errorf("database error: %s", err)
		}
	}
	return
}

func GetPatientList() ([]ViewPatient, error) {
	query := `SELECT id, type, name, chip_id, owner, owner_phone FROM patients ORDER BY updated_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []ViewPatient
	for rows.Next() {
		var patient ViewPatient
		err = rows.Scan(&patient.ID, &patient.Type, &patient.Name, &patient.ChipId, &patient.Owner, &patient.OwnerPhone)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func CreatePatient(patient *ViewPatient) error {
	slog.Info("CreatePatient", "patient", patient)
	query := `INSERT INTO patients (type, name, gender, age, chip_id, weight, castrated, note, owner, owner_phone) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, patient.Type, patient.Name, patient.Gender, patient.Age, patient.ChipId, patient.Weight, patient.Castrated, patient.Note, patient.Owner, patient.OwnerPhone)
	if err != nil {
		return fmt.Errorf("failed to create patient: %s", err)
	}
	// Get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %s", err)
	}
	patient.ID = &id
	return nil
}

func UpdatePatient(id int64, changes *ViewPatient) error {
	slog.Info("UpdatePatient", "id", id, "changes", changes)

	// Initialize the base query and parameters slice
	query := "UPDATE patients SET "
	params := []interface{}{}

	// Dynamically build the query based on non-null fields
	if changes.Type != nil {
		query += "type = ?, "
		params = append(params, *changes.Type)
	}
	if changes.Name != nil {
		query += "name = ?, "
		params = append(params, *changes.Name)
	}
	if changes.Gender != nil {
		query += "gender = ?, "
		params = append(params, *changes.Gender)
	}
	if changes.Age != nil {
		query += "age = ?, "
		params = append(params, *changes.Age)
	}
	if changes.ChipId != nil {
		query += "chip_id = ?, "
		params = append(params, *changes.ChipId)
	}
	if changes.Weight != nil {
		query += "weight = ?, "
		params = append(params, *changes.Weight)
	}
	if changes.Castrated != nil {
		query += "castrated = ?, "
		params = append(params, *changes.Castrated)
	}
	if changes.Note != nil {
		query += "note = ?, "
		params = append(params, *changes.Note)
	}
	if changes.Owner != nil {
		query += "owner = ?, "
		params = append(params, *changes.Owner)
	}
	if changes.OwnerPhone != nil {
		query += "owner_phone = ?, "
		params = append(params, *changes.OwnerPhone)
	}

	// Add the updated_at field and the WHERE clause
	query += "updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	params = append(params, id)

	// Execute the query
	_, err := db.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to update patient: %s", err)
	}
	return nil
}

func DeletePatient(patientId string) error {
	query := `DELETE FROM patients WHERE id = ?`
	_, err := db.Exec(query, patientId)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %s", err)
	}
	return nil
}

func CreateProcedure(procedure *ViewProcedure) error {
	slog.Info("CreateProcedure", "procedure", procedure)
	query := `INSERT INTO procedures (type, date, details, patient_id) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, procedure.Type, procedure.Date, procedure.Details, procedure.PatientID)
	if err != nil {
		return fmt.Errorf("failed to create procedure: %s", err)
	}

	// Get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %s", err)
	}
	procedure.ID = &id
	return nil
}

func UpdateProcedure(id int64, changes *ViewProcedure) error {
	slog.Info("UpdateProcedure", "id", id, "changes", changes)

	// Initialize the base query and parameters slice
	query := "UPDATE procedures SET "
	params := []interface{}{}

	// Dynamically build the query based on non-null fields
	if changes.Type != nil {
		query += "type = ?, "
		params = append(params, *changes.Type)
	}
	if changes.Date != nil {
		query += "date = ?, "
		params = append(params, *changes.Date)
	}
	if changes.Details != nil {
		query += "details = ?, "
		params = append(params, *changes.Details)
	}

	// Add the updated_at field and the WHERE clause
	query += "updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	params = append(params, id)

	// Execute the query
	_, err := db.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to update procedure: %s", err)
	}
	return nil
}

func DeleteProcedure(procedureId string) error {
	query := `DELETE FROM procedures WHERE id = ?`
	_, err := db.Exec(query, procedureId)
	if err != nil {
		return fmt.Errorf("failed to delete procedure: %s", err)
	}
	return nil
}

func GetPatientTypes() ([]ViewSetting, error) {
	query := `SELECT id, value, type, idx FROM settings WHERE type = ? ORDER BY idx`
	rows, err := db.Query(query, PatientType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []ViewSetting
	for rows.Next() {
		var setting ViewSetting
		err = rows.Scan(&setting.ID, &setting.Value, &setting.Type, &setting.Index)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func GetProcedureTypes() ([]ViewSetting, error) {
	query := `SELECT id, value, type, idx FROM settings WHERE type = ? ORDER BY idx`
	rows, err := db.Query(query, ProcedureType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []ViewSetting
	for rows.Next() {
		var setting ViewSetting
		err = rows.Scan(&setting.ID, &setting.Value, &setting.Type, &setting.Index)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func CreateSetting(setting ViewSetting) error {
	query := `INSERT INTO settings (value, type, idx) VALUES (?, ?, ?)`
	_, err := db.Exec(query, setting.Value, setting.Type, setting.Index, setting.ID)
	if err != nil {
		return fmt.Errorf("failed to update setting: %s", err)
	}
	return nil
}

func UpdateSetting(setting ViewSetting) error {
	// TODO: make fields optional.
	query := `UPDATE settings SET value = ?, type = ?, idx = ?`
	_, err := db.Exec(query, setting.Value, setting.Type, setting.Index, setting.ID)
	if err != nil {
		return fmt.Errorf("failed to update setting: %s", err)
	}
	return nil
}

func DeleteSetting(settingId string) error {
	query := `DELETE FROM settings WHERE id = ?`
	_, err := db.Exec(query, settingId)
	if err != nil {
		return fmt.Errorf("failed to delete setting: %s", err)
	}
	return nil
}
