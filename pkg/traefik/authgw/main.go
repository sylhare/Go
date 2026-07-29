// Command authgw is the forward-auth gateway and OAuth client Traefik calls before every protected request.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var (
	mu       sync.Mutex
	sessions = map[string]string{}
)

func sid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	var (
		authorizeURL = env("AUTHORIZE_URL", "http://id.localhost/authorize")
		tokenURL     = env("TOKEN_URL", "http://idp:9000/token")
		userinfoURL  = env("USERINFO_URL", "http://idp:9000/userinfo")
		clientID     = env("CLIENT_ID", "gateway")
		clientSecret = env("CLIENT_SECRET", "s3cret")
		callbackURL  = env("CALLBACK_URL", "http://auth.localhost/callback")
		cookieDomain = env("COOKIE_DOMAIN", "localhost")
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("gw_session"); err == nil {
			mu.Lock()
			email, ok := sessions[c.Value]
			mu.Unlock()
			if ok {
				w.Header().Set("X-Auth-User", email)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		target := fmt.Sprintf("%s://%s%s",
			r.Header.Get("X-Forwarded-Proto"),
			r.Header.Get("X-Forwarded-Host"),
			r.Header.Get("X-Forwarded-Uri"))

		u, _ := url.Parse(authorizeURL)
		q := u.Query()
		q.Set("response_type", "code")
		q.Set("client_id", clientID)
		q.Set("redirect_uri", callbackURL)
		q.Set("scope", "openid email")
		q.Set("state", target)
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		target := r.URL.Query().Get("state")

		resp, err := http.PostForm(tokenURL, url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {callbackURL},
			"client_id":     {clientID},
			"client_secret": {clientSecret},
		})
		if err != nil {
			http.Error(w, "token exchange failed", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		var tok struct {
			AccessToken string `json:"access_token"`
		}
		json.NewDecoder(resp.Body).Decode(&tok)

		req, _ := http.NewRequest(http.MethodGet, userinfoURL, nil)
		req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
		ui, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "userinfo failed", http.StatusBadGateway)
			return
		}
		defer ui.Body.Close()
		var info struct {
			Email string `json:"email"`
		}
		json.NewDecoder(ui.Body).Decode(&info)

		id := sid()
		mu.Lock()
		sessions[id] = info.Email
		mu.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     "gw_session",
			Value:    id,
			Path:     "/",
			Domain:   cookieDomain,
			HttpOnly: true,
		})
		if target == "" {
			target = "/"
		}
		http.Redirect(w, r, target, http.StatusFound)
	})

	fmt.Println("authgw listening on :8081")
	http.ListenAndServe(":8081", mux)
}
