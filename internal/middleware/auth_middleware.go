package middleware

import (
	"net/http"
	"strings"

	"github.com/SamuelMarcano15/boleta-backend/internal/security"
	"github.com/gin-gonic/gin"
)

// AuthRequerido intercepta la petición y exige un token válido
func AuthRequerido() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Buscar el token en los Headers de la petición
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado: Falta el token de autorización"})
			c.Abort() // Detiene la petición aquí mismo
			return
		}

		// 2. El estándar de la industria es enviar el token así: "Bearer eyJhbGci..."
		// Separamos la palabra "Bearer" del token real
		partes := strings.Split(authHeader, " ")
		if len(partes) != 2 || partes[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado: Formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := partes[1]

		// 3. Validar el token usando nuestro servicio de seguridad
		usuarioID, err := security.ValidarJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso denegado: Token inválido o expirado"})
			c.Abort()
			return
		}

		// 4. (Opcional pero muy pro): Guardar el ID en el contexto de Gin
		// Así, los handlers (ej. Swipe, Feed) saben exactamente quién está haciendo la petición
		// sin tener que mandar el ID en el JSON o en la URL.
		c.Set("usuario_id", usuarioID)

		// 5. Dejarlo pasar al Handler final
		c.Next()
	}
}