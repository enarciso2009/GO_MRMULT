package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type EmpresaService struct{}

func NewEmpresaService() *EmpresaService {
	return &EmpresaService{}
}

// Listar traz os usuarios.

func (e *EmpresaService) Listar(empresaID *uint) ([]models.Empresa, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var empresas []models.Empresa
	query := db.Model(&models.Empresa{})

	if empresaID != nil {
		query = query.Where("ID = ?", *empresaID)
	}

	err = query.Find(&empresas).Error
	return empresas, err
}

// Incluir criar Empresas

func (e *EmpresaService) Incluir(em *models.Empresa) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(em).Error
}

// Alterar atualizar os dados da Empresa baseado na chave primaria ID

func (e *EmpresaService) Alterar(em *models.Empresa, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	// Segurança extra do Django: garante que o usuario que esta alterando pertence a empresa que esta sendo alterada
	query := db.Model(&models.Empresa{}).Where("ID = ?", empresaID)
	if empresaID != nil {
		query = query.Where("ID = ?", *empresaID)
	}

	// Atualizar apenas os campos enviados

	return query.Updates(em).Error
}

// Excluir deleta o registro usando o ID

func (e *EmpresaService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	query := db.Where("ID = ?", id)
	if empresaID != nil {
		query = query.Where("ID = ?", *empresaID)
	}
	return query.Delete(&models.Empresa{}).Error
}
