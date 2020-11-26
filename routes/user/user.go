package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"ssd-coursework/app"
	"ssd-coursework/routes/templates"

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
	req.Header.Add("authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InBJZ1ZMbUM5LXlnUWIzSWtZM0tpaSJ9.eyJpc3MiOiJodHRwczovL2Rldi1vNmxucTZkZy5ldS5hdXRoMC5jb20vIiwic3ViIjoiSWl4NzJLc1JnYWhiNm5FR3NSbjl3TDZHY3kwTjI1cHZAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vZGV2LW82bG5xNmRnLmV1LmF1dGgwLmNvbS9hcGkvdjIvIiwiaWF0IjoxNjA2NDE4NTI1LCJleHAiOjE2MDY1MDQ5MjUsImF6cCI6IklpeDcyS3NSZ2FoYjZuRUdzUm45d0w2R2N5ME4yNXB2Iiwic2NvcGUiOiJyZWFkOmNsaWVudF9ncmFudHMgY3JlYXRlOmNsaWVudF9ncmFudHMgZGVsZXRlOmNsaWVudF9ncmFudHMgdXBkYXRlOmNsaWVudF9ncmFudHMgcmVhZDp1c2VycyB1cGRhdGU6dXNlcnMgZGVsZXRlOnVzZXJzIGNyZWF0ZTp1c2VycyByZWFkOnVzZXJzX2FwcF9tZXRhZGF0YSB1cGRhdGU6dXNlcnNfYXBwX21ldGFkYXRhIGRlbGV0ZTp1c2Vyc19hcHBfbWV0YWRhdGEgY3JlYXRlOnVzZXJzX2FwcF9tZXRhZGF0YSByZWFkOnVzZXJfY3VzdG9tX2Jsb2NrcyBjcmVhdGU6dXNlcl9jdXN0b21fYmxvY2tzIGRlbGV0ZTp1c2VyX2N1c3RvbV9ibG9ja3MgY3JlYXRlOnVzZXJfdGlja2V0cyByZWFkOmNsaWVudHMgdXBkYXRlOmNsaWVudHMgZGVsZXRlOmNsaWVudHMgY3JlYXRlOmNsaWVudHMgcmVhZDpjbGllbnRfa2V5cyB1cGRhdGU6Y2xpZW50X2tleXMgZGVsZXRlOmNsaWVudF9rZXlzIGNyZWF0ZTpjbGllbnRfa2V5cyByZWFkOmNvbm5lY3Rpb25zIHVwZGF0ZTpjb25uZWN0aW9ucyBkZWxldGU6Y29ubmVjdGlvbnMgY3JlYXRlOmNvbm5lY3Rpb25zIHJlYWQ6cmVzb3VyY2Vfc2VydmVycyB1cGRhdGU6cmVzb3VyY2Vfc2VydmVycyBkZWxldGU6cmVzb3VyY2Vfc2VydmVycyBjcmVhdGU6cmVzb3VyY2Vfc2VydmVycyByZWFkOmRldmljZV9jcmVkZW50aWFscyB1cGRhdGU6ZGV2aWNlX2NyZWRlbnRpYWxzIGRlbGV0ZTpkZXZpY2VfY3JlZGVudGlhbHMgY3JlYXRlOmRldmljZV9jcmVkZW50aWFscyByZWFkOnJ1bGVzIHVwZGF0ZTpydWxlcyBkZWxldGU6cnVsZXMgY3JlYXRlOnJ1bGVzIHJlYWQ6cnVsZXNfY29uZmlncyB1cGRhdGU6cnVsZXNfY29uZmlncyBkZWxldGU6cnVsZXNfY29uZmlncyByZWFkOmhvb2tzIHVwZGF0ZTpob29rcyBkZWxldGU6aG9va3MgY3JlYXRlOmhvb2tzIHJlYWQ6YWN0aW9ucyB1cGRhdGU6YWN0aW9ucyBkZWxldGU6YWN0aW9ucyBjcmVhdGU6YWN0aW9ucyByZWFkOmVtYWlsX3Byb3ZpZGVyIHVwZGF0ZTplbWFpbF9wcm92aWRlciBkZWxldGU6ZW1haWxfcHJvdmlkZXIgY3JlYXRlOmVtYWlsX3Byb3ZpZGVyIGJsYWNrbGlzdDp0b2tlbnMgcmVhZDpzdGF0cyByZWFkOnRlbmFudF9zZXR0aW5ncyB1cGRhdGU6dGVuYW50X3NldHRpbmdzIHJlYWQ6bG9ncyByZWFkOmxvZ3NfdXNlcnMgcmVhZDpzaGllbGRzIGNyZWF0ZTpzaGllbGRzIHVwZGF0ZTpzaGllbGRzIGRlbGV0ZTpzaGllbGRzIHJlYWQ6YW5vbWFseV9ibG9ja3MgZGVsZXRlOmFub21hbHlfYmxvY2tzIHVwZGF0ZTp0cmlnZ2VycyByZWFkOnRyaWdnZXJzIHJlYWQ6Z3JhbnRzIGRlbGV0ZTpncmFudHMgcmVhZDpndWFyZGlhbl9mYWN0b3JzIHVwZGF0ZTpndWFyZGlhbl9mYWN0b3JzIHJlYWQ6Z3VhcmRpYW5fZW5yb2xsbWVudHMgZGVsZXRlOmd1YXJkaWFuX2Vucm9sbG1lbnRzIGNyZWF0ZTpndWFyZGlhbl9lbnJvbGxtZW50X3RpY2tldHMgcmVhZDp1c2VyX2lkcF90b2tlbnMgY3JlYXRlOnBhc3N3b3Jkc19jaGVja2luZ19qb2IgZGVsZXRlOnBhc3N3b3Jkc19jaGVja2luZ19qb2IgcmVhZDpjdXN0b21fZG9tYWlucyBkZWxldGU6Y3VzdG9tX2RvbWFpbnMgY3JlYXRlOmN1c3RvbV9kb21haW5zIHVwZGF0ZTpjdXN0b21fZG9tYWlucyByZWFkOmVtYWlsX3RlbXBsYXRlcyBjcmVhdGU6ZW1haWxfdGVtcGxhdGVzIHVwZGF0ZTplbWFpbF90ZW1wbGF0ZXMgcmVhZDptZmFfcG9saWNpZXMgdXBkYXRlOm1mYV9wb2xpY2llcyByZWFkOnJvbGVzIGNyZWF0ZTpyb2xlcyBkZWxldGU6cm9sZXMgdXBkYXRlOnJvbGVzIHJlYWQ6cHJvbXB0cyB1cGRhdGU6cHJvbXB0cyByZWFkOmJyYW5kaW5nIHVwZGF0ZTpicmFuZGluZyBkZWxldGU6YnJhbmRpbmcgcmVhZDpsb2dfc3RyZWFtcyBjcmVhdGU6bG9nX3N0cmVhbXMgZGVsZXRlOmxvZ19zdHJlYW1zIHVwZGF0ZTpsb2dfc3RyZWFtcyBjcmVhdGU6c2lnbmluZ19rZXlzIHJlYWQ6c2lnbmluZ19rZXlzIHVwZGF0ZTpzaWduaW5nX2tleXMgcmVhZDpsaW1pdHMgdXBkYXRlOmxpbWl0cyBjcmVhdGU6cm9sZV9tZW1iZXJzIHJlYWQ6cm9sZV9tZW1iZXJzIGRlbGV0ZTpyb2xlX21lbWJlcnMiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.fNQCS0TIEcCXmR7nFVPnNejA3oxXKwmUqPimbkPB_15KfarwgxmahSYWIcrQOxIQ19GqDwKSol7cNtTq17UqzZCj4PoKiutL8x_VM-86SOdMFs5ebJ2yuMZQmT3jalfJNSgCG98sV1Jrmp1zBsi-V8Yli3a9c2Mdw2oO5x31c89yaRzuawKnOGBBg5im7WviuNWLBl-dEtYq3ONlBbdQO-6mFgvnEHm2EY6bblYef39qs5j2Efre6SYgeVM2ebxAo6qPL7Uljg37j9xUIuuT2Xevqu6KKNPlFFNrlF5O6Lbm2tiWzholqNCT_b90dBZ3aWUoOzhLDaMukW3BfDdedw")

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
