package handlers

import (
	"net/http"

	"github.com/SamuelMarcano15/boleta-backend/internal/models"
	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SwipeDTO es lo que nos envía Flutter cuando el usuario desliza
type SwipeDTO struct {
	FEvaluador string `json:"f_evaluador" binding:"required"` // Quien hace el swipe (Samuel)
	FEvaluado  string `json:"f_evaluado" binding:"required"`  // A quien están viendo (Maria)
	TipoAccion string `json:"tipo_accion" binding:"required,oneof=ACEPTAR RECHAZAR"` // La decisión
}

// ProcesarSwipe maneja la ruta POST /swipes
func ProcesarSwipe(c *gin.Context) {
	var input SwipeDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Convertir los strings a UUIDs
	evaluadorID, err1 := uuid.Parse(input.FEvaluador)
	evaluadoID, err2 := uuid.Parse(input.FEvaluado)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los IDs de usuario no son UUIDs válidos"})
		return
	}

	// No puedes darte swipe a ti mismo
	if evaluadorID == evaluadoID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No puedes evaluarte a ti mismo"})
		return
	}

	// Preparar el modelo
	nuevoSwipe := models.Swipe{
		FEvaluador: evaluadorID,
		FEvaluado:  evaluadoID,
		TipoAccion: input.TipoAccion,
	}

	// Mandar a la capa lógica (La transacción)
	huboMatch, err := repository.RegistrarSwipe(&nuevoSwipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el swipe"})
		return
	}

	// Responder a Flutter
	respuesta := gin.H{
		"mensaje":          "Acción registrada con éxito",
		"match_encontrado": huboMatch,
	}

	if huboMatch {
		respuesta["mensaje"] = "¡Coronaste! Tienes un nuevo Match."
	}

	c.JSON(http.StatusCreated, respuesta)
}