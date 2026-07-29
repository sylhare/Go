// errors serves the custom error pages. Traefik's errors middleware forwards the
// failed request here with the path from its `query` template (for example
// /500.html), and returns whatever this service produces to the original client
// while keeping the original status code.
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := strings.TrimSuffix(strings.Trim(r.URL.Path, "/"), ".html")
		if status == "" {
			status = "error"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>%[1]s</title></head>
<body style="font-family:system-ui,sans-serif;text-align:center;padding-top:4rem">
  <h1>%[1]s</h1>
  <p>The gateway caught this one. Nothing leaked from the backend.</p>
</body>
</html>
`, status)
	})

	fmt.Println("errors listening on :8082")
	http.ListenAndServe(":8082", mux)
}
