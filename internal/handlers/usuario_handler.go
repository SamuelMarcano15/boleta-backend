package handlers

import (
	"net/http"
	"time"
	"errors"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/SamuelMarcano15/boleta-backend/internal/models"
	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/SamuelMarcano15/boleta-backend/internal/security"
)

// RegistroUsuarioDTO define exactamente que campos esperamos del frontend.

type RegistroUsuarioDTO struct{
	Nombre             string   `json:"nombre" binding:"required"`
	FechaNacimiento    string   `json:"fecha_nacimiento" binding:"required"` // Formato esperado: YYYY-MM-DD
	Genero             string   `json:"genero" binding:"required,oneof=M F Otro Todos"`
	PaisOrigen         string   `json:"pais_origen" binding:"required,len=2"` // Ej: VE
	BuscandoGenero     string   `json:"buscando_genero" binding:"required,oneof=M F Otro Todos"`
	BuscandoIntencion  string   `json:"buscando_intencion" binding:"required,oneof=Citas Juego Amistad 'Solo Tops'"`
	
	// Opcionales
	OrientacionSexual  *string  `json:"orientacion_sexual"`
	EstadoProvincia    *string  `json:"estado_provincia"`
	Telefono           *string  `json:"telefono"`

	Correo string `json:"correo" binding:"required,email"`
	Clave  string `json:"clave" binding:"required,min=6"`
}

// PreferenciasDTO define qué campos del algoritmo se pueden modificar
type PreferenciasDTO struct {
	BuscandoGenero    *string `json:"buscando_genero" binding:"omitempty,oneof=M F Otro Todos"`
	RangoEdadMin      *int16  `json:"rango_edad_min" binding:"omitempty,min=18"`
	RangoEdadMax      *int16  `json:"rango_edad_max" binding:"omitempty,max=99,gtefield=RangoEdadMin"`
	DistanciaMaximaKm *int16  `json:"distancia_maxima_km" binding:"omitempty,min=1"`
	BuscandoIntencion *string `json:"buscando_intencion" binding:"omitempty"`
	EstadoPreferido   *string `json:"estado_preferido"`
}

// CrearUsuario maneja la ruta POST /usuarios

func CrearUsuario(c *gin.Context){
	var input RegistroUsuarioDTO

	//1. Validar el JSON entrante.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// 2. Parsear la fecha de nacimiento (De String a time.Time)
	// En Go, el formato de referencia para parsear fechas siempre es "2006-01-02"
	fechaNac, err := time.Parse("2006-01-02", input.FechaNacimiento)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha invalido. Usa YYYY-MM-DD"})
		return
	}

	hash, err := security.HashearClave(input.Clave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	

	// 3. Mapear el DTO a nuestro Modelo Real
	nuevoUsuario := models.Usuario{
		Nombre:            input.Nombre,
		FechaNacimiento:   fechaNac,
		Genero:            input.Genero,
		OrientacionSexual: input.OrientacionSexual,
		PaisOrigen:        input.PaisOrigen,
		EstadoProvincia:   input.EstadoProvincia,
		Telefono:          input.Telefono,
		BuscandoGenero:    input.BuscandoGenero,
		BuscandoIntencion: input.BuscandoIntencion,
		Correo:            input.Correo,
		Clave:             hash,
		// Los valores por defecto (como rangos de edad) los pondrá Postgres/GORM
	}

	// 4. Guardar en la Base de Datos
	if err := repository.CrearUsuario(&nuevoUsuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el usuario"})
		return
	}

	// 5. Responder al frontend con código 201 (Created) y el ID generado
	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Usuario creado con éxito",
		"id":      nuevoUsuario.ID,
	})

}

// ActualizarPreferencias maneja la ruta PUT /usuarios/:id/preferencias
func ActualizarPreferencias(c *gin.Context) {
	// Obtener el ID dinámico de la URL
	idUsuario := c.Param("id")

	var input PreferenciasDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Construir el mapa de actualizaciones solo con los campos que no son nulos
	updates := make(map[string]interface{})

	if input.BuscandoGenero != nil {
		updates["buscando_genero"] = *input.BuscandoGenero
	}
	if input.RangoEdadMin != nil {
		updates["rango_edad_min"] = *input.RangoEdadMin
	}
	if input.RangoEdadMax != nil {
		updates["rango_edad_max"] = *input.RangoEdadMax
	}
	if input.DistanciaMaximaKm != nil {
		updates["distancia_maxima_km"] = *input.DistanciaMaximaKm
	}
	if input.BuscandoIntencion != nil {
		updates["buscando_intencion"] = *input.BuscandoIntencion
	}
	if input.EstadoPreferido != nil {
		updates["estado_preferido"] = *input.EstadoPreferido
	}

	// Si el JSON vino vacío
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se enviaron datos para actualizar"})
		return
	}

	// Enviar al repositorio
	if err := repository.ActualizarPreferencias(idUsuario, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar preferencias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Preferencias actualizadas con éxito",
	})
}

// ObtenerUsuario maneja la ruta GET /usuarios/:id
func ObtenerUsuario(c *gin.Context) {
	id := c.Param("id")

	usuario, err := repository.ObtenerUsuarioPorID(id)
	if err != nil {
		// Verificamos si el error es porque el registro no existe
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario"})
		return
	}

	c.JSON(http.StatusOK, usuario)
}