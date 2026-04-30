package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitor mantiene el limitador de tasa y la última vez que fue visto.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	// Para zonas públicas (identificado por IP)
	publicVisitors = &sync.Map{}
	// Para zonas protegidas (identificado por User ID)
	protectedVisitors = &sync.Map{}
)

func init() {
	// Iniciar las goroutines de auto-limpieza en background
	go iniciarLimpieza(publicVisitors, 5*time.Minute, 10*time.Minute)
	go iniciarLimpieza(protectedVisitors, 5*time.Minute, 10*time.Minute)
}

// RateLimiterConfig define la configuración del limitador
type RateLimiterConfig struct {
	Rate     rate.Limit
	Burst    int
	KeyFunc  func(c *gin.Context) string
	Visitors *sync.Map
}

// LimitadorPublico es para endpoints como login o registro.
// Identifica por IP y usa un límite estricto para evitar fuerza bruta.
// 1 petición por segundo, ráfaga máxima de 5.
func LimitadorPublico() gin.HandlerFunc {
	config := RateLimiterConfig{
		Rate:  rate.Limit(1),
		Burst: 5,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		Visitors: publicVisitors,
	}
	return crearLimitador(config)
}

// LimitadorProtegido es para los endpoints del uso regular de la app.
// Identifica por el usuario_id del JWT para mayor precisión.
// 10 peticiones por segundo, ráfaga máxima de 20.
func LimitadorProtegido() gin.HandlerFunc {
	config := RateLimiterConfig{
		Rate:  rate.Limit(10),
		Burst: 20,
		KeyFunc: func(c *gin.Context) string {
			// Intentamos obtener el ID inyectado por AuthRequerido
			usuarioID, existe := c.Get("usuario_id")
			if existe {
				if idStr, ok := usuarioID.(string); ok && idStr != "" {
					return idStr
				}
			}
			// Fallback a IP si algo raro pasa y no hay ID
			return c.ClientIP()
		},
		Visitors: protectedVisitors,
	}
	return crearLimitador(config)
}

// LimitadorWebSocket es específico para intentar conectarse al WebSocket.
// Muy estricto para evitar ataques de conexión al servidor WS.
// 1 conexión cada 5 segundos (0.2 req/s), ráfaga de 2.
func LimitadorWebSocket() gin.HandlerFunc {
	config := RateLimiterConfig{
		Rate:  rate.Limit(0.2), 
		Burst: 2,
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		Visitors: publicVisitors, // Reutilizamos el mapa de IPs públicas
	}
	return crearLimitador(config)
}

// crearLimitador es el motor principal que intercepta las peticiones
func crearLimitador(config RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := config.KeyFunc(c)

		// Buscar o crear el limiter para este usuario/IP de forma segura
		v, exists := config.Visitors.Load(key)
		var vis *visitor
		if !exists {
			vis = &visitor{
				limiter:  rate.NewLimiter(config.Rate, config.Burst),
				lastSeen: time.Now(),
			}
			config.Visitors.Store(key, vis)
		} else {
			vis = v.(*visitor)
			vis.lastSeen = time.Now() // Actualizamos la última interacción
		}

		// --- Headers informativos de Rate Limiting (Estándar de la industria) ---
		
		// X-RateLimit-Limit: Tokens máximos en ráfaga
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Burst))
		
		// X-RateLimit-Remaining: Tokens disponibles en este microsegundo
		remaining := vis.limiter.Tokens()
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(math.Round(remaining))))

		// ¿Pasa el peaje?
		if !vis.limiter.Allow() {
			// Calculamos cuánto tiempo tiene que esperar el usuario
			res := vis.limiter.Reserve()
			if !res.OK() {
				// Fallback si la reserva excede los límites (no debería pasar normalmente)
				c.Header("Retry-After", "5")
			} else {
				delay := res.Delay()
				res.Cancel() // Solo queremos leer el delay, cancelamos la reserva
				
				retryAfterSeconds := int(math.Ceil(delay.Seconds()))
				if retryAfterSeconds < 1 {
					retryAfterSeconds = 1
				}
				c.Header("Retry-After", strconv.Itoa(retryAfterSeconds))
			}

			// Rechazar con nuestro formato global
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   true,
				"codigo":  http.StatusTooManyRequests,
				"mensaje": "Has enviado demasiadas peticiones. Por favor, espera un momento.",
			})
			return
		}

		c.Next()
	}
}

// iniciarLimpieza recorre los visitantes y borra los inactivos para no agotar la RAM (Memory Leak)
func iniciarLimpieza(visitors *sync.Map, intervalo time.Duration, expiracion time.Duration) {
	for {
		time.Sleep(intervalo)
		now := time.Now()
		
		visitors.Range(func(key, value interface{}) bool {
			vis := value.(*visitor)
			if now.Sub(vis.lastSeen) > expiracion {
				visitors.Delete(key)
			}
			return true // Continuar iterando
		})
	}
}
