package models

import (
	"time"
)

// 1. EMPRESA

type Empresa struct {
	ID   uint   `gorm:"primaryKey"`
	Nome string `gorm:"type:varchar(100); not null"`
	CNPJ string `gorm:"type:varchar(18);unique;not null"`
}

// 2. GRUPO REFEICAOz

type GrupoRefeicao struct {
	IDGrup    uint       `gorm:"primaryKey;column:id_grup"`
	Nome      *string    `gorm:"type:varchar(100)"`
	EmpresaID *uint      `gorm:"column:empresa_id"`
	Empresa   *Empresa   `gorm:"foreignKey:EmpresaID"`
	Refeicoes []Refeicao `gorm:"many2many:inter_grup_refs;joinForeignKey:id_grup;joinReferences:id_ref"`
}

// 3. REFEICAO

type Refeicao struct {
	IDRef      uint     `gorm:"primaryKey;column:id_ref"`
	Nome       string   `gorm:"type:varchar(100); not null"`
	HoraInicio *string  `gorm:"type:time without time zone;column:hora_inicio"`
	HoraFim    *string  `gorm:"type:time without time zone;column:hora_fim"`
	EmpresaID  *uint    `gorm:"column:empresa_id"`
	Empresa    *Empresa `gorm:"foreignKey:EmpresaID"`
	// Relacionamento com o historico de valores comerciais
	HistoricoPrecos []PrecoRefeicao `gorm:"foreignKey:RefeicaoID"`
}

// ObterPrecoAtual varre a lista pre-carregada e traz o valor que esta sem data_fim (ativo)
func (r Refeicao) ObterPrecoAtual() float64 {
	for _, p := range r.HistoricoPrecos {
		if p.DataFim == nil {
			return p.Valor
		}
	}
	return 0.00 // Caso nao encontre nenhumn preco ativo
}

// ObterInicioVigencia traz a data formatada em padrao brasileiro para a tabela
func (r Refeicao) ObterInicioVigencia() string {
	for _, p := range r.HistoricoPrecos {
		if p.DataFim == nil {
			return p.DataInicio.Format("02/01/2006")
		}
	}
	return "--/--/----"
}

// 4. PRECO REFEICAO historico de precos

type PrecoRefeicao struct {
	IDPreco    uint       `gorm:"primaryKey;column:id_preco"`
	RefeicaoID uint       `gorm:"column:id_ref;not null"`
	Valor      float64    `gorm:"type:numeric(5,2); not null"`
	DataInicio time.Time  `gorm:"type:date;not null;column:data_inicio"`
	DataFim    *time.Time `gorm:"type:date;column:data_fim"` // NULL = Valor ativo atual
}

// 5. INTER_GRUP_REF tabela intermediaria do ManyToMany

type InterGrupRef struct {
	IDInter         uint          `gorm:"primaryKey;column:id_inter"`
	GrupoRefeicaoID uint          `gorm:"column:id_grup;uniqueIndex:idx_grup_ref"`
	GrupoRefeicao   GrupoRefeicao `gorm:"foreignKey:GrupoRefeicaoID"`
	RefeicaoID      uint          `gorm:"column:id_ref;uniqueIndex:idx_grup_ref"`
	Refeicao        Refeicao      `gorm:"foreignKey:RefeicaoID"`
	EmpresaID       *uint         `gorm:"column:empresa_id"`
	Empresa         *Empresa      `gorm:"foreignKey:EmpresaID"`
}

// 6. FUNCIONARIOS

type Funcionario struct {
	ID              uint          `gorm:"primaryKey"`
	Matricula       string        `gorm:"type:varchar(15);not null"`
	Nome            *string       `gorm:"type:varchar(100)"`
	Admissao        *time.Time    `gorm:"type:date"`
	Departamento    *string       `gorm:"type:varchar(100)"`
	CentroDeCusto   *string       `gorm:"type:varchar(100);column:centro_de_custo"`
	Cargo           *string       `gorm:"type:varchar(50)"`
	Documento       *string       `gorm:"type:varchar(50)"`
	Credencial      *string       `gorm:"type:varchar(50)"`
	GrupoRefeicaoID uint          `gorm:"column:grup_ref_id;not null"`
	GrupoRefeicao   GrupoRefeicao `gorm:"foreignKey:GrupoRefeicaoID"`
	TipoPessoa      int           `gorm:"default:1;not null"`
	Ativo           *string       `gorm:"type:varchar(50)"`
	EmpresaID       *uint         `gorm:"column:empresa_id"`
	Empresa         *Empresa      `gorm:"foreignKey:EmpresaID"`
	Equipamentos    []Equipamento `gorm:"many2many:Funcionario_equipamento;joinForeignKey:funcionarios_id;joinReferences:id_equip"`
}

func (f Funcionario) FormatarAdmissao() string {
	if f.Admissao == nil {
		return "--"
	}
	return f.Admissao.Format("02/01/2006")
}

// 7. VISITANTE
type Visitante struct {
	ID              uint          `gorm:"primaryKey"`
	Matricula       *string       `gorm:"type:varchar(100)"`
	Nome            *string       `gorm:"type:varchar(100)"`
	Documento       *string       `gorm:"type:varchar(50)"`
	Credencial      *string       `gorm:"type:varchar(50)"`
	FuncionarioID   uint          `gorm:"column:func_id;not null"`
	Funcionario     Funcionario   `gorm:"foreignKey:FuncionarioID"`
	GrupoRefeicaoID uint          `gorm:"column:grup_ref_id;not null"`
	GrupoRefeicao   GrupoRefeicao `gorm:"foreignKey:GrupoRefeicaoID"`
	DataInicio      *time.Time    `gorm:"type:date"`
	HoraInicio      *string       `gorm:"type:time without time zone;column:hora_inicio"`
	DataFim         *time.Time    `gorm:"type:date"`
	HoraFim         *string       `gorm:"type:time without time zone;column:hora_fim"`
	Motivo          *string       `gorm:"type:varchar(50)"`
	TipoPessoa      int           `gorm:"default:2;not null"`
	EmpresaID       *uint         `gorm:"column:empresa_id"`
	Empresa         *Empresa      `gorm:"foreignKey:EmpresaID"`
	Equipamentos    []Equipamento `gorm:"many2many:Visitante_equipamento;joinForeignKey:visitantes_id;joinReferences:id_equip"`
}

// FormatarAdmissao verifica se o ponteiro não é nulo e formata a data para DD/MM/AAAA
func (v Visitante) FormatarDataInicio() string {
	if v.DataInicio == nil {
		return "--"
	}
	// O Go usa uma data de referência específica para formatação: 02/01/2006
	return v.DataInicio.Format("02/01/2006")
}

func (v Visitante) FormatarDataFim() string {
	if v.DataFim == nil {
		return "--"
	}
	return v.DataFim.Format("02/01/2006")
}

// 8. TERCEIRO
type Terceiro struct {
	ID              uint          `gorm:"primaryKey"`
	Matricula       *string       `gorm:"type:varchar(100)"`
	Nome            *string       `gorm:"type:varchar(100)"`
	EmpTer          *string       `gorm:"type:varchar(100);column:emp_ter"`
	Documento       *string       `gorm:"type:varchar(50)"`
	Credencial      *string       `gorm:"type:varchar(50)"`
	FuncionarioID   uint          `gorm:"column:func_id;not null"`
	Funcionario     Funcionario   `gorm:"foreignKey:FuncionarioID"`
	GrupoRefeicaoID uint          `gorm:"column:grup_ref_id;not null"`
	GrupoRefeicao   GrupoRefeicao `gorm:"foreignKey:GrupoRefeicaoID"`
	DataInicio      *time.Time    `gorm:"type:date"`
	HoraInicio      *string       `gorm:"type:time without time zone;column:hora_inicio"`
	DataFim         *time.Time    `gorm:"type:date"`
	HoraFim         *string       `gorm:"type:time without time zone;column:hora_fim"`
	TipoPessoa      int           `gorm:"default:3;not null"`
	EmpresaID       *uint         `gorm:"column:empresa_id"`
	Empresa         *Empresa      `gorm:"foreignKey:EmpresaID"`
	Equipamentos    []Equipamento `gorm:"many2many:Terceiro_equipamento;joinForeignKey:terceiros_id;joinReferences:id_equip"`
}

// 9. EQUIPAMENTO
type Equipamento struct {
	IDEquip   uint     `gorm:"primaryKey;column:id_equip"`
	Nome      *string  `gorm:"type:varchar(100)"`
	IP        *string  `gorm:"type:varchar(15)"`
	Mask      *string  `gorm:"type:varchar(100)"`
	EmpresaID *uint    `gorm:"column:empresa_id"`
	Empresa   *Empresa `gorm:"foreignKey:EmpresaID"`
}

// 10. EVENTO
type Evento struct {
	ID         uint       `gorm:"primaryKey"` // O GORM precisa de uma PK numérica interna idealmente
	IDEvento   string     `gorm:"type:varchar(15);column:id_evento;not null"`
	Matricula  *string    `gorm:"type:varchar(15)"`
	Nome       *string    `gorm:"type:varchar(100)"`
	Data       *time.Time `gorm:"type:date"`
	Hora       *string    `gorm:"type:time"`
	EquipID    *string    `gorm:"type:varchar(10);column:equip_id"`
	EquipNome  *string    `gorm:"type:varchar(100);column:equip_nome"`
	TipoPessoa *int       `gorm:"column:tipo_pessoa"`
	EmpresaID  *uint      `gorm:"column:empresa_id"`
	Empresa    *Empresa   `gorm:"foreignKey:EmpresaID"`
}

// 11. USUÁRIO
type Usuario struct {
	IDUser    string   `gorm:"primaryKey;type:varchar(10);column:id_user"`
	Nome      string   `gorm:"type:varchar(100);not null"`
	Email     string   `gorm:"type:varchar(100);not null"`
	Usuario   string   `gorm:"type:varchar(100);not null"`
	Senha     string   `gorm:"type:varchar(100);not null"`
	Permissao *string  `gorm:"type:varchar(50)"`
	EmpresaID *uint    `gorm:"column:empresa_id"`
	Empresa   *Empresa `gorm:"foreignKey:EmpresaID"`
}

// 12. PARAMETRO
type Parametro struct {
	ID             uint     `gorm:"primaryKey"`
	IDParam        string   `gorm:"type:varchar(15);column:id_param;not null"`
	Nome           *string  `gorm:"type:varchar(100)"`
	ModPadraoUsu   bool     `gorm:"default:false;column:mod_padrao_usu"`
	ModCreditoUsu  bool     `gorm:"default:false;column:mod_credito_usu"`
	ModPadraoVisi  bool     `gorm:"default:false;column:mod_padrao_visi"`
	ModCreditoVisi bool     `gorm:"default:false;column:mod_credito_visi"`
	EmpresaID      *uint    `gorm:"column:empresa_id"`
	Empresa        *Empresa `gorm:"foreignKey:EmpresaID"`
}
