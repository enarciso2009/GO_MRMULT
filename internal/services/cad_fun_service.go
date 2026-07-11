package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type CadFunService struct{}

func NewCadFunService() *CadFunService {
	return &CadFunService{}
}

//Listar os Funcionarios filtrando pela empresa do usuario

func (f *CadFunService) Listar(EmpresaID *uint) ([]models.Funcionario, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var cadfun []models.Funcionario
	query := db.Model(&models.Funcionario{})

	// Se n'ao for superuser(tiver ID de Empresa), aplica o filtro igual ao do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", EmpresaID)
	}

	err = query.Find(&cadfun).Error
	return cadfun, err
}
