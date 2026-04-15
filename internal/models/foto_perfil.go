package models

import (
	"time"

	"github.com/google/uuid"
)

type FotoPerfil struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FUsuario    uuid.UUID `gorm:"type:uuid;not null" json:"f_usuario"`
	UrlFoto     string    `gorm:"not null" json:"url_foto"`
	OrdenVisual int16     `gorm:"default:1;not null" json:"orden_visual"` // 1 para principal, 2, 3...
	FechaSubida time.Time `gorm:"autoCreateTime" json:"fecha_subida"`
}

// TableName sobreescribe el nombre de la tabla por defecto de GORM.
// Sin esto, GORM buscaría una tabla llamada "foto_perfils".
func (FotoPerfil) TableName() string {
	return "fotos_perfil"
}