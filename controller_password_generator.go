package main

import (
	"io/ioutil"
	"net/http"
)

func generatepassword() string {
	url := "https://www.dinopass.com/password/simple"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
