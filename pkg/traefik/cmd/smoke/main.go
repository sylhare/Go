// Command smoke runs an end-to-end check of the Traefik forward-auth gateway.
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
	client = &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	pass, fail int
)

func main() {
	fmt.Printf("smoke-testing gateway at %s\n", base)

	r := get("app.localhost", "/", "")
	checkStatus("unauth request is redirected", http.StatusFound, r.status)
	loc := r.header.Get("Location")
	check("redirect targets the provider", "id.localhost/authorize", loc)
	check("target URL parked in state", "state=http%3A%2F%2Fapp.localhost%2F", loc)

	r = get("id.localhost", "/authorize?"+query(loc), "")
	cb := r.header.Get("Location")
	check("provider returns an authorization code", "code=", cb)

	r = get("auth.localhost", "/callback?"+query(cb), "")
	sid := cookieValue(r, "gw_session")
	checkTrue("callback sets a session cookie", sid != "")
	check("callback returns to the original app", "http://app.localhost/", r.header.Get("Location"))

	jar := "gw_session=" + sid
	r = get("app.localhost", "/", jar)
	checkStatus("authenticated request is allowed", http.StatusOK, r.status)
	check("success path serves the hello page", "Hello from app-a", r.body)
	check("backend sees Remote-User", "Remote-User: alice@example.com", r.body)

	r = get("app.localhost", "/status/503", jar)
	checkStatus("original 503 preserved", http.StatusServiceUnavailable, r.status)
	check("error path serves the custom page", "<h1>503</h1>", r.body)

	r = get("admin.localhost", "/", jar)
	check("admin host serves its own hello page", "Hello from admin-b", r.body)

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

// get issues a GET to base+path with an overridden Host header.
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
