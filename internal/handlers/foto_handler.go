package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SamuelMarcano15/boleta-backend/internal/models"
	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubirFoto maneja la ruta POST /usuarios/:id/fotos
func SubirFoto(c *gin.Context) {
	idUsuarioStr := c.Param("id")

	// 1. Convertir el ID de la URL (string) al tipo UUID nativo
	usuarioID, err := uuid.Parse(idUsuarioStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	// 2. Extraer el archivo del formulario (la clave 'foto' es la que usará Flutter)
	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró ninguna foto en la petición"})
		return
	}

	// 3. Generar un nombre único para el archivo y asegurar que la carpeta exista
	extension := filepath.Ext(file.Filename) // Ej: .jpg, .png
	nuevoNombre := uuid.New().String() + extension
	rutaCarpeta := "uploads/fotos"
	rutaFinal := filepath.Join(rutaCarpeta, nuevoNombre)

	// Crear la carpeta física si no existe (Permisos 0755 es el estándar seguro en Linux)
	os.MkdirAll(rutaCarpeta, 0755)

	// 4. Guardar el archivo en el disco duro del servidor
	if err := c.SaveUploadedFile(file, rutaFinal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la foto en el servidor"})
		return
	}

	// 5. Preparar la URL pública que devolveremos (Simulando lo que haría Cloudflare R2 a futuro)
	// Como estamos en local, guardamos la ruta relativa. En prod, podríamos inyectar el dominio base.
	urlPublica := fmt.Sprintf("/%s", filepath.ToSlash(rutaFinal)) // Resultado: /uploads/fotos/uuid.jpg

	// 6. Guardar en PostgreSQL
	nuevaFoto := models.FotoPerfil{
		FUsuario:    usuarioID,
		UrlFoto:     urlPublica,
		OrdenVisual: 1, // Por ahora todas serán 1. Luego puedes hacer un count para saber el orden.
	}

	if err := repository.GuardarFoto(&nuevaFoto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar en base de datos"})
		return
	}

	// 7. Éxito
	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Foto subida con éxito",
		"foto":    nuevaFoto,
	})
}