package logout

import (
	"net/http"
	"net/url"
	"ssd-coursework/routes/middlewares"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	domain := "dev-o6lnq6dg.eu.auth0.com"

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host + "/login")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", "UmOdRUfJnUUgNB5zHhlltXjTLqQJp5GM")
	logoutUrl.RawQuery = parameters.Encode()
	middlewares.ClearSession(w, r)
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
