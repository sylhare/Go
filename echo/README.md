# Go echo server

From https://echo.labstack.com/docs/quick-start

Start the app with:

```bash
go run server.go
```

The app should be at http://localhost:1323/

## Rest Open API

Install the [openapi codegen][1] for go with:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest  
```

Generate the OpenAPI spec with:

```bash
oapi-codegen -generate types,server,spec -package api -o api/api.gen.go api/api.yml
```

[1] https://github.com/oapi-codegen/oapi-codegen