package auth

import (
	"net/http"
)

var users = map[string]struct {
	Password string
	Role     string
}{
	"ahmed":  {Password: "admin123", Role: "admin"},
	"muna":   {Password: "pass123", Role: "manager"},
	"bilal":  {Password: "employee", Role: "employee"},
	"guest":  {Password: "guest", Role: "guest"},
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]
	if !ok || user.Password != password {
		http.Redirect(w, r, "/login?error=invalid", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: username,
		Path:  "/",
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func GetUserRole(r *http.Request) string {
	cookie, err := r.Cookie("user")
	if err != nil {
		return "guest"
	}
	user := users[cookie.Value]
	return user.Role
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if GetUserRole(r) != role {
			http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

