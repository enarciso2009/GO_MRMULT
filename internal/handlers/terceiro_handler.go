package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
	"strconv"
	"time"
)

var tmplTerceiro = template.Must(template.ParseFiles("web/templates/terceiro.html"))

// TelaTerceiroDados segue rigorosamente o padrao das outras strucks de tela
type TelaTerceiroDados struct {
	Terceiros         []models.Terceiro
	FuncionariosDispo []models.Funcionario
	GruposDispo       []models.GrupoRefeicao
	EquipamentosDispo []models.Equipamento
	IsErro            bool
	Mensagem          string
}

// TerceiroHandler Servico para carregar os funcionarios, grupos de refeicao e equipamentos disponiveis
type TerceiroHandler struct {
	serviceTerceiro    *services.TerceiroService
	serviceFuncionario *services.FuncionarioService
	serviceGrupo       *services.GrupoRefeicaoService
	serviceEquipamento *services.EquipamentoService
}

func NewTerceiroHandler(st *services.TerceiroService, sf *services.FuncionarioService, sg *services.GrupoRefeicaoService, se *services.EquipamentoService) *TerceiroHandler {
	return &TerceiroHandler{
		serviceTerceiro:    st,
		serviceFuncionario: sf,
		serviceGrupo:       sg,
		serviceEquipamento: se,
	}
}

func (h *TerceiroHandler) GerenciarTerceiro(w http.ResponseWriter, r *http.Request) {
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
		idTerStr := r.PostForm.Get("id_ter")

		matriculaInput := r.PostForm.Get("matricula")
		nomeInput := r.PostForm.Get("nome")
		empterInput := r.PostForm.Get("empter")
		documentoInput := r.PostForm.Get("documento")
		credencialInput := r.PostForm.Get("credencial")
		funcionarioidInput := r.PostForm.Get("func_id")
		gruporefeicaoidInput := r.PostForm.Get("grup_ref_id")
		datainicioInput := r.PostForm.Get("datainicio")
		horainicioInput := r.PostForm.Get("horainicio")
		datafimInput := r.PostForm.Get("datafim")
		horafimInput := r.PostForm.Get("horafim")

		// Conversao dos IDs numericos usando strconv
		grupoRefIDUint, _ := strconv.ParseUint(gruporefeicaoidInput, 10, 32)
		grupoRefID := uint(grupoRefIDUint)

		funcionarioIDUint, _ := strconv.ParseUint(funcionarioidInput, 10, 32)
		funcionarioID := uint(funcionarioIDUint)

		var idTer uint = 0
		if acao == "alterar" {
			idUint, _ := strconv.ParseUint(idTerStr, 10, 32)
			idTer = uint(idUint)
		}

		// Montagem do Objeto Base de Terceiros
		terceiro := models.Terceiro{
			ID:              idTer,
			Matricula:       &matriculaInput,
			Nome:            &nomeInput,
			EmpTer:          &empterInput,
			Documento:       &documentoInput,
			Credencial:      &credencialInput,
			FuncionarioID:   funcionarioID,
			GrupoRefeicaoID: grupoRefID,
			HoraInicio:      &horainicioInput,
			HoraFim:         &horafimInput,
		}

		if datainicioInput != "" {
			if d, err := time.Parse("2006-01-02", datainicioInput); err == nil {
				terceiro.DataInicio = &d
			}
		} // <-- CHAVE CORRIGIDA AQUI

		if datafimInput != "" {
			if d, err := time.Parse("2006-01-02", datafimInput); err == nil {
				terceiro.DataFim = &d
			}
		}

		// =============
		// ACAO: INCLUIR
		// =============

		if acao == "incluir" {
			equipamentosSelecionadosStr := r.PostForm["Equipamentos"]

			// VALIDACAO OBRIGATORIA: Se nao marcou nenhum equipamento, barra o cadastro
			if len(equipamentosSelecionadosStr) == 0 {
				h.renderizarComErro(w, &empresaPadraoID, "Erro: Voce deve selecionar ao menos um equipamento para o terceiro")
				return
			}

			novoTe := models.Terceiro{
				Matricula:       &matriculaInput,
				Nome:            &nomeInput,
				EmpTer:          &empterInput,
				Documento:       &documentoInput,
				Credencial:      &credencialInput,
				FuncionarioID:   funcionarioID,
				GrupoRefeicaoID: grupoRefID,
				HoraInicio:      &horainicioInput,
				HoraFim:         &horafimInput,
				EmpresaID:       &empresaPadraoID,
			}

			if datainicioInput != "" {
				if d, err := time.Parse("2006-01-02", datainicioInput); err == nil {
					novoTe.DataInicio = &d // <-- CORRIGIDO: Vinculando a novoTe e nao terceiro antigo
				}
			}

			if datafimInput != "" {
				if d, err := time.Parse("2006-01-02", datafimInput); err == nil {
					novoTe.DataFim = &d // <-- CORRIGIDO: Vinculando a novoTe e nao terceiro antigo
				}
			}

			// Varre as strings recebidas da tela e adiciona na lista virtual do visitante
			for _, idStr := range equipamentosSelecionadosStr {
				idUint, _ := strconv.ParseUint(idStr, 10, 32)
				novoTe.Equipamentos = append(novoTe.Equipamentos, models.Equipamento{
					IDEquip: uint(idUint),
				})
			}

			// O service executa o Db.Create salvando o terceiro e criando as amarracoes no banco de dados
			err = h.serviceTerceiro.Incluir(&novoTe)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao incluir terceiro: "+err.Error())
				return
			}
		}

		if acao == "alterar" {
			err = h.serviceTerceiro.Alterar(&terceiro, &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao alterar terceiro: "+err.Error()) // <-- CORRIGIDO .Error()
				return
			}
		}

		if acao == "excluir" {
			idUint, _ := strconv.ParseUint(idTerStr, 10, 32)
			err = h.serviceTerceiro.Excluir(uint(idUint), &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao excluir terceiro: "+err.Error())
				return
			}
		}

		http.Redirect(w, r, "/terceiro", http.StatusSeeOther)
		return
	}

	// ==================================
	// OPERACAO 2: EXIBICAO DA TELA (GET)
	// ==================================

	listaTerceiro, errT := h.serviceTerceiro.Listar(&empresaPadraoID)
	listaFuncionario, errF := h.serviceFuncionario.Listar(&empresaPadraoID)
	listaGrupos, errG := h.serviceGrupo.Listar(&empresaPadraoID)

	// CORRIGIDO: Busca equipamentos usando a variável padrão de ID ativa no escopo GET
	listaEquipamentos, errE := h.serviceEquipamento.Listar(&empresaPadraoID)

	data := TelaTerceiroDados{
		Terceiros:         listaTerceiro,
		FuncionariosDispo: listaFuncionario,
		GruposDispo:       listaGrupos,
		EquipamentosDispo: listaEquipamentos,
	}

	if errT != nil || errF != nil || errG != nil || errE != nil {
		data.IsErro = true
		data.Mensagem = "Erro ao carregar dados do banco de dados relacional."
	}

	tmplTerceiro.Execute(w, data)
}

func (h *TerceiroHandler) renderizarComErro(w http.ResponseWriter, empresaID *uint, msg string) {
	listaTerceiro, _ := h.serviceTerceiro.Listar(empresaID)
	listaFuncionario, _ := h.serviceFuncionario.Listar(empresaID)
	listaGrupos, _ := h.serviceGrupo.Listar(empresaID) //
	listaEquipamentos, _ := h.serviceEquipamento.Listar(empresaID)

	data := TelaTerceiroDados{
		Terceiros:         listaTerceiro,
		FuncionariosDispo: listaFuncionario,
		GruposDispo:       listaGrupos,
		EquipamentosDispo: listaEquipamentos,
		IsErro:            true,
		Mensagem:          msg,
	}

	tmplTerceiro.Execute(w, data)
}
