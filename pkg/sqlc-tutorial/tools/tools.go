//go:build tools
package tools

import (
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)

// SQL code generation
//go:generate sqlc generate
