package models

import (
	"time"

	"github.com/google/uuid"
)

type Mensaje struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FMatch     uuid.UUID `gorm:"type:uuid;not null" json:"f_match"`
	FRemitente uuid.UUID `gorm:"type:uuid;not null" json:"f_remitente"`
	Contenido  string    `gorm:"not null" json:"contenido"`
	Leido      bool      `gorm:"default:false" json:"leido"` // En Go, el valor por defecto de un bool es false
	FechaEnvio time.Time `gorm:"autoCreateTime" json:"fecha_envio"`
}