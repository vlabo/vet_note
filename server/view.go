package main

import "github.com/google/uuid"

type ViewProcedure struct {
	Id      uuid.UUID `json:"id"`
	Type    string    `json:"type"`
	Date    string    `json:"date"`
	Details string    `json:"details"`
}

type ViewListPatient struct {
	Id     uuid.UUID `json:"id"`
	Type   string    `json:"type"`
	Name   string    `json:"name"`
	ChipId string    `json:"chipId"`
	Owner  string    `json:"owner"`
	Phone  string    `json:"phone"`
}

type ViewPatient struct {
	Id           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;"`
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
	return Patient{
		Id:           p.Id,
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
