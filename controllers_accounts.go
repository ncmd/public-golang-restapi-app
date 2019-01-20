package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func c_accounts_add_activity_objectives_to_user_in_firestore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		fmt.Println(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	if r.Method != "OPTIONS" {

		var account Account
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&account)
		if err != nil {
			panic(err)
		}

		fmt.Println("Updating Activity!")
		_, err = client.Collection("accounts").Doc(account.ID).Set(context.Background(), map[string]interface{}{
			"activity": []interface{}{
				map[string]interface{}{
					"runbookid":     account.Activity.Runbookid,
					"runbooktitle":  account.Activity.Runbooktitle,
					"runbookstatus": account.Activity.Runbookstatus,
					"runbookobjectives": []interface{}{
						map[string]interface{}{
							"objectivetitle": account.Activity.Runbookobjectives.Objectivetitle,
						},
					},
				},
			},
		}, firestore.MergeAll)
		if err != nil {
			fmt.Println(err)
		}

	}
}

func c_accounts_invite_user_create_account_in_firebase(w http.ResponseWriter, r *http.Request) {
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

		var password = generatepassword()
		sendgridEmail("Cavalry_User", account.Email, password)

		params := (&auth.UserToCreate{}).Email(account.Email).EmailVerified(false).Password(password).Disabled(false)
		u, err := client.CreateUser(context.Background(), params)
		if err != nil {
			fmt.Println("error creating user:", err)
		}
		log.Printf("Successfully created user: %v\n", u.UID)
		c_accounts_create_account_in_firebase(u.UID, account.Email, account.Username, account.OrganizationName)
		var data = u.UID
		js, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)

	}
}

func c_accounts_create_user_account_in_firebase(w http.ResponseWriter, r *http.Request) {
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

		var password = generatepassword()

		// Send Email to this email address
		sendgridEmail(account.Username, account.Email, password)

		params := (&auth.UserToCreate{}).Email(account.Email).EmailVerified(false).Password(password).Disabled(false)
		u, err := client.CreateUser(context.Background(), params)
		if err != nil {
			fmt.Println("error creating user:", err)
		}
		log.Printf("Successfully created user: %v\n", u.UID)

		c_accounts_create_account_in_firebase(u.UID, account.Email, account.Username, "")

	}
}

func c_accounts_create_account_in_firebase(userid string, useremail string, username string, organizationname string) {
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		fmt.Println(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	if organizationname != "" {
		client.Collection("accounts").Doc(userid).Set(context.Background(), map[string]interface{}{
			"accountid":          userid,
			"emailaddress":       useremail,
			"username":           username,
			"organizationname":   organizationname,
			"organizationmember": true,
			"activity":           []interface{}{},
		})
	} else {
		client.Collection("accounts").Doc(userid).Set(context.Background(), map[string]interface{}{
			"accountid":          userid,
			"emailaddress":       useremail,
			"username":           username,
			"organizationname":   "",
			"organizationmember": false,
			"activity":           []interface{}{},
		})
	}

}

// Get account information in firestore (id, email, plan)
func c_accounts_get_user_account_information_in_firestore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		fmt.Println(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	var account Account
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&account)
	if err != nil {
		panic(err)
	}

	// Client usually makes 2 requests to backend; 1 options; 1 post, need to filter out options request
	if r.Method != "OPTIONS" {
		// fmt.Println("AccountID:", account.ID)
		// Creating in 'accounts' collection, with the ID generated by Firebase Authentication UID
		dsnap, err := client.Collection("accounts").Doc(account.ID).Get(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		m := dsnap.Data()
		// fmt.Printf("Document data: %#v\n", m)

		js, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}
}

func c_accounts_create_account_with_email_in_firestore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		fmt.Println(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	var account Account
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&account)
	if err != nil {
		panic(err)
	}
	if r.Method != "OPTIONS" {
		// Creating in 'accounts' collection, with the ID generated by Firebase Authentication UID
		client.Collection("accounts").Doc(account.ID).Set(context.Background(), map[string]interface{}{
			"accountid":    account.ID,
			"emailaddress": account.Email,
			"username":     account.Username,
		})
	}
}

// Add account information in firestore (id, email, plan)
func controllers_accounts_create_user_account_in_firestore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		fmt.Println(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	var account Account
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&account)
	if err != nil {
		panic(err)
	}

	// Client usually makes 2 requests to backend; 1 options; 1 post, need to filter out options request
	if r.Method != "OPTIONS" {
		// Creating in 'accounts' collection, with the ID generated by Firebase Authentication UID
		client.Collection("accounts").Doc(account.ID).Set(context.Background(), map[string]interface{}{
			"accountid":    account.ID,
			"username":     account.Username,
			"emailaddress": account.Email,
			"plan":         account.Plan,
		})
	}
}
