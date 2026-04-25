package routes

import (
	"net/http"
	"github.com/SamuelMarcano15/boleta-backend/internal/handlers"
	"github.com/SamuelMarcano15/boleta-backend/internal/middleware" 
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	// 1. ZONA DE PROTECCIÓN GLOBAL
	// Añadimos nuestro atrapador de errores a nivel global
	r.Use(middleware.ManejoGlobalErrores())

	// 2. MANEJO DE RUTAS NO ENCONTRADAS (404)
	// Si alguien pide una URL que no programamos, Gin usa esto por defecto
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"codigo":  http.StatusNotFound,
			"mensaje": "La ruta que buscas no existe en Boleta. Verifica la URL.",
		})
	})
	
	// ¡OJO! Mantenemos tu ruta estática para que las fotos sigan funcionando
	r.Static("/uploads", "./uploads")

	v1 := r.Group("/api/v1")
	{
		// ==========================================
		// 🟢 ZONA PÚBLICA (No requiere Token)
		// ==========================================
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
		}

		usuariosPublicos := v1.Group("/usuarios")
		{
			usuariosPublicos.POST("/", handlers.CrearUsuario) // Registro
		}

		// ==========================================
		// 🔴 ZONA PROTEGIDA (Requiere Token)
		// ==========================================
		protegido := v1.Group("/")
		protegido.Use(middleware.AuthRequerido()) // <-- ¡El Cadenero se pone aquí!
		{
			// Todo lo que esté dentro de este bloque, pedirá Token obligatoriamente

			usuariosPrivados := protegido.Group("/usuarios")
			{
				usuariosPrivados.GET("/:id", handlers.ObtenerUsuario)
				usuariosPrivados.PUT("/:id/preferencias", handlers.ActualizarPreferencias)
				usuariosPrivados.POST("/:id/fotos", handlers.SubirFoto)
			}

			feed := protegido.Group("/feed")
			{
				feed.GET("/:id_usuario", handlers.ObtenerFeed)
			}

			swipes := protegido.Group("/swipes")
			{
				swipes.POST("/", handlers.ProcesarSwipe)
			}

			chat := protegido.Group("/matches")
			{
				chat.GET("/usuario/:id_usuario", handlers.ObtenerBandejaEntrada)
				chat.GET("/:id_match/mensajes", handlers.ObtenerHistorialChat)
			}
		}

		// Nota: El WebSocket a veces requiere una validación de token distinta 
		// (por query param en vez de header), por ahora lo dejamos libre
		v1.GET("/ws/:id_usuario", handlers.ConectarWebSocket)
	}
}