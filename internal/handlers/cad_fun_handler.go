package handlers

import (
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

var tmplCadFun = template.Must(template.ParseFiles("web/templates/cadfun.html"))

// CadFun agrupa as variaveis que o HTML vai renderizar

type CadFunHandler struct {
	cadFunService *services.CadFunService
}

func NewCadFunHandler(service *services.CadFunService) *CadFunHandler {
	return &CadFunHandler{cadFunService: service}

}

// DadosCadFunTela leva as informacoes consolidadas para a tela html
type DadosCadFunTela struct {
	Fun      []models.Funcionario
	Mensagem string
	IsErro   bool
}

func (cf *CadFunHandler) ExibirCadFun(w http.ResponseWriter, r *http.Request) {
	var empresaLogadaID uint = 1

	lista, err := cf.cadFunService.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao carregar dados do funcionario", http.StatusInternalServerError)
		return
	}

	tmplCadFun.Execute(w, DadosCadFunTela{Fun: lista})

}

// Funcao auxiliar para renderizar a pagina mostrando erro

func (cf *CadFunHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := cf.cadFunService.Listar(empID)
	tmplCadFun.Execute(w, DadosCadFunTela{
		Fun:      lista,
		Mensagem: msg,
		IsErro:   true,
	})
}
