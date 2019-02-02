package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var sendgridkey = ""

func main() {

	// Determines which stripe and sendgrid key to use from 'APP_ENV'
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if variable[0] == "APP_ENV" {
			if variable[1] == "local " {
				sendgridkey = config.SendgridLocalKey
			} else if variable[1] == "local" {
				sendgridkey = config.SendgridLocalKey
			} else if variable[1] == "prod " {
				sendgridkey = config.SendgridProdKey
			} else if variable[1] == "prod" {
				sendgridkey = config.SendgridProdKey
			} else {
				fmt.Println("NO APP_ENV found!")
			}
		}
	}

	// hello world
	http.HandleFunc("/", helloworld)

	// response pong
	http.HandleFunc("/api/ping", pingpong)

	// SIGNUP - REGISTER
	http.HandleFunc("/api/signup/register", controller_signup_register_email_address)
	// http.HandleFunc("/api/signup/survey", controller_signup_firebase_authentication)
	// ACCOUNT - LOGIN
	// http.HandleFunc("/api/login", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - VERIFICATION
	// http.HandleFunc("/api/account/verification/personal-id", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/approval/stage1", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/approval/stage2", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/approval/stage3", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - INFORMATION
	// http.HandleFunc("/api/account/information", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - FREEZE
	// http.HandleFunc("/api/account/freeze", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/unfreeze", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - REPORT
	// http.HandleFunc("/api/account/report/stripe", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/report/paypal", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/report/apple-appstore", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/report/google-playstore", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - HISTORY
	// http.HandleFunc("/api/account/history/all", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/history/orders", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/history/login", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - DELETE
	// http.HandleFunc("/api/account/delete", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/delete/review", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/delete/verify", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/delete/finalize", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - MULTIFACTOR - PHONE
	// http.HandleFunc("/api/account/phone/check", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/phone/change", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/phone/verify", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - MULTIFACTOR - QUESTIONS
	// http.HandleFunc("/api/account/questions/change", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/questions/verify", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - MULTIFACTOR - DEVICES
	// http.HandleFunc("/api/account/device/check", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/device/change", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/device/verify", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - USERNAME
	// http.HandleFunc("/api/account/username/check", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/username/change", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - EMAIL
	// http.HandleFunc("/api/account/email/check", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/email/change", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/email/verify", controllerFirebaseAuthenticationAddUser)
	// // ACCOUNT - MULTIFACTOR - PASSWORD
	// http.HandleFunc("/api/account/password/check", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/password/reset", controllerFirebaseAuthenticationAddUser)
	// http.HandleFunc("/api/account/password/change", controllerFirebaseAuthenticationAddUser)

	// ACCOUNT - PORTFOLIO
	// /api/account/portfolio/cash
	// /api/account/portfolio/equity
	// /api/account/portfolio/stocks/list

	// ACCOUNT - BANKING - CREDITCARD
	// /api/account/creditcard/link
	// /api/account/creditcard/remove

	// ACCOUNT - BANKING - BANK
	// /api/account/banking/verify
	// /api/account/banking/link
	// /api/account/banking/remove
	// /api/account/banking/withdraw
	// /api/account/banking/withdraw/review
	// /api/account/banking/deposit
	// /api/account/banking/deposit/review
	// /api/account/banking/history
	// /api/account/banking/lock

	// ACCOUNT - SUBSCRIPTION
	// /api/account/subscription/start
	// /api/account/subscription/end

	// ACCOUNT - EQUITY
	// /api/account/portfolio
	// /api/account/order/history

	// ACCOUNT - REFERRAL
	// /api/account/referral
	// /api/account/referral/apply

	// DECISION
	// /api/decision/add
	// /api/decision/result
	// /api/decision/delete
	// /api/decision/archive
	// /api/decision/modify
	// /api/decision/start
	// /api/decision/end
	// /api/decision/metrics

	// DECISION - COMMENT
	// /api/decision/comment/add
	// /api/decision/comment/delete
	// /api/decision/comment/modify
	// /api/decision/comment/reply
	// /api/decision/comment/vote
	// /api/decision/vote

	// STOCKS - LIST
	// /api/stocks/list
	// /api/stocks/list/detailed
	// /api/stocks/list/metrics
	// /api/stocks/rate
	// /api/stocks/watch
	// /api/stocks/equity

	// STOCKS - DETAILS
	// /api/stocks/details/ratings
	// /api/stocks/details/demand
	// /api/stocks/details/supply
	// /api/stocks/details/metrics
	// /api/stocks/details/finances
	// /api/stocks/details/social
	// /api/stocks/details/schedule

	// STOCKS - CHAT
	// /api/stocks/chat

	// STOCKS - EXCHANGE
	// /api/stocks/exchange
	// /api/stocks/exchange/ask
	// /api/stocks/exchange/accept

	// STOCKS - ORDER
	// /api/stocks/order/buy
	// /api/stocks/order/sell
	// /api/stocks/order/review
	// /api/stocks/order/buyback
	// /api/stocks/order/history

	// GDPR - https://www.mobiloud.com/blog/gdpr-compliant-mobile-app/
	// /api/gdpr/consent/opt-in
	// /api/gdpr/consent/opt-out
	// /api/gdpr/data/pull
	// /api/gdpr/data/erase
	// /api/gdpr/github
	// /api/gdpr/google-firebase
	// /api/gdpr/stripe
	// /api/gdpr/paypal
	// /api/gdpr/travis
	// /api/gdpr/google-playstore
	// /api/gdpr/apple-appstore
	// /api/gdpr/paypal
	// /api/gdpr/heroku
	// /api/gdpr/privacy-policy
	// /api/gdpr/legal-terms
	// /api/gdpr/subject-access-request
	// /api/gdpr/data-breach-notification
	// /api/gdpr/data-encryption
	// /api/gdpr/data-storage
	// /api/gdpr/data-collection-justification

	// listen on port, http handler
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Setting a Default port to 8000 to be used locally
	}
	fmt.Println(http.ListenAndServe(":"+port, nil))
	http.ListenAndServe(":"+port, nil)
}

// Setting up Credentials from config.toml file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("./config/config.toml", &c); err != nil {
		fmt.Println(err)
	}
}

var config = Config{}

func init() {
	config.Read()
}
