package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/SamuelMarcano15/boleta-backend/internal/handlers"
)

// SetupRouter configura todos los endpoints de la API
func SetupRouter(r *gin.Engine) {
	
	r.Static("/uploads", "./uploads")

	// Grupo de rutas para la API v1
	v1 := r.Group("/api/v1")
	{
		// Modulo: Gestión de Perfil
		usuarios := v1.Group("/usuarios")
		{
			usuarios.POST("/", handlers.CrearUsuario)
			usuarios.GET("/:id", handlers.ObtenerUsuario)
			usuarios.PUT("/:id/preferencias", handlers.ActualizarPreferencias)
			usuarios.POST("/:id/fotos", handlers.SubirFoto)
		}

		// Modulo: El Motor "Boleta" (Discovery & Match)
		feed := v1.Group("/feed")
		{
			// GET /api/v1/feed/:id_usuario
			feed.GET("/:id_usuario", handlers.ObtenerFeed)
		}

		// Modulo: Interacciones Sociales
		swipes := v1.Group("/swipes")
		{
			swipes.POST("/", handlers.ProcesarSwipe)
		}
	}
}