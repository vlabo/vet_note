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
	"github.com/stephenafamo/bob"
	_ "modernc.org/sqlite"
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

	patient = ViewPatientFromModel(dbPatient)

	procedures, err := dbPatient.Procedures().All(context.Background(), db)
	if err != nil {
		return
	}

	for _, proc := range procedures {
		patient.Procedures = append(patient.Procedures, ViewProcedureFromModel(proc))
	}

	return
}

func GetProcedure(id string) (procedure ViewProcedure, err error) {
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}
	dbProcedure, err := models.FindProcedure(context.Background(), db, int32(idInt))
	procedure = ViewProcedureFromModel(dbProcedure)

	return
}

func GetPatientList() ([]ViewPatient, error) {
	dbPatients, err := models.Patients.Query().All(context.Background(), db)
	if err != nil {
		return nil, err
	}

	var patients []ViewPatient
	for _, p := range dbPatients {
		patient := ViewPatientFromModel(p)
		procedures, err := p.Procedures().All(context.Background(), db)
		if err != nil {
			continue
		}

		for _, proc := range procedures {
			patient.Procedures = append(patient.Procedures, ViewProcedureFromModel(proc))
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

func CreatePatient(patient models.PatientSetter) (int64, error) {
	slog.Info("CreatePatient", "patient", patient)
	patient.CreatedAt = omitnull.From(time.Now())
	patient.UpdatedAt = omitnull.From(time.Now())
	result, err := bob.Exec(context.Background(), db, models.Patients.Insert(&patient))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdatePatient(patient models.PatientSetter) error {
	slog.Info("UpdatePatient", "id", patient.ID.GetOrZero())
	patient.UpdatedAt = omitnull.From(time.Now())

	_, err := models.Patients.Update(patient.UpdateMod(), models.UpdateWhere.Patients.ID.EQ(patient.ID.GetOrZero())).Exec(context.Background(), db)
	return err
}

func DeletePatient(patientId int32) error {
	_, err := models.Patients.Delete(models.DeleteWhere.Patients.ID.EQ(patientId)).Exec(context.Background(), db)
	return err
}

func CreateProcedure(procedure models.ProcedureSetter) error {
	procedure.CreatedAt = omitnull.From(time.Now())
	procedure.UpdatedAt = omitnull.From(time.Now())
	_, err := models.Procedures.Insert(&procedure).Exec(context.Background(), db)
	return err
}

func UpdateProcedure(procedure models.ProcedureSetter) error {
	slog.Info("UpdatePatient", "id", procedure.ID.GetOrZero())
	procedure.UpdatedAt = omitnull.From(time.Now())

	_, err := models.Procedures.Update(procedure.UpdateMod(), models.UpdateWhere.Procedures.ID.EQ(procedure.ID.GetOrZero())).Exec(context.Background(), db)
	return err
}

func DeleteProcedure(procedureId int32) error {
	_, err := models.Procedures.Delete(models.DeleteWhere.Procedures.ID.EQ(procedureId)).Exec(context.Background(), db)
	return err
}

func GetPatientTypes() ([]ViewSetting, error) {
	dbSettings, err := models.Settings.Query(models.SelectWhere.Settings.Type.EQ("PatientType")).All(context.Background(), db)
	if err != nil {
		return nil, err
	}

	settings := make([]ViewSetting, 0, len(dbSettings))

	for _, s := range dbSettings {
		settings = append(settings, ViewSettingFromModel(s))
	}

	return settings, nil
}

func GetProcedureTypes() ([]ViewSetting, error) {
	dbSettings, err := models.Settings.Query(models.SelectWhere.Settings.Type.EQ("ProcedureType")).All(context.Background(), db)
	if err != nil {
		return nil, err
	}

	settings := make([]ViewSetting, 0, len(dbSettings))

	for _, s := range dbSettings {
		settings = append(settings, ViewSettingFromModel(s))
	}

	return settings, nil
}

func CreateSetting(setting ViewSetting) error {
	setter := setting.AsSetter()
	_, err := models.Settings.Insert(&setter).Exec(context.Background(), db)

	return err
}

func UpdateSetting(setting ViewSetting) error {
	setter := setting.AsSetter()
	_, err := models.Settings.Update(setter.UpdateMod(), models.UpdateWhere.Settings.ID.EQ(setting.ID.GetOrZero())).Exec(context.Background(), db)
	return err
}

func DeleteSetting(settingId int32) error {
	_, err := models.Settings.Delete(models.DeleteWhere.Settings.ID.EQ(settingId)).Exec(context.Background(), db)
	return err
}
