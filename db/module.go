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
	"github.com/shopspring/decimal"
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

func CreatePatient(patient models.PatientSetter) error {
	slog.Info("CreatePatient", "patient", patient)
	patient.CreatedAt = omitnull.From(time.Now())
	patient.UpdatedAt = omitnull.From(time.Now())
	_, err := models.Patients.Insert(&patient).Exec(context.Background(), db)
	return err
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
		settings = append(settings, ViewSetting{
			ID:    s.ID,
			Type:  PatientType,
			Value: s.Value.GetOrZero(),
			Index: s.Idx.GetOrZero(),
		})
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
		settings = append(settings, ViewSetting{
			ID:    s.ID,
			Type:  ProcedureType,
			Value: s.Value.GetOrZero(),
			Index: s.Idx.GetOrZero(),
		})
	}

	return settings, nil
}

func CreateSetting(setting ViewSetting) error {
	setter := models.SettingSetter{
		Value:     omitnull.From(setting.Value),
		Type:      omitnull.From(string(setting.Type)),
		Idx:       omitnull.From(setting.Index),
		CreatedAt: omitnull.From(time.Now()),
		UpdatedAt: omitnull.From(time.Now()),
	}
	_, err := models.Settings.Insert(&setter).Exec(context.Background(), db)

	return err
}

func UpdateSetting(setting ViewSetting) error {
	setter := models.SettingSetter{
		Value:     omitnull.From(setting.Value),
		Type:      omitnull.From(string(setting.Type)),
		Idx:       omitnull.From(setting.Index),
		CreatedAt: omitnull.From(time.Now()),
		UpdatedAt: omitnull.From(time.Now()),
	}
	_, err := models.Settings.Update(setter.UpdateMod(), models.UpdateWhere.Settings.ID.EQ(setting.ID)).Exec(context.Background(), db)
	return err
}

func DeleteSetting(settingId int32) error {
	_, err := models.Settings.Delete(models.DeleteWhere.Settings.ID.EQ(settingId)).Exec(context.Background(), db)
	return err
}
