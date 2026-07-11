package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type TerceiroService struct{}

func NewTerceiroService() *TerceiroService {
	return &TerceiroService{}
}

func (t *TerceiroService) Listar(EmpresaID *uint) ([]models.Terceiro, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var terceiros []models.Terceiro

	query := db.Model(&models.Terceiro{}).Preload("GrupoRefeicao").Preload("Equipamentos").Preload("Funcionario")

	// Se nao for superuser (tiver ID da empresa), aplica o filtro igual o do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", *EmpresaID)
	}

	err = query.Find(&terceiros).Error
	return terceiros, err
}

// Incluir salva novo terceiro no banco de dados
func (t *TerceiroService) Incluir(te *models.Terceiro) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(te).Error
}

// Alterar deletar o registro usando o ID

func (t *TerceiroService) Alterar(te *models.Terceiro, EmpresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	// Seguranca extra, garante que o terceiro pertence a empresa do usuario que esta cadastrando
	query := db.Model(&models.Terceiro{}).Where("id = ?", te.ID)
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", *EmpresaID)
	}

	// Atualiza apenas os campos enviados
	return query.Updates(te).Error
}

// Excluir deletar o registro usando o ID

func (t *TerceiroService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	query := db.Where("id = ?", id)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	return query.Delete(&models.Terceiro{}).Error
}
