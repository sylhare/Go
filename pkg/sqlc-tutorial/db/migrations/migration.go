package migrations

import (
	"fmt"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"net/url"
	"path/filepath"
	"runtime"
)

var migrationsDir string

// Called when the script is run
func init() {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsDir = filepath.Join(basePath)
}

func Run(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return err
	}
	db := dbmate.New(u)
	db.MigrationsDir = []string{migrationsDir}

	query := u.Query()
	if query.Get("sslmode") == "" {
		query.Set("sslmode", "disable")
		u.RawQuery = query.Encode()
	}

	err = db.CreateAndMigrate()
	if err != nil {
		fmt.Println("Migration failed", err)
		return err
	}

	return nil
}
