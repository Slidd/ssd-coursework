package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	return sessionUsername.(string)
}

// GetAuth0UserRole will return the role for a user from the Auth0 User database
func GetAuth0UserRole(w http.ResponseWriter, r *http.Request, name string) string {
	userID := GetUserIDFromName(w, r, name)
	url := "https://dev-o6lnq6dg.eu.auth0.com/api/v2/users/" + userID + "/roles"
	req, _ := http.NewRequest("GET", url, nil)
	// This is using a test API token, this should be update to pull in a new auto refreshed one

	accessToken := getAuthMgmtAPIToken()

	req.Header.Add("authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var data []map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}
	role := fmt.Sprint(data[0]["name"])
	return role
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

func GetUserRole(w http.ResponseWriter, r *http.Request) interface{} {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		log.Error(err)
	}
	sessionValues := session.Values["profile"].(map[string]interface{})["http://localhost:3000/roles"]
	fmt.Println(sessionValues)
	return sessionValues
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

// GetAllUsers Returns all the names of the users
func GetAllUsers(w http.ResponseWriter, r *http.Request) []string {
	accessToken := getAuthMgmtAPIToken()

	url := "https://dev-o6lnq6dg.eu.auth0.com/api/v2/users?q=name:*&include_fields=true&fields=name"
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
	var usernames []string
	for i := range data {
		usernames = append(usernames, fmt.Sprint(data[i]["name"]))
	}
	return usernames

}

type users struct {
	name string
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
	auth0MGMTID := os.Getenv("AUTH0_MGMT_API_CLIENT_ID")
	auth0MGMTSecret := os.Getenv("AUTH0_MGMT_API_CLIENT_SECRET")
	payload := strings.NewReader("{\"client_id\":\"" + auth0MGMTID + "\",\"client_secret\":\"" + auth0MGMTSecret + "\",\"audience\":\"https://dev-o6lnq6dg.eu.auth0.com/api/v2/\",\"grant_type\":\"client_credentials\"}")
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
