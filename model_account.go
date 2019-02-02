package main

type Account struct {
	ID       string `json:"accountid"`
	Username string `json:"username"`
	Plan     string `json:"plan"`
	Email    string `json:"email"`
}
