package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type FuncionarioService struct{}

func NewFuncionarioService() *FuncionarioService {
	return &FuncionarioService{}
}

// Listar traz todos os Funcionarios da empresa e pre carrega seus funcionarios vinculados

func (s *FuncionarioService) Listar(EmpresaID *uint) ([]models.Funcionario, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var funcionarios []models.Funcionario

	query := db.Model(&models.Funcionario{}).Preload("GrupoRefeicao").Preload("Equipamentos")

	// Se nao for superuser (tiver ID da Empresa), aplica o filtro igual ao do Django
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", *EmpresaID)
	}

	err = query.Find(&funcionarios).Error
	return funcionarios, err
}

// Incluir salva um novo funcionario no banco de dados

func (s *FuncionarioService) Incluir(fu *models.Funcionario) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(fu).Error
}

// Alterar deleta o registro usando o ID
func (s *FuncionarioService) Alterar(fu *models.Funcionario, EmpresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	// Seguranca extra, garante que o funcionario pertence a empresa do usuario
	query := db.Model(&models.Funcionario{}).Where("id = ?", fu.ID)
	if EmpresaID != nil {
		query = query.Where("empresa_id = ?", *EmpresaID)
	}

	// Atualizar apenas os campos enviados
	return query.Updates(fu).Error

}

// Excluir deletar o registro usando o ID

func (s *FuncionarioService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	query := db.Where("id = ?", id)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	return query.Delete(&models.Funcionario{}).Error
}
