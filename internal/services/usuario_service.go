package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
)

type UsuarioService struct{}

func NewUsuarioService() *UsuarioService {
	return &UsuarioService{}
}

// Listar traz os usuarios.

func (u *UsuarioService) Listar(empresaID *uint) ([]models.Usuario, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}
	var usuarios []models.Usuario
	query := db.Model(&models.Usuario{})

	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}

	err = query.Find(&usuarios).Error
	return usuarios, err

}

// Incluir cria o usuario

func (u *UsuarioService) Incluir(us *models.Usuario) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	return db.Create(us).Error
}

// Alterar atualizar os dados do Usuario baseado na chave primaria id_user

func (u *UsuarioService) Alterar(us *models.Usuario, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	// Seeguranca extra do Django: garante que o usuario pertence a empresa do usuario
	query := db.Model(&models.Usuario{}).Where("id_user = ?", us.IDUser)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}

	// Atualizar apenas os campos enviados

	return query.Updates(us).Error
}

// Excluir deleta o registro usando o ID
func (u *UsuarioService) Excluir(id uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	query := db.Where("id_user = ?", id)
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	return query.Delete(&models.Usuario{}).Error
}
