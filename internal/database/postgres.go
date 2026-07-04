package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Conectar estabelece a conesão com o GORM e ativa o AutoMigrate
func Conectar() (*gorm.DB, error) {
	//Se ja estiver conectado, reaproveitar a conexão
	if DB != nil {
		return DB, nil
	}

	strConexão := "host=localhost port=5432 user=postgres password=Ev323232@ dbname=go_postgres sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(strConexão), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	return DB, nil

}
