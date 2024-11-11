//go:build tools

package tools

import (
	_ "github.com/vektra/mockery/v2"
)

// For mocks
//go:generate go run github.com/vektra/mockery/v2 --config=mockery.yaml
