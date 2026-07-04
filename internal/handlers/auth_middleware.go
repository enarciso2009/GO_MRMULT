package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta que deve ser IDENTICA  a do seu arquivo auth_service.go
var jwtChaveSecretaVerificacao = []byte("SUA_CHAVE_SECRETA_SUPER_COMPLEXA_E_LONGA_123456")

// RequererAutenticacao é o cadeado que faz o papel do seu LoginRequiredMixin do Django
func RequererAutenticacao(proximoHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1. Tenta buscar o cookie de sessão no navegador do usuario
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// Se o cookie Não existir, expulsa o usuario para a tela de login na hora!
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// 2. Se o cookie existir, extrai o texto do Token de dentro dele
		tokenString := cookie.Value

		// 3. Valida se o Token é verdadeiro e não foi falsificado por um hacker
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtChaveSecretaVerificacao, nil
		})

		// 4. Se o Token estiver quebrado, vencido ou com erro de leitura
		if err != nil || !token.Valid {
			// Redireciona o usuario para a tela de Acesso Negado
			http.Redirect(w, r, "/acesso-negado", http.StatusSeeOther)
			return
		}

		// 5. Se passou por todas as travas, o usuario esta logado! Libera o acesso para a pagina
		proximoHandler(w, r)
	}
}
