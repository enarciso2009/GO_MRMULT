package main

import (
	"fmt"
	"mrmult/internal/database"
	"mrmult/internal/handlers"
	"mrmult/internal/models"
	"net/http"
)

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
	refeicaoHandler := handlers.NewRefeicaoHandler()
	sobreHandler := handlers.NewSobreHandler()

	// Rotas

	// 1 Cria um roteador isolado (O seu aruvo urls.py do Django)
	roteador := http.NewServeMux()

	// 2. Rota de arquivos estaticos (CSS, imagens)
	fs := http.FileServer(http.Dir("web/static"))
	roteador.Handle("/static/", http.StripPrefix("/static", fs))

	// 3. Rotas publicas (Qualquer um pode acessar sem login)
	// Adicionamos o "GET" ou "POST" antes para o Go isolar os metodos (recurso do Go Moderno)
	roteador.HandleFunc("GET /login", authHandler.ExibirLogin)
	roteador.HandleFunc("POST /login", authHandler.Login)
	roteador.HandleFunc("GET /acesso_negado", authHandler.ExibirAcessoNegado)
	roteador.HandleFunc("GET /logout", authHandler.Logout)
	roteador.HandleFunc("GET /sobre", sobreHandler.ExibirSobre)

	// 4. ROTA PROTEGIDA (O Middleware vai inspecionar o acesso antes de liberar)
	// Repare o homeHandler.ExibirHome envelopado pelo RequererAutenticacao
	roteador.HandleFunc("GET /{$}", handlers.RequererAutenticacao(homeHandler.ExibirHome))
	roteador.HandleFunc("GET /equipamento", handlers.RequererAutenticacao(equipamentoHandler.GerenciarEquipamento))
	roteador.HandleFunc("POST /equipamento", handlers.RequererAutenticacao(equipamentoHandler.GerenciarEquipamento))
	roteador.HandleFunc("GET /refeicao", handlers.RequererAutenticacao(refeicaoHandler.GerenciarRefeicao))
	roteador.HandleFunc("POST /refeicao", handlers.RequererAutenticacao(refeicaoHandler.GerenciarRefeicao))

	fmt.Println("Servidor rodando com Dashboard em http://localhost:8080")
	err = http.ListenAndServe(":8080", roteador)
	if err != nil {
		panic(err)
	}
}
