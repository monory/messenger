package web

import (
	"io"
	"net/http"
	"net/url"

	"encoding/json"
	"log"

	"fmt"
	"io/ioutil"
)

func verifyCaptcha(captchaResponse []string) bool {
	if len(captchaResponse) == 0 {
		return false
	}

	response := captchaResponse[0]
	resp, _ := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {""}, "response": {response}})

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var googleResponse map[string]interface{}

	json.Unmarshal(body, &googleResponse)
	log.Printf("Captcha request: %v\n", googleResponse)

	return googleResponse["success"].(bool)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	requestString := fmt.Sprintf("%#v\n\n", *r)
	formString := fmt.Sprintf("%#v\n\n", r.Form)

	captchaVerification := verifyCaptcha(r.Form["g-recaptcha-response"])

	captchaVerificationString := fmt.Sprintf("%#v\n\n", captchaVerification)

	io.WriteString(w, requestString)
	io.WriteString(w, formString)
	io.WriteString(w, captchaVerificationString)

}
