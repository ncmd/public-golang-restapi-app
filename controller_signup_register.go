package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// Add account information in firestore
func controller_signup_register_email_address(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		fmt.Println(err)
	}
	if r.Method != "OPTIONS" {
		client, err := app.Auth(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		var account Account
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&account)
		if err != nil {
			panic(err)
		}
		// Generate password for account
		var password = generatepassword()
		// Send Email with temporary credentials
		sendgrid_email_credentials(account.Username, account.Email, password)
		params := (&auth.UserToCreate{}).Email(account.Email).EmailVerified(false).Password(password).Disabled(false)
		u, err := client.CreateUser(context.Background(), params)
		if err != nil {
			fmt.Println("error creating user:", err)
		}
		log.Printf("Successfully created user: %v\n", u.UID)
	}
}
