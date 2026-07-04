package handlers

import (
	"encoding/json"
	"html/template"
	"mrmult/internal/services"
	"net/http"
)

var tmplHome = template.Must(template.ParseFiles("web/templates/home.html"))

// DadosHome agrupa as variáveis que o HTML vai renderizar
type HomeHandler struct {
	dashboardService *services.DashboardService
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{
		dashboardService: services.NewDashboardService(),
	}
}

// DadosHome leva as informações consolidadas para a tela html
type DadosHome struct {
	UsuarioNome string
	TotalFunc   int64
	TotalVisit  int64
	TotalTerc   int64
	TotalGeral  int64
}

func (h *HomeHandler) ExibirHome(w http.ResponseWriter, r *http.Request) {
	// [Simulação do Mixin]: Se o usuario nao estivesse logado, fariamos isso:
	// http.Redirect(w, r, "/login", http.StatusSeeOther)

	// Simulação de uma empresa logada (ID 1). Se for superuser, deixamos como nil
	var empresaLogadaID uint = 1

	// 1. Buca os totais reais contando direto do banco de dados via GORM
	tFunc, tVisit, tTerc, err := h.dashboardService.BuscarTotaisDoDia(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao carregar dados do painel", http.StatusInternalServerError)
		return
	}

	// 2. [SUPORTE AJAX]: Se a Url tiver ?ajax=true, responde com JSON puro igual ao Django
	if r.URL.Query().Get("ajax") != "" {
		w.Header().Set("Content-Type", "application/json")
		// Cria o mapa identico ao JsonResponse do Python
		json.NewEncoder(w).Encode(map[string]int64{
			"total_func":  tFunc,
			"total_visit": tVisit,
			"total_terc":  tTerc,
		})
		return
	}

	// 3. Renderização normal da pagina HTML caso não sejja uma requisição AJAX
	dados := DadosHome{
		UsuarioNome: "Everton Narciso",
		TotalFunc:   tFunc,
		TotalVisit:  tVisit,
		TotalTerc:   tTerc,
		TotalGeral:  tFunc + tVisit + tTerc,
	}

	tmplHome.Execute(w, dados)

}
