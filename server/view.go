package main

import (
	"strconv"

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

func (p *ViewPatient) asPatient() Patient {
	id, _ := strconv.ParseUint(p.Id, 10, 64)
	return Patient{
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

func (p *ViewProcedure) asProcedure() Procedure {
	id, _ := strconv.ParseUint(p.Id, 10, 64)
	return Procedure{
		Model: gorm.Model{
			ID: uint(id),
		},
		Type:    p.Type,
		Date:    p.Date,
		Details: p.Details,
	}
}
