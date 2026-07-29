# Traefik forward-auth gateway

A runnable example of using [Traefik](https://doc.traefik.io/traefik/) as a single edge router that:

- authenticates every request against an external OAuth 2.0 provider (authorization-code flow),
- sends the user back to the exact URL they first asked for, whichever app that is,
- passes the authenticated identity to the backend as a header,
- serves custom error pages instead of leaking a backend's 5xx response.

The gateway logic lives entirely in Traefik plus one small auth service. The apps behind it stay
oblivious: they receive already-authenticated requests and never see the login.

## What's in the box

Six services, wired through Traefik's **file provider** (no Docker socket needed):

| Service   | Port   | Role                                                              |
|-----------|--------|-------------------------------------------------------------------|
| `traefik` | 80     | Edge router. Web entrypoint on `:80`, dashboard on `:8080`.       |
| `echo-a`  | 8080   | Backend for `app.localhost`. Prints the request it received.      |
| `echo-b`  | 8080   | Backend for `admin.localhost`. Same image, different name.        |
| `authgw`  | 8081   | Forward-auth gateway + OAuth client (`/auth`, `/callback`).       |
| `idp`     | 9000   | Throwaway OAuth 2.0 provider standing in for Google/Okta/Auth0.   |
| `errors`  | 8082   | Serves the custom error pages (`/{status}.html`).                 |

Two Traefik middlewares do the work, both defined in `dynamic.yml`:

- `gateway-auth` (`forwardAuth`) calls `authgw` before every protected request. A `2xx` lets the
  request through; anything else is returned to the client verbatim, which is how a `302` becomes a login.
- `custom-errors` (`errors`) catches `500-599` from a backend and fetches a replacement page from the
  `errors` service, keeping the original status code.

## How a request flows

The first request has no session, so it is bounced through the OAuth authorization-code flow and comes
back with a session cookie. Every later request carries that cookie and is let straight through with the
identity injected as `Remote-User`.

`authgw`'s `/auth` endpoint never sees the request body — only the
[`X-Forwarded-*` headers](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) Traefik copies in
(`X-Forwarded-Method` / `-Proto` / `-Host` / `-Uri`), which it uses to rebuild the exact URL the user
asked for and park it in the OAuth `state`. Those `X-Forwarded-*` names are set by
[Traefik's ForwardAuth middleware](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) and are
kept as-is — the `X-` prefix is deprecated by [RFC 6648](https://www.rfc-editor.org/rfc/rfc6648), but they
are a long-standing de-facto standard we do not control. The one header we *do* set, `Remote-User`, drops
the prefix. On success `/auth` sets `Remote-User` on its response, and Traefik propagates that header to
the backend because it is listed under the middleware's
[`authResponseHeaders`](https://doc.traefik.io/traefik/middlewares/http/forwardauth/#authresponseheaders).

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant T as Traefik
    participant G as authgw
    participant I as idp (OAuth)
    participant B as echo (backend)

    Note over C,B: First request — no session
    C->>T: GET app.localhost/
    T->>G: forwardAuth /auth
    G-->>T: 302 → idp/authorize (target in state)
    T-->>C: 302
    C->>I: GET id.localhost/authorize
    I-->>C: 302 → auth.localhost/callback?code&state
    C->>G: GET /callback?code&state
    G->>I: POST /token (code → access_token)
    G->>I: GET /userinfo (Bearer token)
    G-->>C: Set-Cookie gw_session; 302 → state

    Note over C,B: Later requests — with session cookie
    C->>T: GET app.localhost/ (Cookie gw_session)
    T->>G: forwardAuth /auth
    G-->>T: 200 + Remote-User
    T->>B: forward with Remote-User
    B-->>C: 200 hello page
```

The provider (`idp`) auto-approves as `alice@example.com`, so there is no real login screen to click
through. Point `authgw` at a real provider by overriding its env vars (see below).

## Run it

```bash
docker compose up -d --build
```

Host port `:80` is often taken, so the compose file maps Traefik's web entrypoint to **`:8000`**
(change it to `80:80` if the port is free). The routers match on `*.localhost` hostnames, so tests
pass an explicit `Host` header and no `/etc/hosts` edits are needed.

The Traefik dashboard is at <http://127.0.0.1:8080/dashboard/> once the stack is up.

### End-to-end smoke test

`cmd/smoke` is a small Go program that walks the full authorization-code flow, captures the session
cookie, then asserts identity injection, per-host routing, and the error page. It exits non-zero if any
check fails, so it works as a CI gate too.

```bash
go run ./cmd/smoke
# smoke-testing gateway at http://127.0.0.1:8000
#   ok   unauth request is redirected
#   ok   backend sees Remote-User
#   ok   original 503 preserved
#   ...
# passed: 13  failed: 0
```

Override the target with `BASE` (for the kind cluster below, `BASE=http://127.0.0.1:8001`).

### Tear down

```bash
docker compose down
```

## Run it on Kubernetes (kind)

The same setup runs on a cluster: the `dynamic.yml` routers and middlewares become Traefik CRDs under
`k8s/`. This was verified on a local [kind](https://kind.sigs.k8s.io/) cluster with Traefik v3.7.6
installed from the official Helm chart.

```bash
# 1. Cluster with a host port mapped to a NodePort for the web entrypoint.
cat > kind-config.yaml <<'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - { containerPort: 30080, hostPort: 8001, protocol: TCP }
EOF
kind create cluster --name traefik-demo --config kind-config.yaml

# 2. Build the Go images and load them into the cluster.
#    NB: use a tag WITHOUT ":latest" here. zsh treats `$s:latest` as the `:l`
#    (lowercase) modifier and mangles the tag; ":v1" is safe.
for s in echo authgw idp errors; do
  docker build --build-arg SERVICE="$s" -t "traefik-demo/$s:v1" .
  kind load docker-image "traefik-demo/$s:v1" --name traefik-demo
done

# 3. Install Traefik. The chart enables the kubernetesCRD provider and installs the CRDs.
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm install traefik traefik/traefik -n traefik --create-namespace \
  --set service.type=NodePort --set ports.web.nodePort=30080 --set ports.web.exposedPort=80

# 4. Apply the apps + CRDs, then run the same smoke test against port 8001.
kubectl apply -f k8s/
```

Then re-run the smoke test against the cluster with `BASE=http://127.0.0.1:8001 go run ./cmd/smoke`.
Tear the whole thing down with `kind delete cluster --name traefik-demo`.

Note the `apiVersion`: Traefik v3 uses `traefik.io/v1alpha1`. v2 used `traefik.containo.us/v1alpha1`,
and copying an old manifest is the usual reason an IngressRoute never shows up. On `helm upgrade`, Helm
does not update CRDs, so from chart v34 the CRDs ship in a separate `traefik-crds` chart you upgrade
first.

## Point it at a real provider

`authgw` reads its OAuth settings from env vars (defaults target the bundled `idp`). Set these on the
`authgw` service in `docker-compose.yml` to use Google, Okta, Auth0, etc.:

| Env var         | Meaning                                                        |
|-----------------|----------------------------------------------------------------|
| `AUTHORIZE_URL` | Provider's authorization endpoint (browser-facing).           |
| `TOKEN_URL`     | Token endpoint (server-to-server).                            |
| `USERINFO_URL`  | UserInfo endpoint (server-to-server).                         |
| `CLIENT_ID`     | OAuth client ID.                                              |
| `CLIENT_SECRET` | OAuth client secret.                                         |
| `CALLBACK_URL`  | This gateway's callback, e.g. `https://auth.example.com/callback`. |
| `COOKIE_DOMAIN` | Shared parent domain so the session covers every app (e.g. `example.com`). |

For real single sign-on across subdomains, all protected apps must share the parent domain in
`COOKIE_DOMAIN`; otherwise each app gets its own cookie and the session is not shared.

## Layout

```
traefik.yml        # install (static) config: entrypoints + file provider
dynamic.yml        # routing (dynamic) config: routers, middlewares, services
docker-compose.yml # the six services + port mapping
Dockerfile         # one image, built per service via the SERVICE build arg
authgw/main.go     # forward-auth endpoint + OAuth client + callback
idp/main.go        # minimal OAuth 2.0 provider (/authorize, /token, /userinfo)
echo/main.go       # backend that echoes headers; /status/{code} returns that code
errors/main.go     # static custom error pages
cmd/smoke/main.go  # end-to-end smoke test: walks the OAuth flow and asserts the result
k8s/               # the same setup as Traefik CRDs: apps, middlewares, ingressroutes
```

All four Go services are standard-library only, one `main.go` each. Build a single service locally with
`go build ./authgw` (or `echo`, `idp`, `errors`).
