package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func CreateMigration(name string) error {
	dir := "cmd/migrate/migrations"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	timestamp := time.Now().Unix()
	up := filepath.Join(dir, fmt.Sprintf("%d_%s.up.sql", timestamp, name))
	down := filepath.Join(dir, fmt.Sprintf("%d_%s.down.sql", timestamp, name))

	os.WriteFile(up, []byte("-- write your UP migration here\n"), 0644)
	os.WriteFile(down, []byte("-- write your DOWN migration here\n"), 0644)

	fmt.Println("Created:")
	fmt.Println("  ", up)
	fmt.Println("  ", down)

	return nil
}

func RunMigrations(app *application, versionArg string, action string) {

	dbURL := app.config.db.addr

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("DB pool create error: ", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Migration driver error:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Migration instance error:", err)
	}	

	switch action {
		case "up":
			fmt.Println("Running migrations UP...")
			if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {	
				log.Fatal("Migration up error:", err)
			}
			fmt.Println("Migrations applied successfully.")
		
		case "down":
			fmt.Println("Running migrations DOWN...")
			if err := m.Down(); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Migration DOWN complete")

		case "force":
			v, _ := strconv.Atoi(versionArg)

			fmt.Println("Forcing migration version to:", v)
			if err := m.Force(v); err != nil {
				log.Fatal(err)
			}

		default:
			fmt.Println("Unknown migration command. Use: up | down | force 1.0.0")
		}

}