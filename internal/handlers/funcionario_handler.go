package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
	"strconv"
	"time"
)

var tmplFuncionario = template.Must(template.ParseFiles("web/templates/funcionario.html"))

// TelaFuncionarioData segue rigorosamente o padrão das outras structs de tela
type TelaFuncionarioData struct {
	Funcionarios      []models.Funcionario
	GruposDispo       []models.GrupoRefeicao
	EquipamentosDispo []models.Equipamento
	IsErro            bool
	Mensagem          string
}

type FuncionarioHandler struct {
	serviceFuncionario *services.FuncionarioService
	serviceGrupo       *services.GrupoRefeicaoService
	serviceEquipamento *services.EquipamentoService // Serviço para carregar os equipamentos disponíveis
}

func NewFuncionarioHandler(sf *services.FuncionarioService, sg *services.GrupoRefeicaoService, se *services.EquipamentoService) *FuncionarioHandler {
	return &FuncionarioHandler{
		serviceFuncionario: sf,
		serviceGrupo:       sg,
		serviceEquipamento: se,
	}
}

func (h *FuncionarioHandler) GerenciarFuncionario(w http.ResponseWriter, r *http.Request) {
	var empresaPadraoID uint = 1

	// ==========================================
	// OPERAÇÃO 1: PROCESSAMENTO DE FORMULÁRIO (POST)
	// ==========================================
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.renderizarComErro(w, &empresaPadraoID, "Erro ao processar formulário: "+err.Error())
			return
		}

		acao := r.PostForm.Get("acao")
		idFunStr := r.PostForm.Get("id_fun")

		matriculaInput := r.PostForm.Get("matricula")
		nomeInput := r.PostForm.Get("nome")
		admissaoInput := r.PostForm.Get("admissao")
		departamentoInput := r.PostForm.Get("departamento")
		centrodecustoInput := r.PostForm.Get("centrodecusto")
		cargoInput := r.PostForm.Get("cargo")
		documentoInput := r.PostForm.Get("documento")
		credencialInput := r.PostForm.Get("credencial")
		gruporefeicaoidInput := r.PostForm.Get("grup_ref_id")
		ativoInput := r.PostForm.Get("ativo")

		// Conversão dos IDs Numéricos usando strconv
		grupoRefIDUint, _ := strconv.ParseUint(gruporefeicaoidInput, 10, 32)
		grupoRefID := uint(grupoRefIDUint)

		var idFun uint = 0
		if acao == "alterar" {
			idUint, _ := strconv.ParseUint(idFunStr, 10, 32)
			idFun = uint(idUint)
		}

		// Montagem do Objeto Base de Funcionário
		funcionario := models.Funcionario{
			ID:              idFun,
			Matricula:       matriculaInput,
			Nome:            &nomeInput,
			Departamento:    &departamentoInput,
			CentroDeCusto:   &centrodecustoInput,
			Cargo:           &cargoInput,
			Documento:       &documentoInput,
			Credencial:      &credencialInput,
			GrupoRefeicaoID: grupoRefID,
			Ativo:           &ativoInput,
			EmpresaID:       &empresaPadraoID,
		}

		// Tratamento de conversão da data de Admissão
		if admissaoInput != "" {
			if t, err := time.Parse("2006-01-02", admissaoInput); err == nil {
				funcionario.Admissao = &t
			}
		}

		// Executa as operações através da camada de serviços
		if acao == "incluir" {
			equipamentosSelecionadosStr := r.PostForm["equipamentos"]

			// VALIDACAO OBRIGATORIA: Se nao marcou nenhum equipamento, barra o cadastro
			if len(equipamentosSelecionadosStr) == 0 {
				h.renderizarComErro(w, &empresaPadraoID, "Erro: Voce deve selecionar ao menos um equipamento para o funcionario")
				return
			}

			novoFu := models.Funcionario{
				Matricula:       matriculaInput,
				Nome:            &nomeInput,
				Departamento:    &departamentoInput,
				CentroDeCusto:   &centrodecustoInput,
				Cargo:           &cargoInput,
				Documento:       &documentoInput,
				Credencial:      &credencialInput,
				GrupoRefeicaoID: grupoRefID,
				Ativo:           &ativoInput,
				EmpresaID:       &empresaPadraoID,
			}

			if admissaoInput != "" {
				if t, err := time.Parse("2006-01-02", admissaoInput); err == nil {
					novoFu.Admissao = &t
				}
			}

			// Varre as strings recebidas da tela e adiciona na lista virtual do Funcionario
			for _, idStr := range equipamentosSelecionadosStr {
				idUint, _ := strconv.ParseUint(idStr, 10, 32)
				novoFu.Equipamentos = append(novoFu.Equipamentos, models.Equipamento{
					IDEquip: uint(idUint),
				})
			}

			// o service executa o DB.Create salvando o funcionario e criando as amarracoes no banco

			err = h.serviceFuncionario.Incluir(&novoFu)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao incluir funcionário: "+err.Error())
				return
			}
		}

		if acao == "alterar" {
			err = h.serviceFuncionario.Alterar(&funcionario, &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao alterar funcionário: "+err.Error())
				return
			}
		}

		if acao == "excluir" {
			idUint, _ := strconv.ParseUint(idFunStr, 10, 32)
			err = h.serviceFuncionario.Excluir(uint(idUint), &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao excluir funcionário: "+err.Error())
				return
			}
		}

		http.Redirect(w, r, "/funcionario", http.StatusSeeOther)
		return
	}

	// ==========================================
	// OPERAÇÃO 2: EXIBIÇÃO DA TELA (GET)
	// ==========================================
	listaFuncionarios, errF := h.serviceFuncionario.Listar(&empresaPadraoID)
	listaGrupos, errG := h.serviceGrupo.Listar(&empresaPadraoID)

	// Busca equipamentos disponíveis para alimentar o select da esquerda
	listaEquipamentos, errE := h.serviceEquipamento.Listar(&empresaPadraoID)

	data := TelaFuncionarioData{
		Funcionarios:      listaFuncionarios,
		GruposDispo:       listaGrupos,
		EquipamentosDispo: listaEquipamentos,
	}

	if errF != nil || errG != nil || errE != nil {
		data.IsErro = true
		data.Mensagem = "Erro ao carregar dados do banco de dados relacional."
	}

	tmplFuncionario.Execute(w, data)
}

func (h *FuncionarioHandler) renderizarComErro(w http.ResponseWriter, empresaID *uint, msg string) {
	listaFuncionarios, _ := h.serviceFuncionario.Listar(empresaID)
	listaGrupos, _ := h.serviceGrupo.Listar(empresaID)
	listaEquipamentos, _ := h.serviceEquipamento.Listar(empresaID)

	data := TelaFuncionarioData{
		Funcionarios:      listaFuncionarios,
		GruposDispo:       listaGrupos,
		EquipamentosDispo: listaEquipamentos,
		IsErro:            true,
		Mensagem:          msg,
	}
	tmplFuncionario.Execute(w, data)
}
