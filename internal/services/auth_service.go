package services

import (
	"errors"
	"mrmult/internal/database"
	"mrmult/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Chave secreta ultra-segura para assinar o token (Em Produção, use usa variavel de ambiente
var jwtChaveSecreta = []byte("SUA_CHAVE_SECRETA_SUPER_COMPLEXA_E_LONGA_123456")

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// HashSenha transforma a senha pura em um bloco criptografado impenetravel
func (s *AuthService) HashSenha(senhaPura string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senhaPura), 12) // Custo 12 é o padrão de alta segurança
	return string(bytes), err
}

// RealizarLogin valida o usuario no banco e gera o Token JWT seguro
func (s *AuthService) RealizarLogin(login, senhaPura string) (string, error) {
	db, err := database.Conectar()
	if err != nil {
		return "", err
	}

	var usuario models.Usuario
	// Busca o usuario pelo campo 'usuario' da tabela
	result := db.Where("usuario = ?", login).First(&usuario)
	if result.Error != nil {
		return "", errors.New("usuario ou senha incorretos")
	}

	// Compara a senha digitada com Hash criptografado que está no banco de dados
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senhaPura))
	if err != nil {
		return "", errors.New("usuario ou senha incorretos")
	}

	// Cria o Token JWT com tempo de expiração (Ex: Expira em 8 horas)
	claims := jwt.MapClaims{
		"user_id": usuario.IDUser,
		"nome":    usuario.Nome,
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtChaveSecreta)

	return tokenString, err

}

// Autenticar verifica as credenciais do usuario no banco de dados

func (s *AuthService) Autenticar(username, senha string) (*models.Usuario, error) {
	db, err := database.Conectar()
	if err != nil {
		return nil, err
	}

	var usuario models.Usuario
	// Busca o usuario pelo campo "Usuario" trazendo a empresa junto Preload
	err = db.Model(&models.Usuario{}).Where("usuario = ?", username).First(&usuario).Error

	if err != nil {
		// Retorna um erro generico amigavel por questoes de seguranca
		return nil, errors.New("Usuario ou senha invalidos")

	}
	// Como suas senhas agora usam Bcrypt, precisamos validar aqui tambem!
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senha))
	if err != nil {
		return nil, errors.New("usuario ou senha invalidos")
	}
	// retorna o ponteiro do usuario encontrado (com o campo Permissao populado)
	return &usuario, nil
}
