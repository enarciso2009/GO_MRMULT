package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
	"strconv"
	"time"
)

// Carrega o HTML apenas uma vez na inicializacao do sistema para alta performance
var tmplRefeicao = template.Must(template.ParseFiles("web/templates/refeicao.html"))

// Tabela Refeicao Data define o formato de dados que a tela Html espera receber
type TelaRefeicaoData struct {
	Refeicoes []models.Refeicao
	IsErro    bool
	Mensagem  string
}

type RefeicaoHandler struct {
	service *services.RefeicaoService
}

// NewRefeicaoHandler cria a instancia do controlador injetando o service de banco
func NewRefeicaoHandler(service *services.RefeicaoService) *RefeicaoHandler {
	return &RefeicaoHandler{service: service}
}

// GerenciarRefeicao centraliza as operacoes de exibicao e processamento da tela
func (h *RefeicaoHandler) GerenciarRefeicao(w http.ResponseWriter, r *http.Request) {
	// IMPORTANTE: Defina o id da empresa padrao para o ecossistema SaaS Multiempresa
	var empresaPadraoID uint = 1

	// =============================================
	// OPERACAO 1: PROCESSAMENTO DE ENVIO (POST)
	// =============================================
	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		idRefStr := r.FormValue("id_ref")
		nome := r.FormValue("nome")
		valorStr := r.FormValue("valor")
		dataInicioStr := r.FormValue("data_inicio")
		horaInicio := r.FormValue("hora_inicio")
		horaFim := r.FormValue("hora_fim")

		// Ajuste de horas para HH:MM:SS
		if len(horaInicio) == 5 {
			horaInicio = horaInicio + ":00"
		}
		if len(horaFim) == 5 {
			horaFim = horaFim + ":00"
		}

		// Converter o valor financeiro para float64
		valor, _ := strconv.ParseFloat(valorStr, 64)

		// Converter a string de data (AAAA-MM-DD) vinda do HTML para time.Time do Go
		var dataInicio time.Time
		if dataInicioStr != "" {
			dataInicio, _ = time.Parse("2006-01-02", dataInicioStr)
		} else {
			dataInicio = time.Now() // Fallback de seguranca
		}

		// ACAO: INCLUIR NOVA REFEICAO
		if acao == "incluir" {
			novaRef := models.Refeicao{
				Nome:       nome,
				HoraInicio: &horaInicio,
				HoraFim:    &horaFim,
				EmpresaID:  &empresaPadraoID,
			}

			err := h.service.Incluir(&novaRef, valor, dataInicio)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao incluir refeicao: "+err.Error())
				return
			}
		}

		// ACAO: ALTERAR REFEICAO EXISTENTE (HISTORICO)
		if acao == "alterar" {
			idRefUint, _ := strconv.ParseUint(idRefStr, 10, 32)
			idRef := uint(idRefUint)

			refAlterado := models.Refeicao{
				IDRef:      idRef,
				Nome:       nome,
				HoraInicio: &horaInicio,
				HoraFim:    &horaFim,
				EmpresaID:  &empresaPadraoID,
			}

			err := h.service.Alterar(&refAlterado, valor, dataInicio, &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao alterar refeicao: "+err.Error())
				return
			}
		}

		// ACAO: EXCLUIR REFEICAO
		if acao == "excluir" {
			idRefUint, _ := strconv.ParseUint(idRefStr, 10, 32)
			idRef := uint(idRefUint)

			err := h.service.Excluir(&idRef, &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao excluir refeicao: "+err.Error())
				return
			}
		}

		// Apos qualquer operacao com sucesso, redireciona via GET para limpar o formulario e evitar reenvio (PRG pattern)
		http.Redirect(w, r, "/refeicao", http.StatusSeeOther)
		return
	}

	// ===============================================
	// OPERACAO 2: RENDERIZACAO LIMPA DA TELA GET
	// ===============================================

	lista, err := h.service.Listar(&empresaPadraoID)
	data := TelaRefeicaoData{
		Refeicoes: lista,
	}

	if err != nil {
		data.IsErro = true
		data.Mensagem = "Erro ao listar refeicoes do banco: " + err.Error()
	}

	tmplRefeicao.Execute(w, data)

}

// Funcao auxiliar interna para renderizar a tela exibindo os alertas de erro na interface
func (h *RefeicaoHandler) renderizarComErro(w http.ResponseWriter, EmpresaID *uint, msg string) {
	lista, _ := h.service.Listar(EmpresaID)
	data := TelaRefeicaoData{
		Refeicoes: lista,
		IsErro:    true,
		Mensagem:  msg,
	}
	tmplRefeicao.Execute(w, data)
}
