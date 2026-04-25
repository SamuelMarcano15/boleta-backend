package security

import (
	"os" 
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// obtenerSecreto extrae la llave en tiempo real del entorno
func obtenerSecreto() []byte {
	secreto := os.Getenv("JWT_SECRET")
	if secreto == "" {
		// Si olvidaste poner el .env, la app hace "panic" y no arranca por seguridad
		panic("CRÍTICO: JWT_SECRET no está configurado en las variables de entorno") 
	}
	return []byte(secreto)
}

// HashearClave convierte la contraseña en un texto ilegible
func HashearClave(clave string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(clave), 14)
	return string(bytes), err
}

// VerificarClave compara la contraseña plana con el Hash de la BD
func VerificarClave(clavePlana, claveHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(claveHash), []byte(clavePlana))
	return err == nil
}

// GenerarJWT crea el "Pase VIP" para el usuario
func GenerarJWT(usuarioID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": usuarioID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Aquí usamos la llave de grado militar que viene del .env
	return token.SignedString(obtenerSecreto())
}

func ValidarJWT(tokenString string) (string, error) {
	// Parseamos el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar que se haya firmado con el algoritmo correcto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// Le entregamos nuestra llave secreta del .env para que verifique la firma
		return obtenerSecreto(), nil
	})

	if err != nil {
		return "", err
	}

	// Si es válido, extraemos el ID del usuario (el "sub" que guardamos antes)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID := claims["sub"].(string)
		return usuarioID, nil
	}

	return "", jwt.ErrSignatureInvalid
}