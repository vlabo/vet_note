package main

import (
	"fmt"

	"vet_note/db"
)

type ViewProcedure struct {
	ID        *int64  `json:"id"`
	Type      *string `json:"type"`
	Date      *string `json:"date"`
	Details   *string `json:"details"`
	PatientID int64   `json:"patientId"`
}

func (vp *ViewProcedure) asProcedure() db.Procedure {
	procedure := db.Procedure{}
	procedure.PatientID = vp.PatientID

	if vp.Type != nil {
		procedure.Type = *vp.Type
	}
	if vp.Date != nil {
		procedure.Date = *vp.Date
	}
	if vp.Details != nil {
		procedure.Details = *vp.Details
	}

	return procedure
}

func (vp *ViewProcedure) fromProcedure(procedure db.Procedure) {
	vp.ID = &procedure.ID
	vp.PatientID = procedure.PatientID

	vp.Type = &procedure.Type
	vp.Date = &procedure.Date
	vp.Details = &procedure.Details
}

type ViewListPatient struct {
	ID     int64  `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	ChipId string `json:"chipId"`
	Owner  string `json:"owner"`
	Phone  string `json:"phone"`
}

func (vlp *ViewListPatient) fromPatient(patient db.Patient) {
	vlp.ID = patient.ID

	vlp.Type = patient.Type
	vlp.Name = patient.Name
	vlp.ChipId = patient.ChipId
	vlp.Owner = patient.Owner
	vlp.Phone = patient.OwnerPhone
}

type ViewPatient struct {
	ID           *int64           `json:"id"`
	Type         *string          `json:"type"`
	Name         *string          `json:"name"`
	Gender       *string          `json:"gender" tstype:"'unknown' | 'male' | 'female'"`
	BirthDate    *string          `json:"birthDate"`
	ChipId       *string          `json:"chipId"`
	Weight       *float64         `json:"weight"`
	Castrated    *bool            `json:"castrated"`
	LastModified *string          `json:"lastModified"`
	Note         *string          `json:"note"`
	Owner        *string          `json:"owner"`
	OwnerPhone   *string          `json:"ownerPhone"`
	Procedures   *[]ViewProcedure `json:"procedures"`
}

func (vp *ViewPatient) asPatient() db.Patient {
	patient := db.Patient{}

	if vp.Type != nil {
		patient.Type = *vp.Type
	}
	if vp.Name != nil {
		patient.Name = *vp.Name
	}
	if vp.Gender != nil {
		patient.Gender = *vp.Gender
	}
	if vp.BirthDate != nil {
		patient.BirthDate = *vp.BirthDate
	}
	if vp.ChipId != nil {
		patient.ChipId = *vp.ChipId
	}
	if vp.Weight != nil {
		patient.Weight = *vp.Weight
	}
	if vp.Castrated != nil {
		patient.Castrated = *vp.Castrated
	}
	if vp.LastModified != nil {
		patient.LastModified = *vp.LastModified
	}
	if vp.Note != nil {
		patient.Note = *vp.Note
	}
	if vp.Owner != nil {
		patient.Owner = *vp.Owner
	}
	if vp.OwnerPhone != nil {
		patient.OwnerPhone = *vp.OwnerPhone
	}

	// if vp.Procedures != nil {
	// 	for _, vpProc := range *vp.Procedures {
	// 		patient.Procedures = append(patient.Procedures, vpProc.asProcedure())
	// 	}
	// }

	return patient
}

func (vp *ViewPatient) fromPatient(patient db.Patient) {
	vp.ID = &patient.ID

	vp.Type = &patient.Type
	vp.Name = &patient.Name
	vp.Gender = &patient.Gender
	vp.BirthDate = &patient.BirthDate
	vp.ChipId = &patient.ChipId
	vp.Weight = &patient.Weight
	vp.Castrated = &patient.Castrated
	vp.LastModified = &patient.LastModified
	vp.Note = &patient.Note
	vp.Owner = &patient.Owner
	vp.OwnerPhone = &patient.OwnerPhone

	if len(patient.Procedures) > 0 {
		var viewProcedures []ViewProcedure
		for _, proc := range patient.Procedures {
			vpProc := ViewProcedure{}
			vpProc.fromProcedure(proc)
			viewProcedures = append(viewProcedures, vpProc)
		}
		vp.Procedures = &viewProcedures
	}
}

type ViewSetting struct {
	ID    *int64         `json:"id"`
	Type  db.SettingType `json:"type" tstype:"'PatientType' | 'ProcedureType'"`
	Value string         `json:"value"`
	Index uint           `json:"index"`
}

func (vs *ViewSetting) asSetting() db.Setting {
	setting := db.Setting{}
	if vs.ID != nil {
		setting.ID = *vs.ID
	}
	setting.Type = vs.Type
	setting.Value = vs.Value
	setting.Idx = vs.Index
	return setting
}

func (vs *ViewSetting) fromSetting(setting db.Setting) {
	vs.ID = &setting.ID
	vs.Type = setting.Type
	vs.Value = setting.Value
	vs.Index = setting.Idx
}

type ViewError struct {
	Error string `json:"error"`
}

func FmtError(format string, args ...any) ViewError {
	return ViewError{
		Error: fmt.Sprintf(format, args...),
	}
}
