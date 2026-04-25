package repository

import (
	"github.com/SamuelMarcano15/boleta-backend/internal/models"
)

// ObtenerMatchesPorUsuario busca todos los matches donde participe el usuario (ya sea como 1 o como 2)
func ObtenerMatchesPorUsuario(usuarioID string) ([]models.Match, error) {
	var matches []models.Match
	// Usamos OR porque el UUID de Samuel podría estar en la columna 1 o en la 2 dependiendo del orden lexicográfico
	result := DB.Where("f_usuario_1 = ? OR f_usuario_2 = ?", usuarioID, usuarioID).
		Order("fecha_match desc").
		Find(&matches)
	return matches, result.Error
}

// ObtenerMensajesPorMatch trae el historial de una conversación
func ObtenerMensajesPorMatch(matchID string) ([]models.Mensaje, error) {
	var mensajes []models.Mensaje
	// Ordenamos ascendente para que el chat se lea de arriba (viejo) hacia abajo (nuevo)
	result := DB.Where("f_match = ?", matchID).
		Order("fecha_envio asc").
		Find(&mensajes)
	return mensajes, result.Error
}

// GuardarMensaje inserta un nuevo mensaje en la base de datos
func GuardarMensaje(mensaje *models.Mensaje) error {
	return DB.Create(mensaje).Error
}