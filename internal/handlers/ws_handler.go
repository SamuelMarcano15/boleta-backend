package handlers

import (
	"log"
	"net/http"

	"github.com/SamuelMarcano15/boleta-backend/internal/models"
	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/SamuelMarcano15/boleta-backend/internal/websockets"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// MensajeEntranteDTO define la estructura que esperamos recibir desde Flutter
type MensajeEntranteDTO struct {
	MatchID     string `json:"match_id"`
	ReceptorID  string `json:"receptor_id"` // A quién se lo estamos enviando
	Contenido   string `json:"contenido"`
}

// Configuración del Upgrader (Convierte HTTP a WebSocket)
var upgrader = websocket.Upgrader{
	// En producción, aquí se restringen los dominios. Para el MVP permitimos todos.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ConectarWebSocket maneja la ruta GET /ws/:id_usuario
func ConectarWebSocket(c *gin.Context) {
	miUsuarioID := c.Param("id_usuario")

	// 1. Convertir la petición HTTP en un túnel WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error al hacer el upgrade a WebSocket:", err)
		return
	}

	// 2. Registrar la conexión en nuestro Manager global
	websockets.ChatManager.AgregarCliente(miUsuarioID, conn)

	// Asegurarnos de limpiar la conexión cuando el usuario cierre la app o pierda internet
	defer func() {
		websockets.ChatManager.RemoverCliente(miUsuarioID)
		conn.Close()
	}()

	// 3. El Bucle Infinito: Escuchar los mensajes que manda este usuario
	for {
		var input MensajeEntranteDTO

		// Leer el JSON que viene por el túnel
		err := conn.ReadJSON(&input)
		if err != nil {
			log.Printf("Conexión cerrada o error de lectura: %v", err)
			break // Rompemos el ciclo y desconectamos al usuario
		}

		//Validar que el mensaje no sea basura o un "ping" vacío
		if input.MatchID == "" || input.Contenido == "" || input.ReceptorID == "" {
			log.Println("Mensaje incompleto recibido, ignorando...")
			continue // Saltamos este mensaje y volvemos a escuchar
		}

		// Convertir IDs a UUID (y verificar que no den error)
		remitenteUUID, err1 := uuid.Parse(miUsuarioID)
		matchUUID, err2 := uuid.Parse(input.MatchID)

		if err1 != nil || err2 != nil {
			log.Println("Error de formato UUID en el mensaje, ignorando...")
			continue
		}

		// 4. Guardar en Base de Datos (Persistencia)
		nuevoMensaje := models.Mensaje{
			FMatch:     matchUUID,
			FRemitente: remitenteUUID,
			Contenido:  input.Contenido,
		}
		
		if err := repository.GuardarMensaje(&nuevoMensaje); err != nil {
			log.Println("Error guardando mensaje en la BD:", err)
			continue
		}

		// 5. El Envío en Tiempo Real: Empujar el mensaje a la otra persona (Si está conectada)
		websockets.ChatManager.EnviarMensajeA(input.ReceptorID, nuevoMensaje)
	}
}