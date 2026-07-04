package handlers

import (
	"fmt"
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
	"time"
)

var tmplRefeicao = template.Must(template.ParseFiles("web/templates/refeicao.html"))

type RefeicaoHandler struct {
	service *services.RefeicaoService
}

func NewRefeicaoHandler() *RefeicaoHandler {
	return &RefeicaoHandler{
		service: services.NewRefeicaoService(),
	}
}

type DadosTelaRefeicao struct {
	Refeicoes []models.Refeicao
	Mensagem  string
	IsErro    bool
}

func (h *RefeicaoHandler) GerenciarRefeicao(w http.ResponseWriter, r *http.Request) {
	var empresaLogadaID uint = 1 // Vinculado a empresa padrão do admin

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		nomeInput := r.FormValue("nome")

		var valorInput float64
		fmt.Sscanf(r.FormValue("valor"), "%f", &valorInput)
		// Capture as strings de data vindas do HTML (formato padrão "AAAA-MM-DD")

		horaInicio := r.FormValue("hora_inicio")
		horaFim := r.FormValue("hora_fim")

		// Crie o layout que ensina o Go a ler esse formato de data
		layoutData := "2006-01-02"

		// Converte data inicio obrigatorio
		dInicio, err1 := time.Parse(layoutData, r.FormValue("data_inicio"))
		var ponteiroDataInicio *time.Time
		if err1 == nil {
			ponteiroDataInicio = &dInicio
		}

		//Converte a data de fim (Opcional)
		dataFimStr := r.FormValue("data_fim")
		var ponteiroDataFim *time.Time
		// Só tentar converter se realmente o usuario escolheu uma data na tela
		if dataFimStr != "" {
			dFim, err2 := time.Parse(layoutData, dataFimStr)
			if err2 == nil {
				ponteiroDataFim = &dFim
			}
		}

		// Se dataFimStr for vazia o ponteiroDataFim continuará sendo 'nil'.
		// Fazendo o GORM gravar como null perfeitamente no banco de dados!

		if acao == "incluir" {
			novaRef := models.Refeicao{
				Nome:       nomeInput,
				Valor:      valorInput,
				DataInicio: ponteiroDataInicio,
				DataFim:    ponteiroDataFim,
				HoraInicio: &horaInicio,
				HoraFim:    &horaFim,
				EmpresaID:  &empresaLogadaID,
			}
			if err := h.service.Incluir(&novaRef); err != nil {
				h.renderizarComErro(w, "Erro ao incluir refeição: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		if acao == "alterar" {
			var idRef uint
			fmt.Sscanf(r.FormValue("id_ref"), "%d", &idRef)

			refEditada := models.Refeicao{
				IDRef:      idRef,
				Nome:       nomeInput,
				Valor:      valorInput,
				DataInicio: ponteiroDataInicio,
				DataFim:    ponteiroDataFim,
				HoraInicio: &horaInicio,
				HoraFim:    &horaFim,
			}
			if err := h.service.Alterar(&refEditada, &empresaLogadaID); err != nil {
				h.renderizarComErro(w, "Erro ao alterar refeição: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		if acao == "excluir" {
			var idRef uint
			fmt.Sscanf(r.FormValue("id_ref"), "%d", &idRef)
			if err := h.service.Excluir(idRef, &empresaLogadaID); err != nil {
				h.renderizarComErro(w, "Erro ao excluir refeição: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		http.Redirect(w, r, "/refeicao", http.StatusSeeOther)
		return
	}

	lista, err := h.service.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao buscar refeições", http.StatusInternalServerError)
		return
	}

	tmplRefeicao.Execute(w, DadosTelaRefeicao{Refeicoes: lista})
}

func (h *RefeicaoHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := h.service.Listar(empID)
	tmplRefeicao.Execute(w, DadosTelaRefeicao{Refeicoes: lista, Mensagem: msg, IsErro: true})
}
