package websockets

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Manager mantiene la lista de clientes conectados
type Manager struct {
	// Usamos un RWMutex para evitar que dos personas se conecten/desconecten exactamente 
	// al mismo milisegundo y rompan el mapa en memoria (Thread Safety).
	sync.RWMutex
	// Un mapa donde la llave es el ID del Usuario y el valor es su conexión de red
	Clients map[string]*websocket.Conn
}

// Creamos una instancia global de nuestro manager
var ChatManager = &Manager{
	Clients: make(map[string]*websocket.Conn),
}

// AgregarCliente registra una nueva conexión
func (m *Manager) AgregarCliente(userID string, conn *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	m.Clients[userID] = conn
	log.Printf("Usuario %s conectado al WebSocket", userID)
}

// RemoverCliente elimina una conexión cuando alguien cierra la app
func (m *Manager) RemoverCliente(userID string) {
	m.Lock()
	defer m.Unlock()
	delete(m.Clients, userID)
	log.Printf("Usuario %s desconectado del WebSocket", userID)
}

// EnviarMensajeA busca si un usuario está conectado y le empuja el mensaje por el túnel
func (m *Manager) EnviarMensajeA(userID string, mensajeJSON interface{}) {
	m.RLock()
	defer m.RUnlock()

	conn, existe := m.Clients[userID]
	if existe {
		// Si está conectado, le escribimos el JSON directamente al túnel
		err := conn.WriteJSON(mensajeJSON)
		if err != nil {
			log.Printf("Error enviando mensaje a %s: %v", userID, err)
			conn.Close()
		}
	}
	// Si no existe (está offline), no hacemos nada. 
	// Lo verá en su historial cuando vuelva a abrir la app.
}