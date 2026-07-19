package handlers

import (
	"fmt"
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

// Carrega o HTML apenas uma vez na inicialização do sistema para alta performance
var tmplUsuario = template.Must(template.ParseFiles("web/templates/usuario.html"))

type UsuarioHandler struct {
	service *services.UsuarioService
}

// NewUsuarioHandler cria a instancia do controlador injetando o service de banco
func NewUsuarioHandler(service *services.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{service: service}

}

// Tabela Usuario Dados define o formato de dados que a tela html espera receber

type DadosTelaUsuario struct {
	Usuarios []models.Usuario
	Mensagem string
	IsErro   bool
}

// GeraUsuario centraliza as operações de exibição e processamento da tela

func (u *UsuarioHandler) GerenciarUsuario(w http.ResponseWriter, r *http.Request) {
	// Simulação de sessão: Usuario logado pertence a empresa Id 1 (se for superuser, deixe nil)
	var empresaLogadaID uint = 1

	/* === MÉTODO POST === */

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		fmt.Printf("Ação recebida no Go: %s\n", acao)

		// Captura os dados textuais do formulario

		nomeInput := r.FormValue("nome")
		emailInput := r.FormValue("email")
		usuarioInput := r.FormValue("usuario")
		senhaInput := r.FormValue("senha")
		permissaoInput := r.FormValue("permissao")
		ativoInput := r.FormValue("ativo")

		// 1. AÇÃO INCLUIR

		if acao == "incluir" {
			novoUs := models.Usuario{
				Nome:      nomeInput,
				Email:     emailInput,
				Usuario:   usuarioInput,
				Senha:     senhaInput,
				Permissao: &permissaoInput,
				Ativo:     &ativoInput,
			}

			if err := u.service.Incluir(&novoUs); err != nil {
				u.renderizarComErro(w, "Erro ao incluir usuario: "+err.Error(), &empresaLogadaID)
				return

			}

		}
		// 2. AÇÃO: ALTERAR

		if acao == "alterar" {
			var idUser uint
			fmt.Sscanf(r.FormValue("id_user"), "%d", &idUser)
			usEditado := models.Usuario{

				Nome:      nomeInput,
				Email:     emailInput,
				Senha:     senhaInput,
				Permissao: &permissaoInput,
				Ativo:     &ativoInput,
			}

			if err := u.service.Alterar(&usEditado, &empresaLogadaID); err != nil {
				u.renderizarComErro(w, "Erro ao alterar Usuario "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// 3. AÇÃO: EXCLUIR

		if acao == "excluir" {
			var idUser uint
			fmt.Sscanf(r.FormValue("id_user"), "%d", &idUser)

			if err := u.service.Excluir(idUser, &empresaLogadaID); err != nil {
				u.renderizarComErro(w, "Erro ao excluir usuario: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// Recarrega a pagina limpando os inputs após o sucesso (Igual o redirect do Django)
		http.Redirect(w, r, "/usuario", http.StatusSeeOther)
		return
	}

	/* === MÉTODO GET === */

	lista, err := u.service.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao buscar usuarios", http.StatusInternalServerError)
		return
	}
	tmplUsuario.Execute(w, DadosTelaUsuario{Usuarios: lista})
}

// Função auxiliar para renderizar a pagina mostrando o erro

func (u *UsuarioHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := u.service.Listar(empID)
	tmplUsuario.Execute(w, DadosTelaUsuario{
		Usuarios: lista,
		Mensagem: msg,
		IsErro:   true,
	})
}
