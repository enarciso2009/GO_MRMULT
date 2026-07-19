package handlers

import (
	"fmt"
	"html/template"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

// Carrega o html apenas uma vez na inicialização do sistema para alta performance

var tmplEmpresa = template.Must(template.ParseFiles("web/templates/empresa.html"))

type EmpresaHandler struct {
	service *services.EmpresaService
}

// NewEmpresa cria a instancia do controlador injetando o service de banco
func NewEmpresaHandler(service *services.EmpresaService) *EmpresaHandler {
	return &EmpresaHandler{service: service}
}

// tabela Empresa Dados define o formato que a tela html espera receber

type DadosTelaEmpresa struct {
	Empresas []models.Empresa
	Mensagem string
	IsErro   bool
}

// Gera Empresa centraliza as operações de exibição e processamento da tela

func (e *EmpresaHandler) GerenciarEmpresa(w http.ResponseWriter, r *http.Request) {
	// Simulação de sessão: Usuario logado pertence a empresa ID 1 ( se for superUser, deixe nil)
	var empresaLogadaID uint = 1

	/* === MÉTODO POST === */

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		fmt.Printf("Ação recebida no Go: %s\n", acao)

		// Captura os dados textuais do formulario

		nomeInput := r.FormValue("nome")
		cnpjInput := r.FormValue("cnpj")
		enderecoInput := r.FormValue("endereco")
		cepInput := r.FormValue("cep")
		bairroInput := r.FormValue("bairro")
		municipioInput := r.FormValue("municipio")
		estadoInput := r.FormValue("estado")

		// 1. AÇÃO INCLUIR

		if acao == "incluir" {
			novoEm := models.Empresa{
				Nome:      nomeInput,
				CNPJ:      cnpjInput,
				Endereco:  enderecoInput,
				Cep:       cepInput,
				Bairro:    bairroInput,
				Municipio: municipioInput,
				Estado:    estadoInput,
			}

			if err := e.service.Incluir(&novoEm); err != nil {
				e.renderizarComErro(w, "Erro ao incluir Empresa: "+err.Error(), &empresaLogadaID)
				return
			}
		}
		// 2. AÇÃO: ALTERAR

		if acao == "alterar" {
			var idEmpr uint
			fmt.Sscanf(r.FormValue("id_empr"), "%d", &idEmpr)
			emEditado := models.Empresa{
				Nome:      nomeInput,
				CNPJ:      cnpjInput,
				Endereco:  enderecoInput,
				Cep:       cepInput,
				Bairro:    bairroInput,
				Municipio: municipioInput,
				Estado:    estadoInput,
			}

			if err := e.service.Alterar(&emEditado, &empresaLogadaID); err != nil {
				e.renderizarComErro(w, "Erro ao alterar Empresa "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// 3. AÇÃO: EXCLUIR

		if acao == "excluir" {
			var idEmpr uint
			fmt.Sscanf(r.FormValue("id_empr"), "%d", &idEmpr)

			if err := e.service.Excluir(idEmpr, &empresaLogadaID); err != nil {
				e.renderizarComErro(w, "Erro ao excluir Empresa: "+err.Error(), &empresaLogadaID)
				return
			}
		}

		// recarrega a pagina limpando os inputs após o sucesso (Igual o redirect do Django)
		http.Redirect(w, r, "/empresa", http.StatusSeeOther)
		return
	}

	/* === MÉTODO GET === */

	lista, err := e.service.Listar(&empresaLogadaID)
	if err != nil {
		http.Error(w, "Erro ao buscar Empresas", http.StatusInternalServerError)
		return
	}
	tmplEmpresa.Execute(w, DadosTelaEmpresa{Empresas: lista})
}

// Função auxiliar para renderizar a pagina mostrando o erro

func (e *EmpresaHandler) renderizarComErro(w http.ResponseWriter, msg string, empID *uint) {
	lista, _ := e.service.Listar(empID)
	tmplEmpresa.Execute(w, DadosTelaEmpresa{
		Empresas: lista,
		Mensagem: msg,
		IsErro:   true,
	})
}
