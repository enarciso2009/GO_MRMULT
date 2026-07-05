package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"

	"gorm.io/gorm"
)

type GrupoRefeicaoService struct{}

func NewGrupoRefeicaoService() *GrupoRefeicaoService {
	return &GrupoRefeicaoService{}
}

// Listar traz todos os grupos da empesa e pre carrega suas refeicoes vinculadas

func (s *GrupoRefeicaoService) Listar(empresaID *uint) ([]models.GrupoRefeicao, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}
	var grupos []models.GrupoRefeicao

	query := db.Preload("Refeicoes") // Puxa as refeicoes associadas da tabela many-to-many
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}

	err = query.Find(&grupos).Error
	return grupos, err
}

// Salvar trata a inclusao e alteracao. O segredo do many-to-many esta no tx.Association
func (s *GrupoRefeicaoService) Salvar(grupo *models.GrupoRefeicao, idsRefeicoes []uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Busca as structs completas das refeicoes selecionadas para vincular
		var refeicoesSelecionadas []models.Refeicao
		if len(idsRefeicoes) > 0 {
			if err := tx.Where("id_ref IN ?", idsRefeicoes).Find(&refeicoesSelecionadas).Error; err != nil {
				return err
			}
		}
		// 2. se for um grupo novo (id zero), cria o registro pai primeiro
		if grupo.IDGrup == 0 {
			if err := tx.Create(grupo).Error; err != nil {
				return err
			}
		} else {
			// Se ja existe, atualiza os dados basicos (Nome)
			if err := tx.Model(&models.GrupoRefeicao{}).Where("id_grup = ?", grupo.IDGrup).Updates(grupo).Error; err != nil {
				return err
			}
		}
		// 3. MAGICA DO GORM: Substitui todas as amarracoes antigas na tabela intermediaria pelas nova
		// Isso limpla o que foi desmarcado e insere o que foi marcado automaticamente
		err := tx.Model(grupo).Association("Refeicoes").Replace(&refeicoesSelecionadas)
		if err != nil {
			return err
		}
		// 4. BONUS DE GOVERNANCA: Garante que linhas criadas na 'inter_grup_refs' herdem o empresa_id
		if grupo.EmpresaID != nil && len(idsRefeicoes) > 0 {
			tx.Model(&models.InterGrupRef{}).Where("id_grup = ?", grupo.IDGrup).Update("empresa_id", *grupo.EmpresaID)
		}

		return nil
	})

}

// Excluir remove o grupo (o GORM limpa a tabela intermediaria automaticamente se configurado)
func (s *GrupoRefeicaoService) Excluir(idGrup uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var grupo models.GrupoRefeicao
		if err := tx.Where("id_grup = ?", idGrup).First(&grupo).Error; err != nil {
			return err
		}
		// Limpa os vinculos many-to-many na tabela intermediaria primeiro
		tx.Model(&grupo).Association("Refeicoes").Clear()

		query := tx.Where("id_grup = ?", idGrup)
		if empresaID != nil {
			query = query.Where("empresa_id = ?", *empresaID)
		}
		return query.Delete(&models.GrupoRefeicao{}).Error
	})
}
