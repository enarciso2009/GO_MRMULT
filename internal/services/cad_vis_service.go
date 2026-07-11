package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type CadVisService struct{}

func NewCadVisService() *CadVisService {
	return &CadVisService{}
}

// Listar os Visitantes Filtrando pela empresa do usuario

func (cv *CadVisService) Listar(EmpresaID *uint) ([]models.Visitante, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var cadvis []models.Visitante
	query := db.Model(&models.Visitante{}).Preload("Funcionario")

	// Se nao for superuser(tiver ID de Empresa), aplica o filtro igual ao do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", EmpresaID)
	}

	err = query.Find(&cadvis).Error
	return cadvis, err
}
