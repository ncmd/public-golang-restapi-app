package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// recaptcha confirmation
func controllers_recaptcha_verify_recaptcha_requests_from_signup_page(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	var data Data
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)

	if (data.Recaptcha == "undefined") || (data.Recaptcha == "") {
		log.Println("Nothing in data.Recaptcha", data.Recaptcha)
	} else {
		log.Println("RECAPTCHA:", data.Recaptcha)

		var secretKey = config.Captchasecret

		url := "https://www.google.com/recaptcha/api/siteverify"

		payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r" +
			"\nContent-Disposition: form-data; " +
			"name=\"secret\"\r\n\r\n" + secretKey + "\r" +
			"\n------WebKitFormBoundary7MA4YWxkTrZu0gW\r" +
			"\nContent-Disposition: form-data; name=\"response\"\r\n" +
			"\r\n" + data.Recaptcha + "\r\n" +
			"------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"remoteip\"\r\n" +
			"\r\n" + r.RemoteAddr + "\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")

		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Cache-Control", "no-cache")

		response := &Response{
			Success:     false,
			ChallengeTS: time.Now().String(),
			Hostname:    "Localhost",
		}

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		log.Println("RESPONSE BODY:", string(body))
		json.Unmarshal(body, &response)
		log.Println("RESPONSE RECAPTCHA:", response.Success)
	}
}
