package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

var tmplCadTer = template.Must(template.ParseFiles("web/templates/cadter.html"))

//CadVis agrupa as variaveis que o html vai renderizar

type CadTerHandler struct {
	cadTerService *services.CadTerService
}

func NewCadTerHandler(service *services.CadTerService) *CadTerHandler {
	return &CadTerHandler{cadTerService: service}
}

// DadosCadVisTela leva as informacoes consolidadas para a tela html
type DadosCadTerTela struct {
	Ter      []models.Terceiro
	Mensagem string
	IsErro   bool
}

func (ct *CadTerHandler) ExibirCadTer(w http.ResponseWriter, r *http.Request) {
	var empresaLogadaID uint = 1

	lista, err := ct.cadTerService.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao carregar dados do Visitante", http.StatusInternalServerError)
		return
	}

	tmplCadTer.Execute(w, DadosCadTerTela{Ter: lista})
}

// Funcao auxiliar para renderizar a pagina mostrando erro

func (ct *CadTerHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := ct.cadTerService.Listar(empID)
	tmplCadVis.Execute(w, DadosCadTerTela{
		Ter:      lista,
		Mensagem: msg,
		IsErro:   true,
	})
}
