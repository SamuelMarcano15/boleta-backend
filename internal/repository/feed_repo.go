package repository

import (
	"time"

	"github.com/SamuelMarcano15/boleta-backend/internal/models"
)

// ObtenerCandidatosFeed busca usuarios que encajen con las preferencias y que no hayan sido evaluados
func ObtenerCandidatosFeed(
	usuarioID string,
	buscandoGenero string,
	intencion string,
	estado string,
	fechaNacimientoMin time.Time, // Fecha del más viejo permitido
	fechaNacimientoMax time.Time, // Fecha del más joven permitido
) ([]models.Usuario, error) {
	
	var candidatos []models.Usuario

	// 1. La subconsulta (Subquery): "Tráeme los IDs de los usuarios que YO ya evalué (Like o Dislike)"
	subQuery := DB.Table("swipes").Select("f_evaluado").Where("f_evaluador = ?", usuarioID)

	// 2. La consulta principal
	result := DB.Preload("Fotos"). // Eager loading: Traer las fotos de los candidatos
		Where("id != ?", usuarioID). // Obvio: No me muestres a mí mismo
		Where("genero = ?", buscandoGenero).
		Where("buscando_intencion = ?", intencion).
		Where("estado_provincia = ?", estado). // El filtro estricto "Boleta"
		Where("fecha_nacimiento BETWEEN ? AND ?", fechaNacimientoMin, fechaNacimientoMax). // Optimización de índice
		Where("id NOT IN (?)", subQuery). // El Anti-Join: Excluir a los ya evaluados
		Limit(20). // Paginación de seguridad: Solo devolver 20 a la vez para no saturar a Flutter
		Find(&candidatos)

	return candidatos, result.Error
}