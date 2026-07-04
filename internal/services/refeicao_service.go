package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type RefeicaoService struct{}

func NewRefeicaoService() *RefeicaoService {
	return &RefeicaoService{}
}

func (s *RefeicaoService) Listar(empresaID *uint) ([]models.Refeicao, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}
	var refeicoes []models.Refeicao
	query := db.Model(&models.Refeicao{})
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	err = query.Find(&refeicoes).Error
	return refeicoes, err
}

func (s *RefeicaoService) Incluir(ref *models.Refeicao) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(ref).Error
}

func (s *RefeicaoService) Alterar(ref *models.Refeicao, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	query := db.Model(&models.Refeicao{}).Where("id_ref = ?", ref.IDRef)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}

	return query.Updates(ref).Error
}

func (s *RefeicaoService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	query := db.Where("id_ref = ?", id)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	return query.Delete(&models.Refeicao{}).Error
}
