package domain

import (
	"time"
)

type Secret struct {
	Id           string
	Title        string
	UpdatedAt    time.Time
	Certificates []Certificate
}
