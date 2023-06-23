package main

import (
	"embed"
	"fmt"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/joho/godotenv"
)

//go:embed files/*.sql
var migFs embed.FS

func main() {
	godotenv.Load()
	helpers.ReadPgxConnEnvs()
	dsn := helpers.ToDsnWithDbName()
	dir := "files"

	dirs, err := migFs.ReadDir(dir)
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
		Dir:     dir,
	}

	driver, err := postgres.New(dsn)
	if err != nil {
		panic(fmt.Sprintf("driver err := %s", err))
	}
	applied, err := migration.Migrate(driver, embededSource, migration.Up, 0)
	if err != nil {
		panic(fmt.Sprintf("Error applying migrations: %s", err.Error()))
	} else {
		println(fmt.Sprintf("last applied %d", applied))
	}
}
