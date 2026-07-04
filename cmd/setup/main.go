package main

import (
	"fmt"
	"mrmult/internal/database"
	"mrmult/internal/models"
	"mrmult/internal/services"
)

func main() {
	db, err := database.Conectar()
	if err != nil {
		panic("Erro ao conectar no banco: " + err.Error())
	}

	fmt.Println("⏳ Inicializando banco com Empresa Master e Usuário Admin...")

	// 1. Cria a Empresa Master (Ela receberá o ID 1 automaticamente por ser a primeira)
	empresaMaster := models.Empresa{
		Nome: "Sistema Master - Administração",
		CNPJ: "00.000.000/0001-00",
	}
	if err := db.Create(&empresaMaster).Error; err != nil {
		panic("Erro ao criar empresa master: " + err.Error())
	}
	fmt.Printf("✔ Empresa Master criada com sucesso! ID: %d\n", empresaMaster.ID)

	// 2. Criptografa a senha do administrador
	authService := services.NewAuthService()
	senhaCriptografada, err := authService.HashSenha("admin123")
	if err != nil {
		panic("Erro ao gerar hash da senha: " + err.Error())
	}

	// 3. Cria o usuário Admin vinculado à Empresa Master (ID 1)
	usuarioAdmin := models.Usuario{
		IDUser:    "USR-001",
		Nome:      "Administrador do Sistema",
		Email:     "admin@empresa.com",
		Usuario:   "admin",
		Senha:     senhaCriptografada,
		Permissao: nil,
		EmpresaID: &empresaMaster.ID, // Vincula ao ID 1 que acabou de nascer
	}

	err = db.Create(&usuarioAdmin).Error
	if err != nil {
		panic("Erro ao salvar usuário no banco: " + err.Error())
	}

	fmt.Println("✅ Banco configurado com sucesso!")
}
