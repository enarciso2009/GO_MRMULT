package main

import (
	"fmt"
	_ "mrmult/docs"
	"mrmult/internal/database"
	"mrmult/internal/handlers"
	"mrmult/internal/models"
	"mrmult/internal/services"
	"net/http"
)

// @title           API de Integração - Módulo Refeitório
// @version         1.0
// @description     API para sincronismo de dados entre o Controle de Acesso e o Refeitório.
// @BasePath        /

func main() {

	// Conexão com o banco de dados
	db, err := database.Conectar()
	if err != nil {
		panic(err)
	}

	// O Automigrate verifica a struct e cria/atualiza a tabela no banco de dados sozinho
	err = db.AutoMigrate(
		&models.Empresa{},
		&models.Refeicao{},
		&models.PrecoRefeicao{},
		&models.GrupoRefeicao{},
		&models.InterGrupRef{},
		&models.Funcionario{},
		&models.Visitante{},
		&models.Terceiro{},
		&models.Equipamento{},
		&models.Evento{},
		&models.Usuario{},
		&models.Parametro{},
	)
	if err != nil {
		panic("Falha ao rodar as migrações automaticas: " + err.Error())
	}
	fmt.Println("Banco de dados sincronizado com sucesso!")

	// INSTANCIA OS CONTROLADORES

	homeHandler := handlers.NewHomeHandler()

	authHandler := handlers.NewAuthHandler()

	equipamentoHandler := handlers.NewEquipamentoHandler()

	// 1. Instancia os servicos necessario (camada de negocios)
	refeicaoService := services.NewRefeicaoService()

	grupoRefeicaoService := services.NewGrupoRefeicaoService()

	funcionarioService := services.NewFuncionarioService()

	equipamentoService := services.NewEquipamentoService()

	visitanteService := services.NewVisitanteService()

	terceiroService := services.NewTerceiroService()

	// 2. Instancia os Handlers passando suas dependencias explicitas
	refeicaoHandler := handlers.NewRefeicaoHandler(refeicaoService)

	sobreHandler := handlers.NewSobreHandler()

	grupoRefeicaoHandler := handlers.NewGrupoRefeicaoHandler(grupoRefeicaoService, refeicaoService)

	funcionarioHandler := handlers.NewFuncionarioHandler(funcionarioService, grupoRefeicaoService, equipamentoService)

	visitanteHandler := handlers.NewVisitanteHandler(visitanteService, funcionarioService, grupoRefeicaoService, equipamentoService)

	terceiroHandler := handlers.NewTerceiroHandler(terceiroService, funcionarioService, grupoRefeicaoService, equipamentoService)

	apiHandler := handlers.NewAPIIntegrationHandler(funcionarioService, visitanteService, terceiroService)

	// Rotas

	// 1 Cria um roteador isolado (O seu aruvo urls.py do Django)
	roteador := http.NewServeMux()

	// 2. Rota de arquivos estaticos (CSS, imagens)
	fs := http.FileServer(http.Dir("web/static"))
	roteador.Handle("/static/", http.StripPrefix("/static", fs))

	// Rotas da API de Integração (O sistema de controle de acesso vai disparar POSTs aqui)
	roteador.HandleFunc("POST /api/integracao/funcionario", apiHandler.ReceberFuncionario)
	roteador.HandleFunc("POST /api/integracao/visitante", apiHandler.ReceberVisitante)
	roteador.HandleFunc("POST /api/integracao/terceiro", apiHandler.ReceberTerceiro)

	// 1. Rota que entrega o arquivo JSON gerado pelo "swag init"
	roteador.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	// 2. Rota que entrega a tela do Swagger usando o arquivo HTML que criamos acima
	roteador.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/swagger/index.html")
	})
	// Redirecionamento caso acessem sem a barra no final
	roteador.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	// Adicionamos o "GET" ou "POST" antes para o Go isolar os metodos (recurso do Go Moderno)

	roteador.HandleFunc("GET /login", authHandler.ExibirLogin)
	roteador.HandleFunc("POST /login", authHandler.Login)

	roteador.HandleFunc("GET /acesso_negado", authHandler.ExibirAcessoNegado)
	roteador.HandleFunc("GET /logout", authHandler.Logout)

	roteador.HandleFunc("GET /sobre", sobreHandler.ExibirSobre)

	roteador.HandleFunc("GET /grupo_refeicao", grupoRefeicaoHandler.GerenciarGrupo)
	roteador.HandleFunc("POST /grupo_refeicao", grupoRefeicaoHandler.GerenciarGrupo)

	roteador.HandleFunc("GET /{$}", handlers.RequererAutenticacao(homeHandler.ExibirHome))

	roteador.HandleFunc("GET /equipamento", handlers.RequererAutenticacao(equipamentoHandler.GerenciarEquipamento))
	roteador.HandleFunc("POST /equipamento", handlers.RequererAutenticacao(equipamentoHandler.GerenciarEquipamento))

	roteador.HandleFunc("GET /refeicao", handlers.RequererAutenticacao(refeicaoHandler.GerenciarRefeicao))
	roteador.HandleFunc("POST /refeicao", handlers.RequererAutenticacao(refeicaoHandler.GerenciarRefeicao))

	roteador.HandleFunc("POST /funcionario", funcionarioHandler.GerenciarFuncionario)
	roteador.HandleFunc("GET /funcionario", funcionarioHandler.GerenciarFuncionario)

	roteador.HandleFunc("POST /visitante", visitanteHandler.GerenciarVisitante)
	roteador.HandleFunc("GET /visitante", visitanteHandler.GerenciarVisitante)

	roteador.HandleFunc("POST /terceiro", terceiroHandler.GerenciarTerceiro)
	roteador.HandleFunc("GET /terceiro", terceiroHandler.GerenciarTerceiro)

	fmt.Println("Servidor rodando com Dashboard em http://localhost:8080")
	err = http.ListenAndServe(":8080", roteador)
	if err != nil {
		panic(err)
	}
}
