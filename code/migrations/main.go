package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/mysql"
)

var migFs embed.FS

func main() {
	db, err := sql.Open("mysql", helpers.BuildPostgresConnUrl())
	if err != nil {
		panic(err)
	}

	driver, err := mysql.NewFromDB(db)
	if err != nil {
		panic(err)
	}

	dirs, err := migFs.ReadDir("sql")
	if err != nil {
		panic(fmt.Sprintf("error reading migration files from embeded: %v", err))
	} else {
		println("List of migrations found: ")
		for _, d := range dirs {
			println(fmt.Sprintf(" - %s", d.Name()))
		}
		println("End of list")
	}

	embededSource := &migration.EmbedMigrationSource{
		EmbedFS: migFs,
		Dir:     "sql",
	}

	applied, err := migration.Migrate(driver, embededSource, migration.Up, 0)
	if err != nil {
		panic(fmt.Sprintf("Error applying migrations: %s", err.Error()))
	} else {
		println(fmt.Sprintf("last applied %d", applied))
	}
}
