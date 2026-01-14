package models

import (
	"context"
	"io/fs"

	"ariga.io/atlas-go-sdk/atlasexec"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and runs Atlas migrations.
// The migrations fs.FS should contain migration files at the root level
// (use fs.Sub if your embed has them in a subdirectory).
func InitDB(path string, migrations fs.FS) (*gorm.DB, error) {
	dbURL := "sqlite://" + path

	// Run Atlas migrations
	if err := runMigrations(context.Background(), dbURL, migrations); err != nil {
		return nil, err
	}

	// Open GORM connection
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// runMigrations applies Atlas versioned migrations
func runMigrations(ctx context.Context, dbURL string, migrations fs.FS) error {
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(migrations),
	)
	if err != nil {
		return err
	}
	defer workdir.Close()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return err
	}

	_, err = client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL: dbURL,
	})
	return err
}
