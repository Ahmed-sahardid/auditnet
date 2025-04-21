package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ahmedsahardid/auditnet/internal/auth"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 AuditNet is live!")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("web/templates/login.html")
			if err != nil {
				http.Error(w, "Error loading template", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
		} else {
			auth.Authenticate(w, r)
		}
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserRole(r)
		tmpl, err := template.ParseFiles("web/templates/dashboard.html")
		if err != nil {
			http.Error(w, "Error loading dashboard", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, struct {
			Username string
			Role     string
		}{Username: r.FormValue("username"), Role: user})
	})

	http.HandleFunc("/admin", auth.RequireRole("admin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🛠️ Admin Panel - Welcome Admin")
	}))

	// ✅ Move this ABOVE the server start
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "user",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	// ✅ Start server LAST
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

