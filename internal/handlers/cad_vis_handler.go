package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

var tmplCadVis = template.Must(template.ParseFiles("web/templates/cadvis.html"))

//CadVis agrupa as variaveis que o html vai renderizar

type CadVisHandler struct {
	cadVisService *services.CadVisService
}

func NewCadVisHandler(service *services.CadVisService) *CadVisHandler {
	return &CadVisHandler{cadVisService: service}
}

// DadosCadVisTela leva as informacoes consolidadas para a tela html
type DadosCadVisTela struct {
	Vis      []models.Visitante
	Mensagem string
	IsErro   bool
}

func (cv *CadVisHandler) ExibirCadVis(w http.ResponseWriter, r *http.Request) {
	var empresaLogadaID uint = 1

	lista, err := cv.cadVisService.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao carregar dados do Visitante", http.StatusInternalServerError)
		return
	}

	tmplCadVis.Execute(w, DadosCadVisTela{Vis: lista})
}

// Funcao auxiliar para renderizar a pagina mostrando erro

func (cv *CadVisHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := cv.cadVisService.Listar(empID)
	tmplCadVis.Execute(w, DadosCadVisTela{
		Vis:      lista,
		Mensagem: msg,
		IsErro:   true,
	})
}
