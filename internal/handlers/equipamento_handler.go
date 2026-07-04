package handlers

import (
	"fmt"
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

var tmplEquipamento = template.Must(template.ParseFiles("web/templates/equipamento.html"))

type EquipamentoHandler struct {
	service *services.EquipamentoService
}

func NewEquipamentoHandler() *EquipamentoHandler {
	return &EquipamentoHandler{
		service: services.NewEquipamentoService(),
	}
}

// DadosTelaEquipamento monta o contexto de renderização

type DadosTelaEquipamento struct {
	Equipamentos []models.Equipamento
	Mensagem     string
	IsErro       bool
}

func (h *EquipamentoHandler) GerenciarEquipamento(w http.ResponseWriter, r *http.Request) {
	//Simulação de sessão: Usuario logado pertence a empresa ID 1 (se for superuser, deixe nil)
	var empresaLogadaID uint = 1

	/* === MÉTODO POST === */

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		fmt.Printf("Ação recebida no Go: %s\n", acao)

		// Captura os dados textuais do formulario
		nomeInput := r.FormValue("nome")
		ipInput := r.FormValue("ip")
		maskInput := r.FormValue("mask")

		// 1. AÇÃO INCLUIR

		if acao == "incluir" {
			novoEq := models.Equipamento{
				Nome:      &nomeInput,
				IP:        &ipInput,
				Mask:      &maskInput,
				EmpresaID: &empresaLogadaID,
			}

			if err := h.service.Incluir(&novoEq); err != nil {
				h.renderizarComErro(w, "Erro ao incluir equipamento: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// 2. AÇÃO: ALTERAR

		if acao == "alterar" {
			var idEquip uint
			fmt.Sscanf(r.FormValue("id_equip"), "%d", &idEquip)

			eqEditado := models.Equipamento{
				IDEquip: idEquip,
				Nome:    &nomeInput,
				IP:      &ipInput,
				Mask:    &maskInput,
			}

			if err := h.service.Alterar(&eqEditado, &empresaLogadaID); err != nil {
				h.renderizarComErro(w, "Erro ao alterar equipamento: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// 3. AÇÃO: EXCLUIR

		if acao == "excluir" {
			var idEquip uint
			fmt.Sscanf(r.FormValue("id_equip"), "%d", &idEquip)

			if err := h.service.Excluir(idEquip, &empresaLogadaID); err != nil {
				h.renderizarComErro(w, "Erro ao excluir equipamento: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// Recarrega a pagina limpando os inputs apos o sucesso (Igual o redirect do Django)
		http.Redirect(w, r, "/equipamento", http.StatusSeeOther)
		return
	}

	/* === METÓDO GET === */

	lista, err := h.service.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao buscar equipamentos", http.StatusInternalServerError)
		return
	}

	tmplEquipamento.Execute(w, DadosTelaEquipamento{Equipamentos: lista})

}

// Função auxiliar para renderizar a pagina mostrando o erro

func (h *EquipamentoHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := h.service.Listar(empID)
	tmplEquipamento.Execute(w, DadosTelaEquipamento{
		Equipamentos: lista,
		Mensagem:     msg,
		IsErro:       true,
	})
}
