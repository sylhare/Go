// Command smoke runs an end-to-end check of the Traefik forward-auth gateway.
// It walks the full OAuth authorization-code flow, then verifies identity
// injection, custom error pages, and per-host routing. It exits non-zero if any
// check fails, so it doubles as a CI gate once the stack is up.
//
//	go run ./cmd/smoke                              # http://127.0.0.1:8000 (compose)
//	BASE=http://127.0.0.1:8001 go run ./cmd/smoke   # kind
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	base = env("BASE", "http://127.0.0.1:8000")
	// Do not follow redirects: the whole point is to inspect each 302 by hand.
	client = &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	pass, fail int
)

func main() {
	fmt.Printf("smoke-testing gateway at %s\n", base)

	// 1. Unauthenticated request is bounced to the provider, target URL in `state`.
	r := get("app.localhost", "/", "")
	checkStatus("unauth request is redirected", http.StatusFound, r.status)
	loc := r.header.Get("Location")
	check("redirect targets the provider", "id.localhost/authorize", loc)
	check("target URL parked in state", "state=http%3A%2F%2Fapp.localhost%2F", loc)

	// 2. The provider issues an authorization code back to the gateway callback.
	r = get("id.localhost", "/authorize?"+query(loc), "")
	cb := r.header.Get("Location")
	check("provider returns an authorization code", "code=", cb)

	// 3. The callback exchanges the code, sets a session cookie, redirects to `state`.
	r = get("auth.localhost", "/callback?"+query(cb), "")
	sid := cookieValue(r, "gw_session")
	checkTrue("callback sets a session cookie", sid != "")
	check("callback returns to the original app", "http://app.localhost/", r.header.Get("Location"))

	// 4. Success path: with the cookie, the request lands on the hello-world page
	//    and the backend has received the authenticated identity.
	jar := "gw_session=" + sid
	r = get("app.localhost", "/", jar)
	checkStatus("authenticated request is allowed", http.StatusOK, r.status)
	check("success path serves the hello page", "Hello from app-a", r.body)
	check("backend sees X-Auth-User", "X-Auth-User: alice@example.com", r.body)

	// 5. Error path: a 5xx from the backend is swapped for the custom error page,
	//    keeping the original status code.
	r = get("app.localhost", "/status/503", jar)
	checkStatus("original 503 preserved", http.StatusServiceUnavailable, r.status)
	check("error path serves the custom page", "<h1>503</h1>", r.body)

	// 6. A different host routes to a different backend on the same session.
	r = get("admin.localhost", "/", jar)
	check("admin host serves its own hello page", "Hello from admin-b", r.body)

	// 7. Without a cookie, the second host also bounces to auth, per its own URL.
	r = get("admin.localhost", "/", "")
	check("admin unauth state is per-URL", "state=http%3A%2F%2Fadmin.localhost%2F", r.header.Get("Location"))

	fmt.Printf("\npassed: %d  failed: %d\n", pass, fail)
	if fail > 0 {
		os.Exit(1)
	}
}

type resp struct {
	status int
	header http.Header
	body   string
}

// get issues a GET to base+path, overriding the Host header so a single
// listener can serve every virtual host in the demo.
func get(host, path, cookie string) resp {
	req, err := http.NewRequest(http.MethodGet, base+path, nil)
	if err != nil {
		die("build request: %v", err)
	}
	req.Host = host
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	r, err := client.Do(req)
	if err != nil {
		die("cannot reach %s (is the stack up?): %v", base, err)
	}
	defer r.Body.Close()
	b, _ := io.ReadAll(r.Body)
	return resp{status: r.StatusCode, header: r.Header, body: string(b)}
}

func check(name, want, got string) {
	if strings.Contains(got, want) {
		ok(name)
		return
	}
	no(name, fmt.Sprintf("want substring %q, got %q", want, trunc(got)))
}

func checkStatus(name string, want, got int) {
	if want == got {
		ok(name)
		return
	}
	no(name, fmt.Sprintf("want status %d, got %d", want, got))
}

func checkTrue(name string, cond bool) {
	if cond {
		ok(name)
		return
	}
	no(name, "condition was false")
}

func ok(name string) { fmt.Printf("  ok   %s\n", name); pass++ }
func no(name, why string) {
	fmt.Printf("  FAIL %s\n       %s\n", name, why)
	fail++
}

func cookieValue(r resp, name string) string {
	for _, c := range (&http.Response{Header: r.header}).Cookies() {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func query(location string) string {
	u, err := url.Parse(location)
	if err != nil {
		die("parse Location %q: %v", location, err)
	}
	return u.RawQuery
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func trunc(s string) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(s) > 160 {
		return s[:160] + "…"
	}
	return s
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}
