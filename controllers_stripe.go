package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sub"
	"google.golang.org/api/option"
)

// get firestore accountid

// get stripe customerid

// subscribe user
func controllers_stripe_generate_password_subscribe_user_to_plan_sendgrid_email_password(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "OPTIONS" {
		var user User

		log.Println("Response Body:", r.Method)
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&user)
		log.Println("User:", user)
		log.Println("Accountid:", user.Accountid)
		log.Println("User Email:", user.Email)
		log.Println("User Source:", user.Source)
		log.Println("User Plan:", user.Plan)
		var password = generatepassword()
		var customerid = createCustomer(user.Email, user.Source)
		w.Write([]byte(password))
		if user.Plan == "1month" {
			log.Println("Selected 1 Month!")
			var subscriptionid = newSubscriber1Month(customerid)
			add_account_subscription_to_account_in_firestore(customerid, user.Email, subscriptionid, user.Plan, user.Accountid)
			// sendEmail(user.Email, password)
			sendgridEmail("user", user.Email, password)
		} else if user.Plan == "12months" {
			log.Println("Selected 12 Months!")
			var subscriptionid = newSubscriber12Months(customerid)
			add_account_subscription_to_account_in_firestore(customerid, user.Email, subscriptionid, user.Plan, user.Accountid)
			// sendEmail(user.Email, password)
			sendgridEmail("user", user.Email, password)
		} else if user.Plan == "beta" {
			log.Println("Selected Beta!")
			var subscriptionid = newSubscriberBeta(customerid)
			add_account_subscription_to_account_in_firestore(customerid, user.Email, subscriptionid, user.Plan, user.Accountid)
			// sendEmail(user.Email, password)
			sendgridEmail("user", user.Email, password)
		} else {
			log.Println("No Plan Selected...")
		}
	}
}

func add_account_subscription_to_account_in_firestore(customerid string, customeremail string, customersubscription string, customerplan string, accountid string) {
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	batch := client.Batch()

	// Set the value
	ref := client.Collection("accounts").Doc(accountid)
	batch.Set(ref, map[string]interface{}{
		"accountid":              accountid,
		"emailaddress":           customeremail,
		"stripeSubscriptionPlan": customerplan,
		"stripeCustomerId":       customerid,
		"stripeSubscriptionId":   customersubscription,
	}, firestore.MergeAll)

	// Commit the batch.
	_, err = batch.Commit(context.Background())
	if err != nil {
		return
	}
	log.Println("Updated Account!")
}

// this teaches you json
func get_stripe_customerId_from_emailaddress(w http.ResponseWriter, r *http.Request) {
	stripe.Key = stripesecretkey

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "OPTIONS" {

		var customer Customer
		log.Println("Response Body:", r.Method)
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&customer)
		url := "https://api.stripe.com/v1/search?query=" + customer.Email + "&prefix=false"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("authorization", "Bearer "+stripe.Key)
		req.Header.Add("content-type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(res.Body)

		// s = byte ; this needs to be converted to buffer
		s, err := getStripeCustomer([]byte(body))

		// create a buffer type
		var buf []byte

		for i, sdata := range s.StripeCustomerData {
			fmt.Println("i:", i, "sdata:", sdata)
		}
		// conver json to buffer using marshal ; getting first element; its customerid
		buf, err = json.Marshal(s.StripeCustomerData[0])
		if err != nil {
			fmt.Println("whoops:", err)
		}

		// w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
		// fmt.Println(s.StripeCustomerData[0], err)
	}

}

// this teaches you json
func get_stripe_subscriptionid_from_customerid(w http.ResponseWriter, r *http.Request) {
	stripe.Key = stripesecretkey

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "OPTIONS" {

		var customer Customer
		log.Println("Response Body:", r.Method)
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&customer)
		url := "https://api.stripe.com/v1/subscriptions"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("authorization", "Bearer "+stripe.Key)
		req.Header.Add("content-type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(res.Body)

		// s = byte ; this needs to be converted to buffer
		// s, _ := getStripeCustomer([]byte(body))
		fmt.Println(string(body))
		// fmt.Println(res.Body)
	}

}

func getStripeSubscription(body []byte) (*StripeSubscriptionResponse, error) {
	var s = new(StripeSubscriptionResponse)

	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return s, err
}

func getStripeCustomer(body []byte) (*StripeCustomerResponse, error) {
	var s = new(StripeCustomerResponse)

	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return s, err
}

func update_stripe_customer_information_to_accounts_firestore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sa := option.WithCredentialsFile("./firestore.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Get a new write batch.
	batch := client.Batch()

	if r.Method != "OPTIONS" {
		var customer Customer
		log.Println("Response Body:", r.Method)
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&customer)

		// Set the value
		ref := client.Collection("accounts").Doc(customer.FirebaseAccountId)
		batch.Set(ref, map[string]interface{}{
			"stripeCustomerId": customer.StripeCustomerId,
		}, firestore.MergeAll)

		// Commit the batch.
		_, err := batch.Commit(context.Background())
		if err != nil {
			return
		}
	}
}

func newSubscriberBeta(customerid string) string {
	messages := make(chan string, 2)
	messages <- customerid
	time.Sleep(time.Second * 3)

	return subscribeCustomer(<-messages, stripebetaplan)
}

func newSubscriber1Month(customerid string) string {
	messages := make(chan string, 2)
	messages <- customerid
	time.Sleep(time.Second * 3)
	return subscribeCustomer(<-messages, stripe1monthplan)
}

func newSubscriber12Months(customerid string) string {
	messages := make(chan string, 2)
	messages <- customerid
	time.Sleep(time.Second * 3)
	return subscribeCustomer(<-messages, stripe12monthsplan)
}

// Creates a Service;
// This should only be created for the first 3 plans (1m, 12m, lifetime)
func createService(serviceName string) string {
	stripe.Key = stripesecretkey

	params := &stripe.ProductParams{
		Name: stripe.String(serviceName),
		Type: stripe.String(string(stripe.ProductTypeService)),
	}
	prod, _ := product.New(params)
	fmt.Println("Product:", prod.ID)
	return prod.ID
}

// Creates a Attach Payment plan for service
func attachPlan(productID string, planNickname string, planAmount int64) string {
	stripe.Key = stripesecretkey

	params := &stripe.PlanParams{
		ProductID: stripe.String(productID),
		Nickname:  stripe.String(planNickname),
		Interval:  stripe.String(string(stripe.PlanIntervalMonth)),
		Currency:  stripe.String("usd"),
		Amount:    stripe.Int64(planAmount),
	}
	// Amount 1000 = $10.00
	p, err := plan.New(params)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("Plan:", p.ID)
	return p.ID
}

func createCustomer(emailaddress string, source string) string {
	stripe.Key = stripesecretkey
	// fmt.Println("Customer email", emailaddress, source)
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(emailaddress),
	}
	customerParams.SetSource(source)
	c, _ := customer.New(customerParams)
	return c.ID

}

func subscribeCustomer(customer string, plan string) string {
	stripe.Key = stripesecretkey

	items := []*stripe.SubscriptionItemsParams{
		{
			Plan: stripe.String(plan),
		},
	}
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customer),
		Items:    items,
	}
	subscription, _ := sub.New(params)

	fmt.Println("Subscription:", subscription.ID)
	return subscription.ID

}

// func getCustomerInfo(customerid string) {
// 	stripe.Key = stripesecretkey
// 	c, err := customer.Get(customerid, nil)
// }

func c_stripe_unsubscribe_account_in_stripe(w http.ResponseWriter, r *http.Request) {
	stripe.Key = stripesecretkey
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	var account Account
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&account)

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	subscription, _ := sub.Update(account.StripeSubscriptionId, params)
	fmt.Println("Unsubcribed:", subscription)
}

func sendInvoice(customerid string, plan string) {
	stripe.Key = stripesecretkey

	items := []*stripe.SubscriptionItemsParams{
		{Plan: stripe.String(plan)},
	}
	params := &stripe.SubscriptionParams{
		Customer:     stripe.String(customerid),
		Items:        items,
		Billing:      stripe.String(string(stripe.InvoiceBillingSendInvoice)),
		DaysUntilDue: stripe.Int64(30),
	}
	subscription, _ := sub.New(params)
	fmt.Println("Subscription:", subscription)
}
