package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ManejoGlobalErrores es un "Recovery" personalizado. 
// Atrapa cualquier colapso del servidor y le avisa a Flutter de forma elegante.
func ManejoGlobalErrores() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Imprimimos el error real en la consola de Go para nosotros los desarrolladores
				log.Printf("[ERROR GRAVE] Servidor colapsó: %v\n", err)

				// Le respondemos a Flutter con nuestro formato estándar de Error 500
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"codigo":  http.StatusInternalServerError,
					"mensaje": "Ocurrió un error interno en el servidor de Boleta. Nuestros ingenieros ya fueron notificados.",
				})
			}
		}()

		// Dejamos que la petición continúe su camino normal
		c.Next()
	}
}