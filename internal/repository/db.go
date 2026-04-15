package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB es la variable global que guarda la conexión activa
var DB *gorm.DB

// ConnectDB inicializa la conexión con la base de datos
func ConnectDB(dsn string) {
	var err error
	
	// Abrimos la conexión
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error crítico al conectar a la base de datos: %v", err)
	}

	fmt.Println("✅ Conexión a PostgreSQL (boleta_bdd) exitosa")

	// Aquí iría el AutoMigrate, pero lo omitimos porque 
	// ya creamos la estructura manual y perfectamente con SQL.
}