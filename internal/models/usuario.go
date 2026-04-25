package models

import (
	"time"

	"github.com/google/uuid"
)

// Usuario mapea exactamente la tabla 'usuarios' de PostgreSQL.
type Usuario struct {
	// Usamos uuid.UUID en lugar de uint para el ID
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	
	Nombre            string    `gorm:"not null" json:"nombre"`
	FechaNacimiento   time.Time `gorm:"type:date;not null" json:"fecha_nacimiento"`
	
	// Los ENUMs de Postgres se mapean como strings en Go, pero le decimos a GORM su tipo real
	Genero            string    `gorm:"type:genero_tipo;not null" json:"genero"` 
	
	// PUNTEROS (*): Estos campos pueden ser NULL en la base de datos
	OrientacionSexual *string   `json:"orientacion_sexual"`
	Biografia         *string   `json:"biografia"`
	
	// Ubicación Real
	PaisOrigen        string    `gorm:"not null" json:"pais_origen"`
	EstadoProvincia   *string   `json:"estado_provincia"`
	Latitud           *float64  `json:"latitud"`
	Longitud          *float64  `json:"longitud"`
	Telefono          *string   `json:"telefono"`
	
	// Preferencias del Algoritmo
	BuscandoGenero    string    `gorm:"type:genero_tipo;not null" json:"buscando_genero"`
	RangoEdadMin      int16     `gorm:"default:18" json:"rango_edad_min"`
	RangoEdadMax      int16     `gorm:"default:99" json:"rango_edad_max"`
	DistanciaMaximaKm int16     `gorm:"default:50" json:"distancia_maxima_km"`
	BuscandoIntencion string    `gorm:"type:intencion_tipo;not null" json:"buscando_intencion"`
	EstadoPreferido   *string   `json:"estado_preferido"`
	
	// Fecha de registro automática
	FechaRegistro     time.Time `gorm:"autoCreateTime" json:"fecha_registro"`

	Fotos []FotoPerfil `gorm:"foreignKey:FUsuario" json:"fotos"`

	Correo string `gorm:"uniqueIndex;not null" json:"correo"`
	Clave  string `gorm:"not null" json:"-"`
}