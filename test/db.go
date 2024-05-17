package test

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func ResetAndSeedDB(db *sqlx.DB) error {
	entries, err := os.ReadDir("./database/migrations/")
	if err != nil {
		return err
	}
	err = dropAllTables(entries, db)
	if err != nil {
		return err
	}
	err = createAllTables(entries, db)
	if err != nil {
		return err
	}
	return nil
}

func dropAllTables(entries []fs.DirEntry, db *sqlx.DB) error {
	files := []string{}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".down.") {
			files = append(files, entry.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	for _, file := range files {
		err := runAllStatements(file, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func createAllTables(entries []fs.DirEntry, db *sqlx.DB) error {
	files := []string{}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".up.") {
			files = append(files, entry.Name())
		}
	}
	for _, file := range files {
		err := runAllStatements(file, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func runAllStatements(file string, db *sqlx.DB) error {
	b, err := os.ReadFile(fmt.Sprintf("./database/migrations/%s", file))
	if err != nil {
		return err
	}
	statements := strings.Split(string(b), "--> statement-breakpoint")
	for _, statement := range statements {
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
	}
	return nil
}
