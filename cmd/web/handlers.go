package main

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "home.page.gohtml", nil)
}
