package model

import "github.com/google/uuid"

type Patient struct {
	Id    uuid.UUID
	Name  string
	Owner string
}
