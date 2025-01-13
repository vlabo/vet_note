package db

import (
	"fmt"
)

type ViewProcedure struct {
	ID        int32  `json:"id"`
	Type      string `json:"type"`
	Date      string `json:"date"`
	Details   string `json:"details"`
	PatientID int64  `json:"patientId"`
}

type ViewListPatient struct {
	ID     int32  `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	ChipId string `json:"chipId"`
	Owner  string `json:"owner"`
	Phone  string `json:"phone"`
}

type ViewPatient struct {
	ID         int32           `json:"id"`
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	Gender     string          `json:"gender" tstype:"'unknown' | 'male' | 'female'"`
	Age        string          `json:"age"`
	ChipID     string          `json:"chipId"`
	Weight     float32         `json:"weight"`
	Castrated  bool            `json:"castrated"`
	Note       string          `json:"note"`
	Owner      string          `json:"owner"`
	OwnerPhone string          `json:"ownerPhone"`
	Procedures []ViewProcedure `json:"procedures"`
}

type ViewSetting struct {
	ID    *int64      `json:"id"`
	Type  SettingType `json:"type" tstype:"'PatientType' | 'ProcedureType'"`
	Value string      `json:"value"`
	Index uint        `json:"index"`
}

type ViewError struct {
	Error string `json:"error"`
}

func FmtError(format string, args ...any) ViewError {
	return ViewError{
		Error: fmt.Sprintf(format, args...),
	}
}
