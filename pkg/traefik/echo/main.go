// Command echo is a backend that greets the authenticated user and echoes the request headers.
package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	name := os.Getenv("NAME")
	if name == "" {
		name = "echo"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		code, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/status/"))
		if err != nil || code < 100 || code > 599 {
			code = 400
		}
		w.WriteHeader(code)
		fmt.Fprintf(w, "%s returning %d\n", name, code)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		who := r.Header.Get("X-Auth-User")
		if who == "" {
			who = "an anonymous visitor"
		}

		var dump strings.Builder
		fmt.Fprintf(&dump, "%s %s%s\n", r.Method, r.Host, r.URL.Path)
		keys := make([]string, 0, len(r.Header))
		for k := range r.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&dump, "%s: %s\n", k, strings.Join(r.Header[k], ","))
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>Hello from %[1]s</title></head>
<body style="font-family:system-ui,sans-serif;padding:2rem">
  <h1>Hello from %[1]s</h1>
  <p>You reached this app as <strong>%[2]s</strong>.</p>
  <pre>%[3]s</pre>
</body>
</html>
`, name, who, dump.String())
	})

	fmt.Printf("%s listening on :8080\n", name)
	http.ListenAndServe(":8080", mux)
}
