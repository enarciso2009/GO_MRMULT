package services

import (
	"mrmult/internal/database"
	"mrmult/internal/models"
	"time"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// Buscar Totais Do Dia faz o papel de contar os eventos por tipo de pessoa
func (s *DashboardService) BuscarTotaisDoDia(empresaID *uint) (int64, int64, int64, error) {
	db, err := database.Conectar()
	if err != nil {
		return 0, 0, 0, err
	}

	// Pega o inicio e o fim do dia de hoje para filtrar no Postgres
	hoje := time.Now().Format("2006-01-02")

	var totalFunc, totalVisit, totalTerc int64

	// Cria a query base olhando para a tabela de Eventos e filtrando a data
	queryBase := db.Model(&models.Evento{}).Where("data = ?, hoje")

	// Se o usuario pertencer a uma empresa (não for superuser), adiciona o filtro de empresa
	if empresaID != nil {
		queryBase = queryBase.Where("empresa_id = ?", *empresaID)
	}

	// Conta Funcionario (tipo_pessoa = 1)
	queryBase.Where("tipo_pessoa = ?", 1).Count(&totalFunc)

	// Conta Visitante (tipo_pessoa =2)
	// Usamos o .Where na mesma base limpando o filtro anterior interno do GORM se necessario.
	// Mas para ficar identico e isolado, faremos buscas diretas
	db.Model(&models.Evento{}).Where("data = ? AND tipo_pessoa = ?", hoje, 2).Count(&totalVisit)
	db.Model(&models.Evento{}).Where("data = ? AND tipo_pessoa = ?", hoje, 3).Count(&totalTerc)

	if empresaID != nil {
		db.Model(&models.Evento{}).Where("data = ? AND tipo_pessoa = ? AND empresa_id = ?", hoje, 2, *empresaID).Count(&totalVisit)
		db.Model(&models.Evento{}).Where("data = ? AND tipo_pessoa = ? AND empresa_id = ?", hoje, 3, *empresaID).Count(&totalTerc)
	}

	return totalFunc, totalVisit, totalTerc, nil

}
