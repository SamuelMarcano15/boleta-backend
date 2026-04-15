package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ObtenerFeed maneja la ruta GET /feed/:id_usuario
func ObtenerFeed(c *gin.Context) {
	usuarioID := c.Param("id_usuario")

	// 1. Obtener al usuario actual para conocer sus preferencias
	usuarioActual, err := repository.ObtenerUsuarioPorID(usuarioID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario actual no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener perfil"})
		return
	}

	// 2. Validación de seguridad (Por si el usuario no ha completado su perfil)
	if usuarioActual.EstadoPreferido == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debes configurar un 'Estado Preferido' antes de buscar."})
		return
	}

	// 3. El Hack del Arquitecto: Traducir Edades a Fechas exactas
	hoy := time.Now()
	// Si el máximo son 35 años, buscamos a la persona que nació hace 35 años (El límite más viejo)
	fechaMin := hoy.AddDate(int(-usuarioActual.RangoEdadMax), 0, 0)
	// Si el mínimo son 21 años, buscamos a la persona que nació hace 21 años (El límite más joven)
	fechaMax := hoy.AddDate(int(-usuarioActual.RangoEdadMin), 0, 0)

	// 4. Ir a la base de datos por los candidatos
	candidatos, err := repository.ObtenerCandidatosFeed(
		usuarioActual.ID.String(),
		usuarioActual.BuscandoGenero,
		usuarioActual.BuscandoIntencion,
		*usuarioActual.EstadoPreferido,
		fechaMin,
		fechaMax,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando el algoritmo de discovery"})
		return
	}

	// 5. Devolver la lista a Flutter
	c.JSON(http.StatusOK, gin.H{
		"cantidad_resultados": len(candidatos),
		"candidatos":          candidatos,
	})
}