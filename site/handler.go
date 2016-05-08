package site

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/monory/messenger/auth"
	"github.com/monory/messenger/inputs/sanitize"
	"github.com/monory/messenger/inputs/validate"
)

var root = "html"

func Listen(db *sql.DB) {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/auth", dbHandler(authHandler, db))
}

func dbHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, root+"/index.html")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	if stat, err := os.Stat(root + r.URL.Path); os.IsNotExist(err) || stat.IsDir() {
		notFoundError(w, err)
		return
	}
	http.ServeFile(w, r, root+r.URL.Path)
}

func authHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		internalServerError(w, err)
		return
	}

	var username, password string

	if len(r.PostForm) != 0 {
		username = r.PostFormValue("username")
		if !validate.Username(username) {
			renderAuth(w, errors.New("Bad username"))
			return
		}
		username = sanitize.Username(username)

		password = r.PostFormValue("password")
		if !validate.Password(password) {
			renderAuth(w, errors.New("Bad password"))
			return
		}
	}

	switch {
	case r.PostFormValue("loginbutton") != "":
		var token auth.UserToken
		token, err = auth.Login(db, username, password)
		if err != nil {
			if _, ok := err.(auth.AuthError); ok {
				renderAuth(w, errors.New("Bad login"))
				return
			}
			internalServerError(w, err)
			return
		}

		var c http.Cookie
		c.Name = "token"
		c.Value = token.String()
		c.MaxAge = 86400 * 7 // a week
		c.Secure = true
		c.HttpOnly = true

		http.SetCookie(w, &c)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	case r.PostFormValue("registerbutton") != "":
		err = auth.Register(db, username, password)
		if err != nil {
			if _, ok := err.(auth.AuthError); ok {
				renderAuth(w, errors.New("Bad register"))
				return
			}
			internalServerError(w, err)
			return
		}
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	renderAuth(w, nil)
}

func renderAuth(w http.ResponseWriter, err error) {
	data := struct {
		Err error
	}{
		err,
	}

	renderTemplate(w, "auth", data)
}
