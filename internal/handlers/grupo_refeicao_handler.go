package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
	"strconv"
)

var tmplGrupoRefeicao = template.Must(template.ParseFiles("web/templates/grupo_refeicao.html"))

// TelaGrupoData envia para a tela os grupos cadastrados e as refeicoes disponiveis para marcar
type TelaGrupoData struct {
	Grupos         []models.GrupoRefeicao
	RefeicoesDispo []models.Refeicao
	IsErro         bool
	Mensagem       string
}

type GrupoRefeicaoHandler struct {
	serviceGrupo    *services.GrupoRefeicaoService
	serviceRefeicao *services.RefeicaoService
}

func NewGrupoRefeicaoHandler(sg *services.GrupoRefeicaoService, sr *services.RefeicaoService) *GrupoRefeicaoHandler {
	return &GrupoRefeicaoHandler{serviceGrupo: sg, serviceRefeicao: sr}
}

// Metodo auxiliar para verificar se uma refeicao epecifica esta dentro de um grupo (usado no js da tela)
// internal/handlers/grupo_refeicao_handler.go

func (h *GrupoRefeicaoHandler) GerenciarGrupo(w http.ResponseWriter, r *http.Request) {
	var empresaPadraoID uint = 1

	// ==========================================
	// OPERAÇÃO 1: PROCESSAMENTO DE FORMULÁRIO (POST)
	// ==========================================
	if r.Method == http.MethodPost {
		// 1. Força o carregamento total do corpo da requisição primeiro
		err := r.ParseForm()
		if err != nil {
			h.renderizarComErro(w, &empresaPadraoID, "Erro ao processar formulário: "+err.Error())
			return
		}

		// 2. MUDANÇA ESSENCIAL: Lemos direto de r.PostForm em vez de r.FormValue
		acao := r.PostForm.Get("acao")
		idGrupStr := r.PostForm.Get("id_grup")
		nome := r.PostForm.Get("nome")

		// 3. Captura o array completo de checkboxes marcados
		refeicoesSelecionadasStr := r.PostForm["refeicoes"]

		var idsRefeicoes []uint
		for _, idStr := range refeicoesSelecionadasStr {
			idUint, _ := strconv.ParseUint(idStr, 10, 32)
			idsRefeicoes = append(idsRefeicoes, uint(idUint))
		}

		// PRINT DE SEGURANÇA: Esse contador agora VAI mudar para maior que 0!
		println("DEBUG - IDs de refeicoes capturadas no formulario:", len(idsRefeicoes))

		if acao == "incluir" || acao == "alterar" {
			var idGrup uint = 0
			if acao == "alterar" {
				idUint, _ := strconv.ParseUint(idGrupStr, 10, 32)
				idGrup = uint(idUint)
			}

			grupo := models.GrupoRefeicao{
				IDGrup:    idGrup,
				Nome:      &nome,
				EmpresaID: &empresaPadraoID,
			}

			// Passa os IDs para a camada de serviços fazer o vínculo no Postgres
			err := h.serviceGrupo.Salvar(&grupo, idsRefeicoes)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao salvar grupo: "+err.Error())
				return
			}
		}

		if acao == "excluir" {
			idUint, _ := strconv.ParseUint(idGrupStr, 10, 32)
			err := h.serviceGrupo.Excluir(uint(idUint), &empresaPadraoID)
			if err != nil {
				h.renderizarComErro(w, &empresaPadraoID, "Erro ao excluir grupo: "+err.Error())
				return
			}
		}

		http.Redirect(w, r, "/grupo_refeicao", http.StatusSeeOther)
		return
	}

	// ==========================================
	// OPERAÇÃO 2: EXIBIÇÃO DA TELA (GET)
	// ==========================================
	listaGrupos, errG := h.serviceGrupo.Listar(&empresaPadraoID)
	listaRefeicoes, errR := h.serviceRefeicao.Listar(&empresaPadraoID)

	data := TelaGrupoData{
		Grupos:         listaGrupos,
		RefeicoesDispo: listaRefeicoes,
	}

	if errG != nil || errR != nil {
		data.IsErro = true
		data.Mensagem = "Erro ao carregar dados do Postgres."
	}

	tmplGrupoRefeicao.Execute(w, data)
}

func (h *GrupoRefeicaoHandler) renderizarComErro(w http.ResponseWriter, empresaID *uint, msg string) {
	listaGrupos, _ := h.serviceGrupo.Listar(empresaID)
	listaRefeicoes, _ := h.serviceRefeicao.Listar(empresaID)
	data := TelaGrupoData{
		Grupos:         listaGrupos,
		RefeicoesDispo: listaRefeicoes,
		IsErro:         true,
		Mensagem:       msg,
	}
	tmplGrupoRefeicao.Execute(w, data)
}
