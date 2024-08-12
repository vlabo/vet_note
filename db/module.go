package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Procedure struct {
	ID        int64
	Type      string
	Date      string
	Details   string
	PatientID int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Patient struct {
	ID           int64
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
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
	Procedures   []Procedure
}

type SettingType string

const (
	PatientType   SettingType = "PatientType"
	ProcedureType SettingType = "ProcedureType"
)

type Setting struct {
	ID        int64
	Value     string
	Type      SettingType
	Idx       uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

var db *sql.DB

func InitializeDB(path string, _ bool) error {
	var err error
	db, err = sql.Open("sqlite", path)
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
        chip_id TEXT,
        weight REAL,
        castrated NUMERIC,
        last_modified TEXT,
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

	return nil
}

func GetPatient(id string) (patient Patient, err error) {
	query := `SELECT id, type, name, gender, birth_date, chip_id, weight, castrated, last_modified, note, owner, owner_phone, created_at, updated_at, deleted_at FROM patients WHERE id = ?`
	row := db.QueryRow(query, id)
	err = row.Scan(&patient.ID, &patient.Type, &patient.Name, &patient.Gender, &patient.BirthDate, &patient.ChipId, &patient.Weight, &patient.Castrated, &patient.LastModified, &patient.Note, &patient.Owner, &patient.OwnerPhone, &patient.CreatedAt, &patient.UpdatedAt, &patient.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("patient with Id %s not found", id)
		} else {
			err = fmt.Errorf("database error: %s", err)
		}
		return
	}

	// Load procedures
	query = `SELECT id, type, date, details, patient_id, created_at, updated_at, deleted_at FROM procedures WHERE patient_id = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var procedure Procedure
		err = rows.Scan(&procedure.ID, &procedure.Type, &procedure.Date, &procedure.Details, &procedure.PatientID, &procedure.CreatedAt, &procedure.UpdatedAt, &procedure.DeletedAt)
		if err != nil {
			return
		}
		patient.Procedures = append(patient.Procedures, procedure)
	}

	return
}

func GetProcedure(id string) (procedure Procedure, err error) {
	query := `SELECT id, type, date, details, patient_id, created_at, updated_at, deleted_at FROM procedures WHERE id = ?`
	row := db.QueryRow(query, id)
	err = row.Scan(&procedure.ID, &procedure.Type, &procedure.Date, &procedure.Details, &procedure.PatientID, &procedure.CreatedAt, &procedure.UpdatedAt, &procedure.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("procedure with Id %s not found", id)
		} else {
			err = fmt.Errorf("database error: %s", err)
		}
	}
	return
}

func GetPatientList() ([]Patient, error) {
	query := `SELECT id, type, name, chip_id, owner, owner_phone FROM patients`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var patient Patient
		err = rows.Scan(&patient.ID, &patient.Type, &patient.Name, &patient.ChipId, &patient.Owner, &patient.OwnerPhone)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func CreatePatient(patient *Patient) error {
	query := `INSERT INTO patients (type, name, gender, birth_date, chip_id, weight, castrated, last_modified, note, owner, owner_phone) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, patient.Type, patient.Name, patient.Gender, patient.BirthDate, patient.ChipId, patient.Weight, patient.Castrated, patient.LastModified, patient.Note, patient.Owner, patient.OwnerPhone)
	if err != nil {
		return fmt.Errorf("failed to create patient: %s", err)
	}
	// Get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %s", err)
	}
	patient.ID = id
	return nil
}

func UpdatePatient(id int64, changes *Patient) error {
	query := `UPDATE patients SET type = ?, name = ?, gender = ?, birth_date = ?, chip_id = ?, weight = ?, castrated = ?, last_modified = ?, note = ?, owner = ?, owner_phone = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, changes.Type, changes.Name, changes.Gender, changes.BirthDate, changes.ChipId, changes.Weight, changes.Castrated, changes.LastModified, changes.Note, changes.Owner, changes.OwnerPhone, id)
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

func CreateProcedure(procedure *Procedure) error {
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
	procedure.ID = id
	return nil
}

func UpdateProcedure(id int64, changes *Procedure) error {
	query := `UPDATE procedures SET type = ?, date = ?, details = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, changes.Type, changes.Date, changes.Details, id)
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

func GetPatientTypes() ([]Setting, error) {
	query := `SELECT id, value, type, idx, created_at, updated_at, deleted_at FROM settings WHERE type = ? ORDER BY idx`
	rows, err := db.Query(query, PatientType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []Setting
	for rows.Next() {
		var setting Setting
		err = rows.Scan(&setting.ID, &setting.Value, &setting.Type, &setting.Idx, &setting.CreatedAt, &setting.UpdatedAt, &setting.DeletedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func GetProcedureTypes() ([]Setting, error) {
	query := `SELECT id, value, type, idx, created_at, updated_at, deleted_at FROM settings WHERE type = ? ORDER BY idx`
	rows, err := db.Query(query, ProcedureType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []Setting
	for rows.Next() {
		var setting Setting
		err = rows.Scan(&setting.ID, &setting.Value, &setting.Type, &setting.Idx, &setting.CreatedAt, &setting.UpdatedAt, &setting.DeletedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func UpdateSetting(setting Setting) error {
	query := `INSERT INTO settings (value, type, idx) VALUES (?, ?, ?)`
	_, err := db.Exec(query, setting.Value, setting.Type, setting.Idx, setting.ID)
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
