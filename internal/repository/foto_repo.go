package repository

import (
	"github.com/SamuelMarcano15/boleta-backend/internal/models"
)

// GuardarFoto inserta un nuevo registro de foto en la base de datos
func GuardarFoto(foto *models.FotoPerfil) error {
	result := DB.Create(foto)
	return result.Error
}