package db

import (
	"fmt"

	"vet_note/db/models"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/shopspring/decimal"
)

type ViewProcedure struct {
	ID        omit.Val[int32]      `json:"id"`
	Type      omitnull.Val[string] `json:"type"`
	Date      omitnull.Val[string] `json:"date"`
	Details   omitnull.Val[string] `json:"details"`
	PatientID omitnull.Val[int32]  `json:"patientId"`
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
	ID         omit.Val[int32]               `json:"id"`
	Type       omitnull.Val[string]          `json:"type"`
	Name       omitnull.Val[string]          `json:"name"`
	Gender     omitnull.Val[string]          `json:"gender" tstype:"'unknown' | 'male' | 'female'"`
	Age        omitnull.Val[decimal.Decimal] `json:"age"`
	ChipID     omitnull.Val[string]          `json:"chipId"`
	Weight     omitnull.Val[float32]         `json:"weight"`
	Castrated  omitnull.Val[decimal.Decimal] `json:"castrated"`
	Note       omitnull.Val[string]          `json:"note"`
	Owner      omitnull.Val[string]          `json:"owner"`
	OwnerPhone omitnull.Val[string]          `json:"ownerPhone"`
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
	ID    omit.Val[int32]      `json:"id"`
	Type  omitnull.Val[string] `json:"type" tstype:"'PatientType' | 'ProcedureType'"`
	Value omitnull.Val[string] `json:"value"`
	Index omitnull.Val[int32]  `json:"index"`
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
