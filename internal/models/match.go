package models

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	// Le decimos explícitamente a GORM cómo se llaman las columnas en Postgres
	FUsuario1  uuid.UUID `gorm:"column:f_usuario_1;type:uuid;not null" json:"f_usuario_1"`
	FUsuario2  uuid.UUID `gorm:"column:f_usuario_2;type:uuid;not null" json:"f_usuario_2"`
	FechaMatch time.Time `gorm:"column:fecha_match;autoCreateTime" json:"fecha_match"`
}