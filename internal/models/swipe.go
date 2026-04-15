package models

import (
	"time"

	"github.com/google/uuid"
)

type Swipe struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FEvaluador  uuid.UUID `gorm:"type:uuid;not null" json:"f_evaluador"` // Quien da el Like/Dislike
	FEvaluado   uuid.UUID `gorm:"type:uuid;not null" json:"f_evaluado"`  // Quien lo recibe
	TipoAccion  string    `gorm:"type:accion_swipe;not null" json:"tipo_accion"` // 'LIKE' o 'DISLIKE'
	FechaAccion time.Time `gorm:"autoCreateTime" json:"fecha_accion"`
}