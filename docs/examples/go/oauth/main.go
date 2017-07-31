package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

type sessionInformation struct {
	ClientID        string
	Secret          string
	RequestedScopes string
}

// This is only an example, in no way should the following code be used in production!
//   It lacks all of the necessary validation and proper handling

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	//When asked to log in, generate a unique state and store the entered values in a cookie to retreive on callback
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		s := &sessionInformation{
			ClientID:        r.FormValue("client_id"),
			Secret:          r.FormValue("secret"),
			RequestedScopes: r.FormValue("requested_scopes"),
		}
		fmt.Printf("Logging in to %s with secret %s and asking for scopes %s\n", s.ClientID, s.Secret, s.RequestedScopes)

		randombytes := make([]byte, 12) //Multiple of 3 to make sure no padding is added
		rand.Read(randombytes)
		state := base64.URLEncoding.EncodeToString(randombytes)

		serializedSessionInformation, _ := json.Marshal(s)
		sessionCookie := &http.Cookie{
			Name:   state,
			Value:  base64.URLEncoding.EncodeToString(serializedSessionInformation),
			MaxAge: 5 * 60, // 5 minutes
		}
		http.SetCookie(w, sessionCookie)

		u, _ := url.Parse("https://dev.itsyou.online:8443/v1/oauth/authorize")
		q := u.Query()
		q.Add("client_id", s.ClientID)
		q.Add("redirect_uri", "http://localhost:8080/callback")
		q.Add("scope", s.RequestedScopes)
		q.Add("state", state)
		q.Add("response_type", "code")
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		code := r.FormValue("code")
		fmt.Println("Code: ", code)
		state := r.FormValue("state")

		sessionCookie, err := r.Cookie(state)
		if err != nil || sessionCookie == nil {
			fmt.Println("Failed to retrieve session cookie: ", err)
			fmt.Println("Cookie: ", sessionCookie)
			return
		}

		fmt.Println("Cookie value: ", sessionCookie.Value)

		decodedCookie, err := base64.URLEncoding.DecodeString(sessionCookie.Value)
		if err != nil {
			fmt.Println("Could not decode session cookie value: ", err)
		}
		s := &sessionInformation{}

		err = json.Unmarshal(decodedCookie, s)
		if err != nil {
			fmt.Println("Failed to convert session cookie: ", err)
			return
		}

		fmt.Println(s)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc := &http.Client{Transport: tr}
		req, err := http.NewRequest("POST", "https://dev.itsyou.online:8443/v1/oauth/access_token", nil)
		if err != nil {
			return
		}
		q := req.URL.Query()
		q.Add("client_id", s.ClientID)
		q.Add("client_secret", s.Secret)
		q.Add("code", code)
		q.Add("redirect_uri", "http://localhost:8080/callback")
		q.Add("state", state)
		req.URL.RawQuery = q.Encode()

		resp, err := hc.Do(req)
		if err != nil {
			fmt.Println("Error while getting the access token: ", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		fmt.Println("Response body: ", string(body))
		return

	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
