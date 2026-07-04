package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type EquipamentoService struct{}

func NewEquipamentoService() *EquipamentoService {
	return &EquipamentoService{}
}

// Listar Busca os equipamentos filtrando pela empresa do usuario
func (s *EquipamentoService) Listar(EmpresaID *uint) ([]models.Equipamento, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var equipamentos []models.Equipamento
	query := db.Model(&models.Equipamento{})

	// Se não for superuser (tiver ID de Empresa), aplica o filtro igual ao do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", EmpresaID)
	}

	err = query.Find(&equipamentos).Error
	return equipamentos, err
}

// Incluir salva um novo equipamento no banco de dados
func (s *EquipamentoService) Incluir(eq *models.Equipamento) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(eq).Error
}

// Alterar atualizar os dados do equipamento baseado na Chave Primaria id_equip
func (s *EquipamentoService) Alterar(eq *models.Equipamento, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	// Segurança extra do Django: garante que o equipamento pertence a empresa do usuario
	query := db.Model(&models.Equipamento{}).Where("id_equip = ?", eq.IDEquip)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}

	// Atualizar apenas os campos enviados
	return query.Updates(eq).Error
}

// Excluir deleta o registro usando o ID
func (s *EquipamentoService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	query := db.Where("id_equip = ?", id)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	return query.Delete(&models.Equipamento{}).Error
}
