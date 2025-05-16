package migrations

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func Up(connectionString string) error {
	d, err := iofs.New(fs, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connectionString)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
