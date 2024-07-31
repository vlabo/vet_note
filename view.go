package main

import (
	"fmt"
	"strconv"

	"vet_note/db"

	"gorm.io/gorm"
)

type ViewProcedure struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Date      string `json:"date"`
	Details   string `json:"details"`
	PatientID string `json:"patientId"`
}

type ViewListPatient struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	ChipId string `json:"chipId"`
	Owner  string `json:"owner"`
	Phone  string `json:"phone"`
}

type ViewPatient struct {
	Id           string          `json:"id"`
	Type         string          `json:"type"`
	Name         string          `json:"name"`
	Gender       string          `json:"gender" tstype:"'unknown' | 'male' | 'female'"`
	BirthDate    string          `json:"birthDate"`
	ChipId       string          `json:"chipId"`
	Weight       float64         `json:"weight"`
	Castrated    bool            `json:"castrated"`
	LastModified string          `json:"lastModified"`
	Note         string          `json:"note"`
	Owner        string          `json:"owner"`
	OwnerPhone   string          `json:"ownerPhone"`
	Procedures   []ViewProcedure `json:"procedures"`
}

func (p *ViewPatient) asPatient() db.Patient {
	id, _ := strconv.ParseUint(p.Id, 10, 64)
	return db.Patient{
		Model: gorm.Model{
			ID: uint(id),
		},
		Type:         p.Type,
		Name:         p.Name,
		Gender:       p.Gender,
		BirthDate:    p.BirthDate,
		ChipId:       p.ChipId,
		Weight:       p.Weight,
		Castrated:    p.Castrated,
		LastModified: p.LastModified,
		Note:         p.Note,
		Owner:        p.Owner,
		OwnerPhone:   p.OwnerPhone,
	}
}

func (p *ViewProcedure) asProcedure() db.Procedure {
	id, _ := strconv.ParseUint(p.Id, 10, 64)
	return db.Procedure{
		Model: gorm.Model{
			ID: uint(id),
		},
		Type:    p.Type,
		Date:    p.Date,
		Details: p.Details,
	}
}

func procedureToViewProcedure(p db.Procedure) ViewProcedure {
	return ViewProcedure{
		Id:        strconv.FormatUint(uint64(p.ID), 10),
		Type:      p.Type,
		Date:      p.Date,
		Details:   p.Details,
		PatientID: strconv.FormatUint(uint64(p.PatientID), 10),
	}
}

func patientToListPatient(p db.Patient) ViewListPatient {
	return ViewListPatient{
		Id:     strconv.FormatUint(uint64(p.ID), 10),
		Type:   p.Type,
		Name:   p.Name,
		ChipId: p.ChipId,
		Owner:  p.Owner,
		Phone:  p.OwnerPhone,
	}
}

func patientToViewPatient(p db.Patient) ViewPatient {
	vp := ViewPatient{
		Id:           strconv.FormatUint(uint64(p.ID), 10),
		Type:         p.Type,
		Name:         p.Name,
		Gender:       p.Gender,
		BirthDate:    p.BirthDate,
		ChipId:       p.ChipId,
		Weight:       p.Weight,
		Castrated:    p.Castrated,
		LastModified: p.LastModified,
		Note:         p.Note,
		Owner:        p.Owner,
		OwnerPhone:   p.OwnerPhone,
		Procedures:   []ViewProcedure{},
	}

	for _, proc := range p.Procedures {
		vp.Procedures = append(vp.Procedures, procedureToViewProcedure(proc))
	}

	return vp
}

type ViewError struct {
	Error string `json:"error"`
}

func FmtError(format string, args ...any) ViewError {
	return ViewError{
		Error: fmt.Sprintf(format, args...),
	}
}
