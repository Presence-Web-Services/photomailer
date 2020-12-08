/*
* photouploader takes files sent via form and uploads them to a google cloud storage bucket.
* The form is passphrase protected, and allows for setting of photo captions and title of group.
* Upon uploading, the firestore database is also updated with the relevant information.
 */
package photomailer

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/presence-web-services/gmailer/v2"
)

var config gmailer.Config

// default important values
var status = http.StatusOK
var errorMessage = ""
var title = ""
var passphrase = ""
var date = ""
var numPhotos = 0

// init loads environment variables and authenticates the gmailer config
func init() {
	loadEnvVars()
	authenticate()
}

// loadEnvVars loads environment variables from a .env file
func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: Could not load environment variables from .env file.")
	}
	config.ClientID = os.Getenv("CLIENT_ID")
	config.ClientSecret = os.Getenv("CLIENT_SECRET")
	config.AccessToken = os.Getenv("ACCESS_TOKEN")
	config.RefreshToken = os.Getenv("REFRESH_TOKEN")
	config.EmailTo = os.Getenv("EMAIL_TO")
	config.EmailFrom = os.Getenv("EMAIL_FROM")
	config.ReplyTo = os.Getenv("EMAIL_FROM")
	config.Subject = os.Getenv("SUBJECT")
}

// authenticate authenticates a gmailer config
func authenticate() {
	err := config.Authenticate()
	if err != nil {
		log.Fatal("Error: Could not authenticate with GMail OAuth using credentials.")
	}
}

// CreateAndRun creates a http server that listens for photo data. You may set the upload beginning offset here.
func CreateAndRun(port string) {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}

// handler checks the title and passphrase for sanity, and then sends email containing title, date, and captions
func handler(response http.ResponseWriter, request *http.Request) {
	defaultValues()
	verifyPost(response, request.Method)
	if status != http.StatusOK {
		http.Error(response, errorMessage, status)
		return
	}
	getFormData(request)
	if status != http.StatusOK {
		http.Error(response, errorMessage, status)
		return
	}
	checkTitle()
	if status != http.StatusOK {
		http.Error(response, errorMessage, status)
		return
	}
	checkPassphrase()
	if status != http.StatusOK {
		http.Error(response, errorMessage, status)
		return
	}
	composeEmail(request)
	if status != http.StatusOK {
		http.Error(response, errorMessage, status)
		return
	}
	sendEmail()
	if status != http.StatusOK {
		http.Error(response, errorMessage, status)
		return
	}
	response.Write([]byte("Email sent successfully!"))
}

// defaultValues sets the status, errorMessage, title, passphrase, date, numphotos, and body to default values
func defaultValues() {
	status = http.StatusOK
	errorMessage = ""
	title = ""
	passphrase = ""
	date = ""
	numPhotos = 0
	config.Body = ""
}

// verifyPost ensures that a POST is sent
func verifyPost(response http.ResponseWriter, method string) {
	if method != "POST" {
		response.Header().Set("Allow", "POST")
		status = http.StatusMethodNotAllowed
		errorMessage = "Error: Method " + method + " not allowed. Only POST allowed."
	}
}

// getFormData populates config struct and hp variable with POSTed data from form submission
func getFormData(request *http.Request) {
	title = request.PostFormValue("title")
	passphrase = request.PostFormValue("passphrase")
	date = request.PostFormValue("date")
	var err error
	numPhotos, err = strconv.Atoi(request.PostFormValue("numPhotos"))
	if err != nil {
		status = http.StatusBadRequest
		errorMessage = "Error: Could not determine number of photos uploaded."
	}
}

// checkTitle ensures the title is filled out
func checkTitle() {
	if title == "" {
		status = http.StatusBadRequest
		errorMessage = "Error: Title not defined."
	}
}

// checkPassphrase checks that the passphrase is correct
func checkPassphrase() {
	if passphrase != "care for your surroundings" {
		status = http.StatusUnauthorized
		errorMessage = "Error: Passphrase incorrect."
	}
}

// composeEmail creates an email with caption, date, and title information included in the body
func composeEmail(request *http.Request) {
	captions := make([]string, numPhotos)
	for i := 0; i < numPhotos; i++ {
		captionName := "caption" + strconv.Itoa(i)
		caption := request.PostFormValue(captionName)
		captions[i] = caption
	}
	config.Body += "Date: " + date + "\n"
	config.Body += "Title: " + title + "\n"
	for i := 0; i < numPhotos; i++ {
		config.Body += "photo" + strconv.Itoa(i) + ": " + captions[i] + "\n"
	}
}

// sendEmail sends an email given a gmailer config
func sendEmail() {
	err := config.Send()
	if err != nil {
		status = http.StatusInternalServerError
		errorMessage = "Error: Internal server error."
		return
	}
}
