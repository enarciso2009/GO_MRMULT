package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type VisitanteService struct{}

func NewVisitanteService() *VisitanteService {
	return &VisitanteService{}
}

// Listar traz todos os visitantes da empresa e pre carrega seus Visitantes vinculados

func (v *VisitanteService) Listar(EmpresaID *uint) ([]models.Visitante, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var visitantes []models.Visitante

	query := db.Model(&models.Visitante{}).Preload("GrupoRefeicao").Preload("Equipamentos").Preload("Funcionario")

	// Se nao for superuser (tiver ID da empresa), aplica o filtro igual o do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", *EmpresaID)

	}

	err = query.Find(&visitantes).Error
	return visitantes, err
}

// Incluir salva novo visitante no banco de dados

func (v *VisitanteService) Incluir(vi *models.Visitante) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(vi).Error
}

// Alterar deleta o registro usando o ID

func (v *VisitanteService) Alterar(vi *models.Visitante, EmpresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	// Seguranca extra, garante que o visitante pertence a empresa do usuario que esta cadastrando
	query := db.Model(&models.Visitante{}).Where("id = ?", vi.ID)
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", *EmpresaID)
	}

	// Atualiza apenas os campos enviados
	return query.Updates(vi).Error
}

// Excluir deletar o registro usando o ID

func (v *VisitanteService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	query := db.Where("id = ?", id)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	return query.Delete(&models.Visitante{}).Error
}
