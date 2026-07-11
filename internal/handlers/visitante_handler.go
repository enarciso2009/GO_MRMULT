package handlers

import (
	"fmt"
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
	"strconv"
	"time"
)

var tmplVisitante = template.Must(template.ParseFiles("web/templates/visitante.html"))

// TelaVisitanteDados segue rigorosamente o padrao das outras strucks de tela

type TelaVisitanteDados struct {
	Visitantes        []models.Visitante
	FuncionariosDispo []models.Funcionario
	GruposDispo       []models.GrupoRefeicao
	EquipamentosDispo []models.Equipamento
	IsErro            bool
	Mensagem          string
}

// Servico para carregar os funcionarios, grupos de refeicao e equipamentos disponiveis
type VisitanteHandler struct {
	serviceVisitante   *services.VisitanteService
	serviceFuncionario *services.FuncionarioService
	serviceGrupo       *services.GrupoRefeicaoService
	serviceEquipamento *services.EquipamentoService
}

func NewVisitanteHandler(sv *services.VisitanteService, sf *services.FuncionarioService, sg *services.GrupoRefeicaoService, se *services.EquipamentoService) *VisitanteHandler {
	return &VisitanteHandler{
		serviceVisitante:   sv,
		serviceFuncionario: sf,
		serviceGrupo:       sg,
		serviceEquipamento: se,
	}
}

func (h *VisitanteHandler) GerenciarVisitante(w http.ResponseWriter, r *http.Request) {
	var empresaPadraoID uint = 1

	// ==============================================
	// OPERACAO 1: PROCESSAMENTO DE FORMULARIO (POST)
	// ==============================================

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.renderizarComErro(w, &empresaPadraoID, "Erro ao processar formulario: "+err.Error())
			return
		}

		acao := r.PostForm.Get("acao")
		idVisStr := r.PostForm.Get("id_vis")

		matriculaInput := r.PostForm.Get("matricula")
		nomeInput := r.PostForm.Get("nome")
		documentoInput := r.PostForm.Get("documento")
		credencialInput := r.PostForm.Get("credencial")
		funcionarioidInput := r.PostForm.Get("func_id")
		gruporefeicaoidInput := r.PostForm.Get("grup_ref_id")
		datainicioInput := r.PostForm.Get("datainicio")
		horainicioInput := r.PostForm.Get("horainicio")
		datafimInput := r.PostForm.Get("datafim")
		horafimInput := r.PostForm.Get("horafim")
		motivoInput := r.PostForm.Get("motivo")

		// Conversao dos IDs numericos usando strconv

		grupoRefIDUint, _ := strconv.ParseUint(gruporefeicaoidInput, 10, 32)
		grupoRefID := uint(grupoRefIDUint)

		funcionarioIDUint, _ := strconv.ParseUint(funcionarioidInput, 10, 32)
		funcionarioID := uint(funcionarioIDUint)

		var idVis uint = 0
		if acao == "alterar" {
			idUint, _ := strconv.ParseUint(idVisStr, 10, 32)
			idVis = uint(idUint)
		}

		// Montagem do Objeto Base de Visitante
		visitante := models.Visitante{
			ID:              idVis,
			Matricula:       &matriculaInput,
			Nome:            &nomeInput,
			Documento:       &documentoInput,
			Credencial:      &credencialInput,
			FuncionarioID:   funcionarioID,
			GrupoRefeicaoID: grupoRefID,
			HoraInicio:      &horainicioInput,
			HoraFim:         &horafimInput,
			Motivo:          &motivoInput,
			EmpresaID:       &empresaPadraoID,
		}
		if datainicioInput != "" {
			if t, err := time.Parse("2006-01-02", datainicioInput); err == nil {
				visitante.DataInicio = &t
			}
		}
		if datafimInput != "" {
			if t, err := time.Parse("2006-01-02", datafimInput); err == nil {
				visitante.DataFim = &t
			}
		}

		// =============
		// ACAO: INCLUIR
		// =============

		if acao == "incluir" {
			equipamentosSelecionadosStr := r.PostForm["equipamentos"]

			// VALIDACAO OBRIGATORIA: Se nao marcou nenhum equipamento, barra o cadastro

			if len(equipamentosSelecionadosStr) == 0 {
				h.renderizarComErro(w, &empresaPadraoID, "Erro: Voce deve selecionar ao menos um equipamento para o Visitante")
				return
			}

			novoVi := models.Visitante{
				Matricula:       &matriculaInput,
				Nome:            &nomeInput,
				Documento:       &documentoInput,
				Credencial:      &credencialInput,
				FuncionarioID:   funcionarioID,
				GrupoRefeicaoID: grupoRefID,
				HoraInicio:      &horainicioInput,
				HoraFim:         &horafimInput,
				EmpresaID:       &empresaPadraoID,
			}

			if datainicioInput != "" {
				if t, err := time.Parse("2006-01-02", datainicioInput); err == nil {
					visitante.DataInicio = &t
				}
			}

			if datafimInput != "" {
				if t, err := time.Parse("2006-01-02", datafimInput); err == nil {
					visitante.DataFim = &t
				}
			}

			// Varre as strings recebidas da tela e adiciona na lista virtual do Visitante
			for _, idStr := range equipamentosSelecionadosStr {
				idUint, _ := strconv.ParseUint(idStr, 10, 32)
				novoVi.Equipamentos = append(novoVi.Equipamentos, models.Equipamento{
					IDEquip: uint(idUint),
				})
			}

			// O service executa o DB.Create salvando o visitante e criando as amarracoes no banco de dados

			err = h.serviceVisitante.Incluir(&novoVi)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao incluir visitante: "+err.Error())
				return
			}

		}

		if acao == "alterar" {
			err = h.serviceVisitante.Alterar(&visitante, &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao alterar visitante: "+err.Error())
				return
			}

		}

		if acao == "excluir" {
			idUint, _ := strconv.ParseUint(idVisStr, 10, 32)
			err = h.serviceVisitante.Excluir(uint(idUint), &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao excluir visitante: "+err.Error())
				return
			}
		}

		http.Redirect(w, r, "/visitante", http.StatusSeeOther)
		return
	}

	// ==================================
	// OPERACAO 2: EXIBICAO DA TELA (GET)
	// ==================================

	listaVisitante, errV := h.serviceVisitante.Listar(&empresaPadraoID)
	listaFuncionario, errF := h.serviceFuncionario.Listar(&empresaPadraoID)
	listaGrupos, errG := h.serviceGrupo.Listar(&empresaPadraoID)

	// Busca equipamentos disponiveis para alimentar o select da esquerda
	listaEquipamentos, errE := h.serviceEquipamento.Listar(&empresaPadraoID)

	data := TelaVisitanteDados{
		Visitantes:        listaVisitante,
		FuncionariosDispo: listaFuncionario,
		GruposDispo:       listaGrupos,
		EquipamentosDispo: listaEquipamentos,
	}

	if errV != nil || errF != nil || errG != nil || errE != nil {
		data.IsErro = true
		// Isso vai te mostrar na tela qual das variáveis deu erro
		data.Mensagem = fmt.Sprintf("Erro DB: V:%v | F:%v | G:%v | E:%v", errV, errF, errG, errE)

		// E printa detalhado no terminal do GoLand
		fmt.Printf("❌ ERRO DETALHADO:\nVisitante: %v\nFuncionario: %v\nGrupo: %v\nEquipamento: %v\n", errV, errF, errG, errE)
		data.Mensagem = "Erro ao carregar dados do banco de dados relacional."

	}

	tmplVisitante.Execute(w, data)

}

func (h *VisitanteHandler) renderizarComErro(w http.ResponseWriter, empresaID *uint, msg string) {
	listaVisitante, _ := h.serviceVisitante.Listar(empresaID)
	listaFuncionario, _ := h.serviceFuncionario.Listar(empresaID)
	listaGrupos, _ := h.serviceGrupo.Listar(empresaID)
	listaEquipamentos, _ := h.serviceEquipamento.Listar(empresaID)

	data := TelaVisitanteDados{
		Visitantes:        listaVisitante,
		FuncionariosDispo: listaFuncionario,
		GruposDispo:       listaGrupos,
		EquipamentosDispo: listaEquipamentos,
		IsErro:            true,
		Mensagem:          msg,
	}
	tmplVisitante.Execute(w, data)
}
