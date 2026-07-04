package handlers

import (
	"html/template"
	"net/http"
)

// Faz o parse uma unica vez na inicializacao do sistema (unindo a tela com a base)
var tmplSobre = template.Must(template.ParseFiles("web/templates/sobre.html"))

type SobreHandler struct{}

// NewSobreHandler cria uma nova instancia do controlador
func NewSobreHandler() *SobreHandler {
	return &SobreHandler{}
}

// ExibirSobre faz o parse do HTML e renderiza a tela
func (h *SobreHandler) ExibirSobre(w http.ResponseWriter, r *http.Request) {
	// executando o template para a chamada da pagina sobre.html
	err := tmplSobre.Execute(w, nil)
	if err != nil {
		http.Error(w, "Erro ao renderizar a tela: "+err.Error(), http.StatusInternalServerError)
	}

}
