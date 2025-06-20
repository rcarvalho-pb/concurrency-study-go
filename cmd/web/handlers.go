package main

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.InfoLog.Println("User not found")
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed login attempt",
			Data:    "invalid login attempt",
		}
		app.SendEmail(msg)
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "successful login")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {

}
