// Command idp is a stand-in OAuth 2.0 provider implementing the authorization-code flow.
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
	codes  = map[string]string{}
	tokens = map[string]string{}
)

func token() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	const user = "alice@example.com"

	mux := http.NewServeMux()

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
