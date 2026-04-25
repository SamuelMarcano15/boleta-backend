package handlers

import (
	"net/http"
	"time"

	"github.com/SamuelMarcano15/boleta-backend/internal/models"
	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// BandejaItemDTO es lo que Flutter necesita para pintar la lista de chats
type BandejaItemDTO struct {
	MatchID       string    `json:"match_id"`
	OtroUsuarioID string    `json:"otro_usuario_id"`
	Nombre        string    `json:"nombre"`
	FechaMatch    time.Time `json:"fecha_match"`
	// Aquí en el futuro agregaremos "UltimoMensaje"
}

// ObtenerBandejaEntrada maneja GET /matches/:id_usuario
func ObtenerBandejaEntrada(c *gin.Context) {
	miID := c.Param("id_usuario")

	matches, err := repository.ObtenerMatchesPorUsuario(miID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la bandeja de entrada"})
		return
	}

	var bandeja []BandejaItemDTO

	// Recorremos los matches para armar la lista bonita para Flutter
	for _, m := range matches {
		// Descubrir quién es la "otra" persona en la relación
		otroID := m.FUsuario1.String()
		if otroID == miID {
			otroID = m.FUsuario2.String()
		}

		// Buscar el nombre de esa otra persona (Podríamos optimizar esto con un JOIN SQL más adelante)
		otroUsuario, _ := repository.ObtenerUsuarioPorID(otroID)
		nombreOtro := "Usuario de Boleta"
		if otroUsuario != nil {
			nombreOtro = otroUsuario.Nombre
		}

		bandeja = append(bandeja, BandejaItemDTO{
			MatchID:       m.ID.String(),
			OtroUsuarioID: otroID,
			Nombre:        nombreOtro,
			FechaMatch:    m.FechaMatch,
		})
	}

	// Si no hay matches, devolvemos un arreglo vacío en lugar de null
	if bandeja == nil {
		bandeja = make([]BandejaItemDTO, 0)
	}

	c.JSON(http.StatusOK, bandeja)
}

// ObtenerHistorialChat maneja GET /matches/:id_match/mensajes
func ObtenerHistorialChat(c *gin.Context) {
	matchID := c.Param("id_match")

	mensajes, err := repository.ObtenerMensajesPorMatch(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar el historial"})
		return
	}

	if mensajes == nil {
		mensajes = make([]models.Mensaje, 0)
	}

	c.JSON(http.StatusOK, mensajes)
}