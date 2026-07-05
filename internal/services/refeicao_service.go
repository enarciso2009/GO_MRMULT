// internal/services/refeicao_service.go
package services

import (
	"errors"
	"mrmult/internal/database"
	"mrmult/internal/models"
	"time"

	"gorm.io/gorm"
)

type RefeicaoService struct{}

func NewRefeicaoService() *RefeicaoService {
	return &RefeicaoService{}
}

// Listar traz as refeicoes e ja carrega o preco que esta ativo no momento (data_fim IS NULL)
func (s *RefeicaoService) Listar(empresaID *uint) ([]models.Refeicao, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}
	var refeicoes []models.Refeicao

	// O Preload carrega apenas o historico de preco que esta vigente (sem data_fim)
	query := db.Preload("HistoricoPrecos", "data_fim IS NULL")
	if empresaID != nil {
		query = query.Where("empresa_id = ?", *empresaID)
	}
	err = query.Find(&refeicoes).Error
	return refeicoes, err
}

// Incluir cria a refeicao e ja abre o primeiro historico de preco dela
func (s *RefeicaoService) Incluir(ref *models.Refeicao, valorInicial float64, dataInicio time.Time) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	// Usamos uma transacao para garantir que insira nas duas tabelas ou em nenhuma
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Salva a refeicao na tabela pai (id_ref sera gerado aqui)
		if err := tx.Create(ref).Error; err != nil {
			return err
		}

		// 2. Cria o primeiro preco na tabela filha apontando para a refeicao criada
		preco := models.PrecoRefeicao{
			RefeicaoID: ref.IDRef,
			Valor:      valorInicial,
			DataInicio: dataInicio,
			DataFim:    nil, // Comeca sem data fim (vigente)
		}

		return tx.Create(&preco).Error
	})
}

// Alterar avalia se o valor mudou. Se mudou, encerra o preco antigo e abre um historico novo
func (s *RefeicaoService) Alterar(ref *models.Refeicao, novoValor float64, dataNovaVigencia time.Time, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Atualiza os dados da refeicao pai (nome, horarios) - Corrigido para usar tx em vez de db
		queryPai := tx.Model(&models.Refeicao{}).Where("id_ref = ?", ref.IDRef)
		if empresaID != nil {
			queryPai = queryPai.Where("empresa_id = ?", *empresaID)
		}
		if err := queryPai.Updates(ref).Error; err != nil {
			return err
		}

		// 2. Busca o preco atual vigente para verificar se houve alteracao no valor financeiro
		var precoAtual models.PrecoRefeicao
		err := tx.Where("id_ref = ? AND data_fim IS NULL", ref.IDRef).First(&precoAtual).Error

		if err == nil {
			// Se o valor digitado for diferente do atual, aplicamos a regra historico
			if precoAtual.Valor != novoValor {

				// CORREÇÃO VISUAL AQUI: Calcula a data primeiro em uma variável isolada
				dataFimCalculada := dataNovaVigencia.AddDate(0, 0, -1)

				// Agora sim aplicamos o ponteiro (&) sobre a variável criada na memória
				precoAtual.DataFim = &dataFimCalculada

				if err := tx.Save(&precoAtual).Error; err != nil {
					return err
				}

				// Cria a nova linha do preco com o novo valor historico
				novoPreco := models.PrecoRefeicao{
					RefeicaoID: ref.IDRef,
					Valor:      novoValor,
					DataInicio: dataNovaVigencia,
					DataFim:    nil,
				}
				if err := tx.Create(&novoPreco).Error; err != nil {
					return err
				}
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Caso a refeicao por algum motivo nao tivesse preco ativo, cria o primeiro
			novoPreco := models.PrecoRefeicao{
				RefeicaoID: ref.IDRef,
				Valor:      novoValor,
				DataInicio: dataNovaVigencia,
				DataFim:    nil,
			}
			return tx.Create(&novoPreco).Error
		} else {
			return err
		}

		return nil
	})
}

// Excluir apaga a refeicao (o banco ou as configuracoes do GORM cuidam do cascade se configurado)
func (s *RefeicaoService) Excluir(idRef *uint, empresaID *uint) error {
	db, err := database.Conectar()
	if err != nil {
		return err
	}
	// Se ponteiro for null por algum motivo, evitamos um panico no sistema retornando um erro limpo
	if idRef == nil {
		return errors.New("O ID da refeicao nao pode ser nulo")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Apaga os historicos de preco associados (usamos *idRef para pegar o valor dentro do ponteiro)
		if err := tx.Where("id_ref = ?", *idRef).Delete(&models.PrecoRefeicao{}).Error; err != nil {
			return err
		}

		query := tx.Where("id_ref = ?", *idRef)
		if empresaID != nil {
			query = query.Where("empresa_id = ?", *empresaID)
		}
		return query.Delete(&models.Refeicao{}).Error
	})
}
