package repository

import (
	"errors"

	"github.com/SamuelMarcano15/boleta-backend/internal/models"
	"gorm.io/gorm"
)

// RegistrarSwipe maneja la lógica de evaluar un perfil y verificar si hay match mutuo
func RegistrarSwipe(swipe *models.Swipe) (bool, error) {
	huboMatch := false

	// Iniciamos una Transacción: Todo este bloque se ejecuta como una operación única y segura
	err := DB.Transaction(func(tx *gorm.DB) error {

		// 1. Guardar la acción (El Swipe) de Samuel a Maria
		if err := tx.Create(swipe).Error; err != nil {
			return err
		}

		// 2. Si la acción fue un "RECHAZAR", terminamos aquí. No hay posibilidad de match.
		if swipe.TipoAccion == "RECHAZAR" {
			return nil
		}

		// 3. El Match Maker: Si fue "ACEPTAR", buscamos si Maria ya había "ACEPTADO" a Samuel antes
		var swipePrevio models.Swipe
		resultado := tx.Where("f_evaluador = ? AND f_evaluado = ? AND tipo_accion = ?", 
			swipe.FEvaluado, swipe.FEvaluador, "ACEPTAR").First(&swipePrevio)

		// Si no hay error, significa que SI encontró un swipe de vuelta -> ¡HAY MATCH!
		if resultado.Error == nil {
			huboMatch = true

			// 4. El Truco del Arquitecto: Ordenar los UUIDs lexicográficamente
			// Esto garantiza que cumpla con el CHECK de Postgres (usuario1 < usuario2)
			u1 := swipe.FEvaluador
			u2 := swipe.FEvaluado

			if u1.String() > u2.String() {
				u1 = swipe.FEvaluado
				u2 = swipe.FEvaluador
			}

			// Creamos el registro oficial del Match ya ordenado
			nuevoMatch := models.Match{
				FUsuario1: u1,
				FUsuario2: u2,
			}

			if err := tx.Create(&nuevoMatch).Error; err != nil {
				return err // Si esto falla, el Swipe original también se cancela (Rollback)
			}
		} else if !errors.Is(resultado.Error, gorm.ErrRecordNotFound) {
			// Si ocurrió un error grave de base de datos distinto a "no lo encontré"
			return resultado.Error
		}

		// Todo salió bien, confirmamos los cambios (Commit)
		return nil
	})

	return huboMatch, err
}