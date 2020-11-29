package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"ssd-coursework/app"
	"ssd-coursework/routes/templates"
	"strings"

	"github.com/prometheus/common/log"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionValues := session.Values["profile"].(map[string]interface{})
	fmt.Println(sessionValues["http://localhost:3000/roles"])
	templates.RenderTemplate(w, "user", session.Values["profile"])
}

// GetSessionUsername will return the name of the current session as a string
func GetSessionUsername(w http.ResponseWriter, r *http.Request) string {
	sessionUsername := extractUserIDFromSession(w, r)
	// fmt.Println(sessionUsername)
	return sessionUsername.(string)
}

// extractUsernameFromSession will extract the name from Auth0 data and return it as an interface
func extractUserIDFromSession(w http.ResponseWriter, r *http.Request) interface{} {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		log.Error(err)
	}
	// fmt.Println(session.Values["profile"])
	sessionUsername := session.Values["profile"].(map[string]interface{})["sub"]
	return sessionUsername
}

func getUserRole(w http.ResponseWriter, r *http.Request) interface{} {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		log.Error(err)
	}
	sessionValues := session.Values["profile"].(map[string]interface{})["http://localhost:3000/roles"]
	fmt.Println(sessionValues)
	return sessionValues
}

// IsDeveloper checks if the logged in user has the Developer role
func IsDeveloper(w http.ResponseWriter, r *http.Request) bool {
	user := getUserRole(w, r)

	if userHasRole(user, "Developer") {
		return true
	}
	return false
}

// IsClient checks if the logged in user has the Client role
func IsClient(w http.ResponseWriter, r *http.Request) bool {
	user := getUserRole(w, r)

	if userHasRole(user, "Client") {
		return true
	}
	return false
}

// IsTester check if the logged in user has the Tester role
func IsTester(w http.ResponseWriter, r *http.Request) bool {
	user := getUserRole(w, r)

	if userHasRole(user, "Tester") {
		return true
	}
	return false
}

// GetUsersNameFromAuth returns a name associated with a userID
func GetUsersNameFromAuth(w http.ResponseWriter, r *http.Request, userID string) string {
	url := "https://dev-o6lnq6dg.eu.auth0.com/api/v2/users/" + userID
	req, _ := http.NewRequest("GET", url, nil)
	// This is using a test API token, this should be update to pull in a new auto refreshed one

	accessToken := getAuthMgmtAPIToken()

	req.Header.Add("authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	// Decode byte data to string in json format
	var data map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}
	name := fmt.Sprint(data["name"])

	defer res.Body.Close()

	return name
}

// GetUserIDFromName returns the ID associated with a name
func GetUserIDFromName(w http.ResponseWriter, r *http.Request, username string) string {
	accessToken := getAuthMgmtAPIToken()

	url := "https://dev-o6lnq6dg.eu.auth0.com/api/v2/users?q=" + username
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var data []map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}
	userID := fmt.Sprint(data[0]["user_id"])
	return userID
}

// getAuthMgmtApiToken will get the Auth0 Management API token to be used in calls to the management API
// Auth0 recycles tokens every 30000 seconds so need to get a new one each time
func getAuthMgmtAPIToken() string {
	url := "https://dev-o6lnq6dg.eu.auth0.com/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"Iix72KsRgahb6nEGsRn9wL6Gcy0N25pv\",\"client_secret\":\"C0-2oxWnzskRJU--7n2Q9OOFS6Sr6ocuzjfzgzQ2LdBDTU2f1frOxC4ZGXpM594V\",\"audience\":\"https://dev-o6lnq6dg.eu.auth0.com/api/v2/\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// Decode byte data to string in json format
	var data map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(data["access_token"])
}

// Helper function to check the values inside the user interface
// Interface is a slice, but first need to reflect this type so Go
// knows how to handle it.
func userHasRole(user interface{}, role string) bool {
	switch reflect.TypeOf(user).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(user)

		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == role {
				return true
			}
		}
	}
	return false
}
