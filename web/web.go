package web

import (
	"fmt"
	"net/http"
	"net/url"

	"encoding/json"
	"log"

	"io/ioutil"

	"database/sql"
	"github.com/monory/messager-backend/database"
)

func verifyCaptcha(response string) bool {
	if len(response) == 0 {
		return false
	}

	resp, _ := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {""}, "response": {response}})

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var googleResponse map[string]interface{}

	json.Unmarshal(body, &googleResponse)
	log.Printf("Captcha request: %v\n", googleResponse)

	return googleResponse["success"].(bool)
}

func Handler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.ParseForm()
	captchaVerification := verifyCaptcha(r.Form["g-recaptcha-response"][0])

	if !captchaVerification {
		fmt.Fprint(w, "Bad captcha")
		return
	}

	if len(r.Form["loginbutton"]) > 0 { // pressed "Log In"
		fmt.Fprint(w, database.CheckUser(db, r.Form["login"][0], r.Form["password"][0]))
	} else { // pressed "Register"
		fmt.Fprint(w, database.AddUser(db, r.Form["login"][0], r.Form["password"][0]))
	}

}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}
