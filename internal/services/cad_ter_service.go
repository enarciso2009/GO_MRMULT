package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type CadTerService struct{}

func NewCadTerService() *CadTerService {
	return &CadTerService{}
}

// Listar os Visitantes Filtrando pela empresa do usuario

func (cv *CadTerService) Listar(EmpresaID *uint) ([]models.Terceiro, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var cadter []models.Terceiro
	query := db.Model(&models.Terceiro{}).Preload("Funcionario")

	// Se nao for superuser(tiver ID de Empresa), aplica o filtro igual ao do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", EmpresaID)
	}

	err = query.Find(&cadter).Error
	return cadter, err
}
