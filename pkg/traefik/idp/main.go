// idp is a throwaway OAuth 2.0 authorization server standing in for an external
// provider (Google, Okta, Auth0, ...). It implements just enough of the
// authorization-code flow to be real: /authorize, /token, /userinfo.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	mu     sync.Mutex
	codes  = map[string]string{} // authorization code -> user email
	tokens = map[string]string{} // access token       -> user email
)

func token() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	const user = "alice@example.com"

	mux := http.NewServeMux()

	// /authorize auto-approves and bounces back to the client's redirect_uri
	// with a one-time code. A real provider shows a login and consent screen here.
	mux.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		redirectURI := q.Get("redirect_uri")
		code := token()

		mu.Lock()
		codes[code] = user
		mu.Unlock()

		u, err := url.Parse(redirectURI)
		if err != nil {
			http.Error(w, "bad redirect_uri", http.StatusBadRequest)
			return
		}
		p := u.Query()
		p.Set("code", code)
		p.Set("state", q.Get("state"))
		u.RawQuery = p.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
	})

	// /token swaps the code for an access token (server-to-server).
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		mu.Lock()
		email, ok := codes[r.Form.Get("code")]
		delete(codes, r.Form.Get("code"))
		if !ok {
			mu.Unlock()
			http.Error(w, "invalid_grant", http.StatusBadRequest)
			return
		}
		at := token()
		tokens[at] = email
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": at,
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})

	// /userinfo returns the profile behind a bearer token.
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		at := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		mu.Lock()
		email, ok := tokens[at]
		mu.Unlock()
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"sub": email, "email": email, "name": "Alice",
		})
	})

	fmt.Println("idp listening on :9000")
	http.ListenAndServe(":9000", mux)
}
