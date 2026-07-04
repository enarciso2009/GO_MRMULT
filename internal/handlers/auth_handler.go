package handlers

import (
	"html/template"
	"mrmult/internal/services"
	"net/http"
	"time"
)

var tmplLogin = template.Must(template.ParseFiles("web/templates/login.html"))
var tmplAcessoNegado = template.Must(template.ParseFiles("web/templates/acesso_negado.html"))
var tmplLogout = template.Must(template.ParseFiles("web/templates/logout.html"))

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// ExibirLogin renderiza a tela de login
func (h *AuthHandler) ExibirLogin(w http.ResponseWriter, r *http.Request) {
	tmplLogin.Execute(w, nil)
}

// Login processa o formulario de entrada
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	usuarioInput := r.FormValue("username")
	senhaInput := r.FormValue("password")

	token, err := h.authService.RealizarLogin(usuarioInput, senhaInput)
	if err != nil {
		// Se falhar, recarrega a pagina de login com uma mensagem de erro
		tmplLogin.Execute(w, "Usuario ou senha invalidos!")
		return
	}

	// A BLINDAGEM DO COOKIE: Configura as travas de segurança maxima
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 8),
		Path:     "/",
		HttpOnly: true,                    // Protege contra roubo via Javascript (XSS)
		Secure:   false,                   // Mude para True quando colocar em produção em https
		SameSite: http.SameSiteStrictMode, // Protege contra ataques CSRF de sites externos
	})

	// Login feito kcom sucesso! Redireciona para a Home do Dashboard
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout limpa o cookie destruindo a sessão imediatamente
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Define uma data no passado para o navegador apagar ele na hora
		Path:     "/",
		HttpOnly: true,
	})
	tmplLogout.Execute(w, nil)
}

func (h *AuthHandler) ExibirAcessoNegado(w http.ResponseWriter, r *http.Request) {
	tmplAcessoNegado.Execute(w, nil)
}
