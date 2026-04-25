package repository

import (
	"github.com/SamuelMarcano15/boleta-backend/internal/models"
)

// CrearUsuario inserta un nuevo usuario en la base de datos

func CrearUsuario(usuario *models.Usuario) error {
	// GORM automáticamente omitirá el ID y dejará que Postgres genere el UUID,
	// y luego GORM llenará el struct con el ID generado y las fechas de creación.
	result := DB.Create(usuario)
	return result.Error
}

// ActualizarPreferencias modifica solo los campos indicados de un usuario
func ActualizarPreferencias(id string, preferencias map[string]interface{}) error {
	// GORM buscará al usuario por ID y actualizará solo las claves del mapa.
	result := DB.Model(&models.Usuario{}).Where("id = ?", id).Updates(preferencias)
	return result.Error
}

// ObtenerUsuarioPorID busca un usuario y carga sus fotos relacionadas
func ObtenerUsuarioPorID(id string) (*models.Usuario, error) {
	var usuario models.Usuario
	
	// .Preload("Fotos") busca las fotos que pertenecen a este usuario
	// .First busca el primer registro que coincida con el ID
	result := DB.Preload("Fotos").First(&usuario, "id = ?", id)
	
	return &usuario, result.Error
}

func ObtenerUsuarioPorCorreo(correo string) (*models.Usuario, error) {
	var usuario models.Usuario
	
	// GORM hace un: SELECT * FROM usuarios WHERE correo = '...' LIMIT 1
	result := DB.Where("correo = ?", correo).First(&usuario)
	
	return &usuario, result.Error
}

