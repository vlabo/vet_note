package db

import (
	"fmt"

	"vet_note/db/models"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/shopspring/decimal"
)

type ViewProcedure struct {
	ID        omit.Val[int32]      `json:"id" tstype:"number"`
	Type      omitnull.Val[string] `json:"type" tstype:"string | null"`
	Date      omitnull.Val[string] `json:"date" tstype:"string | null"`
	Details   omitnull.Val[string] `json:"details" tstype:"string | null"`
	PatientID omitnull.Val[int32]  `json:"patientId" tstype:"number | null"`
}

func (vp ViewProcedure) AsSetter() models.ProcedureSetter {
	return models.ProcedureSetter{
		ID:        vp.ID,
		Type:      vp.Type,
		Date:      vp.Date,
		Details:   vp.Details,
		PatientID: vp.PatientID,
	}
}

func ViewProcedureFromModel(p *models.Procedure) ViewProcedure {
	return ViewProcedure{
		ID:        omit.From(p.ID),
		Type:      omitnull.FromNull(p.Type),
		Date:      omitnull.FromNull(p.Date),
		Details:   omitnull.FromNull(p.Details),
		PatientID: omitnull.FromNull(p.PatientID),
	}
}

type ViewPatient struct {
	ID         omit.Val[int32]               `json:"id" tstype:"number"`
	Type       omitnull.Val[string]          `json:"type" tstype:"string | null"`
	Name       omitnull.Val[string]          `json:"name" tstype:"string | null"`
	Gender     omitnull.Val[string]          `json:"gender" tstype:"'unknown' | 'male' | 'female'"`
	Age        omitnull.Val[decimal.Decimal] `json:"age" tstype:"number | null"`
	ChipID     omitnull.Val[string]          `json:"chipId" tstype:"number | null"`
	Weight     omitnull.Val[float32]         `json:"weight" tstype:"number | null"`
	Castrated  omitnull.Val[decimal.Decimal] `json:"castrated" tstype:"number | null"`
	Note       omitnull.Val[string]          `json:"note" tstype:"string | null"`
	Owner      omitnull.Val[string]          `json:"owner" tstype:"string | null"`
	OwnerPhone omitnull.Val[string]          `json:"ownerPhone" tstype:"string | null"`
	Procedures []ViewProcedure               `json:"procedures"`
}

func (vp ViewPatient) AsSetter() models.PatientSetter {
	return models.PatientSetter{
		ID:         vp.ID,
		Type:       vp.Type,
		Name:       vp.Name,
		Gender:     vp.Gender,
		Age:        vp.Age,
		ChipID:     vp.ChipID,
		Weight:     vp.Weight,
		Castrated:  vp.Castrated,
		Note:       vp.Note,
		Owner:      vp.Owner,
		OwnerPhone: vp.OwnerPhone,
	}
}

func ViewPatientFromModel(p *models.Patient) ViewPatient {
	vp := ViewPatient{
		ID:         omit.From(p.ID),
		Type:       omitnull.FromNull(p.Type),
		Name:       omitnull.FromNull(p.Name),
		Gender:     omitnull.FromNull(p.Gender),
		Age:        omitnull.FromNull(p.Age),
		ChipID:     omitnull.FromNull(p.ChipID),
		Weight:     omitnull.FromNull(p.Weight),
		Castrated:  omitnull.FromNull(p.Castrated),
		Note:       omitnull.FromNull(p.Note),
		Owner:      omitnull.FromNull(p.Owner),
		OwnerPhone: omitnull.FromNull(p.OwnerPhone),
		Procedures: make([]ViewProcedure, 0, 2),
	}
	return vp
}

type ViewSetting struct {
	ID    omit.Val[int32]      `json:"id" tstype:"number"`
	Type  omitnull.Val[string] `json:"type" tstype:"'PatientType' | 'ProcedureType'"`
	Value omitnull.Val[string] `json:"value" tstype:"string | null"`
	Index omitnull.Val[int32]  `json:"index" tstype:"number | null"`
}

func (vs ViewSetting) AsSetter() models.SettingSetter {
	return models.SettingSetter{
		ID:    vs.ID,
		Value: vs.Value,
		Type:  vs.Type,
		Idx:   vs.Index,
	}
}

func ViewSettingFromModel(s *models.Setting) ViewSetting {
	return ViewSetting{
		ID:    omit.From(s.ID),
		Type:  omitnull.FromNull(s.Type),
		Value: omitnull.FromNull(s.Value),
		Index: omitnull.FromNull(s.Idx),
	}
}

type ViewError struct {
	Error string `json:"error"`
}

func FmtError(format string, args ...any) ViewError {
	return ViewError{
		Error: fmt.Sprintf(format, args...),
	}
}
