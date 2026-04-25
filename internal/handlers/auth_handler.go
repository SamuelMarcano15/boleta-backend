package handlers

import (
	"net/http"

	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/SamuelMarcano15/boleta-backend/internal/security"
	"github.com/gin-gonic/gin"
)

type LoginDTO struct {
	Correo string `json:"correo" binding:"required,email"`
	Clave  string `json:"clave" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 1. Buscar al usuario por correo (Necesitarás crear esta función en usuario_repo.go)
	usuario, err := repository.ObtenerUsuarioPorCorreo(input.Correo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// 2. Verificar que la clave coincida
	if coinciden := security.VerificarClave(input.Clave, usuario.Clave); !coinciden {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// 3. Generar el Token JWT
	token, err := security.GenerarJWT(usuario.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	// 4. Devolver el token a Flutter
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Login exitoso",
		"token":   token,
		"usuario": usuario, // Recordatorio: La clave está oculta gracias al json:"-"
	})
}